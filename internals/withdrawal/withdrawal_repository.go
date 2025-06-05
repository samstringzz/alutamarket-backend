package withdrawal

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Chrisentech/aluta-market-api/database"
	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/internals/shared"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/utils"
	"gorm.io/gorm"
)

type Repository interface {
	CreateWithdrawal(ctx context.Context, req *shared.NewWithdrawal) (*shared.Withdrawal, error)
	GetWithdrawal(ctx context.Context, id uint32) (*shared.Withdrawal, error)
	GetStoreWithdrawals(ctx context.Context, storeID uint32) ([]*shared.Withdrawal, error)
	GetPendingWithdrawals(ctx context.Context) ([]*shared.Withdrawal, error)
	ApproveWithdrawal(ctx context.Context, id uint32) error
	RejectWithdrawal(ctx context.Context, id uint32, reason string) error
	CompleteWithdrawal(ctx context.Context, id uint32, paystackTransferID string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

func (r *repository) CreateWithdrawal(ctx context.Context, req *shared.NewWithdrawal) (*shared.Withdrawal, error) {
	// Get store to check balance
	storeRepo := store.NewRepository()
	store, err := storeRepo.GetStore(ctx, req.StoreID)
	if err != nil {
		return nil, fmt.Errorf("failed to get store: %v", err)
	}

	// Update wallet balance before checking
	if err := storeRepo.UpdateWalletBalance(ctx, req.StoreID); err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Get updated store data
	store, err = storeRepo.GetStore(ctx, req.StoreID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated store: %v", err)
	}

	// Check if store has sufficient balance
	if store.Wallet < req.Amount {
		return nil, errors.NewAppError(400, "INSUFFICIENT_BALANCE", "Insufficient balance for withdrawal")
	}

	// Create withdrawal record
	withdrawal := &shared.Withdrawal{
		StoreID:       req.StoreID,
		Amount:        req.Amount,
		Status:        string(shared.StatusPending),
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	}

	if err := r.db.Create(withdrawal).Error; err != nil {
		return nil, fmt.Errorf("failed to create withdrawal: %v", err)
	}

	// Deduct amount from store balance
	// First try to deduct from Paystack balance, then from wallet if needed
	remainingAmount := req.Amount
	if store.PaystackBalance > 0 {
		deductFromPaystack := min(store.PaystackBalance, remainingAmount)
		if err := storeRepo.UpdatePaystackBalance(ctx, req.StoreID, -deductFromPaystack); err != nil {
			r.db.Delete(withdrawal)
			return nil, fmt.Errorf("failed to update Paystack balance: %v", err)
		}
		remainingAmount -= deductFromPaystack
	}

	if remainingAmount > 0 {
		if err := storeRepo.UpdateWallet(ctx, req.StoreID, -remainingAmount); err != nil {
			r.db.Delete(withdrawal)
			return nil, fmt.Errorf("failed to update store wallet: %v", err)
		}
	}

	// Update wallet balance after withdrawal
	if err := storeRepo.UpdateWalletBalance(ctx, req.StoreID); err != nil {
		log.Printf("Warning: failed to update wallet balance after withdrawal: %v", err)
	}

	// Send notification to admin
	utils.SendEmail("admin@alutamarket.com", "New Withdrawal Request", []string{}, map[string]string{
		"store_name":       store.Name,
		"amount":           fmt.Sprintf("%.2f", req.Amount),
		"wallet_balance":   fmt.Sprintf("%.2f", store.Wallet),
		"paystack_balance": fmt.Sprintf("%.2f", store.PaystackBalance),
		"total_balance":    fmt.Sprintf("%.2f", store.Wallet+store.PaystackBalance),
	})

	return withdrawal, nil
}

func (r *repository) GetWithdrawal(ctx context.Context, id uint32) (*shared.Withdrawal, error) {
	var withdrawal shared.Withdrawal
	if err := r.db.First(&withdrawal, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get withdrawal: %v", err)
	}
	return &withdrawal, nil
}

func (r *repository) GetStoreWithdrawals(ctx context.Context, storeID uint32) ([]*shared.Withdrawal, error) {
	var withdrawals []*shared.Withdrawal
	if err := r.db.Where("store_id = ?", storeID).Order("created_at DESC").Find(&withdrawals).Error; err != nil {
		return nil, fmt.Errorf("failed to get store withdrawals: %v", err)
	}
	return withdrawals, nil
}

func (r *repository) GetPendingWithdrawals(ctx context.Context) ([]*shared.Withdrawal, error) {
	var withdrawals []*shared.Withdrawal
	if err := r.db.Where("status = ?", shared.StatusPending).Order("created_at ASC").Find(&withdrawals).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending withdrawals: %v", err)
	}
	return withdrawals, nil
}

func (r *repository) ApproveWithdrawal(ctx context.Context, id uint32) error {
	withdrawal, err := r.GetWithdrawal(ctx, id)
	if err != nil {
		return err
	}

	if withdrawal.Status != string(shared.StatusPending) {
		return errors.NewAppError(400, "INVALID_STATUS", "Withdrawal is not in pending status")
	}

	now := time.Now()
	withdrawal.Status = string(shared.StatusApproved)
	withdrawal.ApprovedAt = &now

	if err := r.db.Save(withdrawal).Error; err != nil {
		return fmt.Errorf("failed to approve withdrawal: %v", err)
	}

	// Update wallet balance after approval
	storeRepo := store.NewRepository()
	if err := storeRepo.UpdateWalletBalance(ctx, withdrawal.StoreID); err != nil {
		log.Printf("Warning: failed to update wallet balance after approval: %v", err)
	}

	// Send notification to store
	store, err := storeRepo.GetStore(ctx, withdrawal.StoreID)
	if err != nil {
		return fmt.Errorf("failed to get store: %v", err)
	}

	utils.SendEmail(store.Email, "Withdrawal Approved", []string{}, map[string]string{
		"amount": fmt.Sprintf("%.2f", withdrawal.Amount),
	})

	return nil
}

func (r *repository) RejectWithdrawal(ctx context.Context, id uint32, reason string) error {
	withdrawal, err := r.GetWithdrawal(ctx, id)
	if err != nil {
		return err
	}

	if withdrawal.Status != string(shared.StatusPending) {
		return errors.NewAppError(400, "INVALID_STATUS", "Withdrawal is not in pending status")
	}

	// Refund the amount to store balance
	storeRepo := store.NewRepository()
	store, err := storeRepo.GetStore(ctx, withdrawal.StoreID)
	if err != nil {
		return fmt.Errorf("failed to get store: %v", err)
	}

	// Refund to Paystack balance first, then to wallet
	remainingAmount := withdrawal.Amount
	if store.PaystackBalance > 0 {
		refundToPaystack := min(store.PaystackBalance, remainingAmount)
		if err := storeRepo.UpdatePaystackBalance(ctx, withdrawal.StoreID, refundToPaystack); err != nil {
			return fmt.Errorf("failed to refund Paystack balance: %v", err)
		}
		remainingAmount -= refundToPaystack
	}

	if remainingAmount > 0 {
		if err := storeRepo.UpdateWallet(ctx, withdrawal.StoreID, remainingAmount); err != nil {
			return fmt.Errorf("failed to refund store wallet: %v", err)
		}
	}

	withdrawal.Status = string(shared.StatusRejected)
	withdrawal.RejectionReason = reason

	if err := r.db.Save(withdrawal).Error; err != nil {
		return fmt.Errorf("failed to reject withdrawal: %v", err)
	}

	utils.SendEmail(store.Email, "Withdrawal Rejected", []string{}, map[string]string{
		"amount": fmt.Sprintf("%.2f", withdrawal.Amount),
		"reason": reason,
	})

	return nil
}

func (r *repository) CompleteWithdrawal(ctx context.Context, id uint32, paystackTransferID string) error {
	withdrawal, err := r.GetWithdrawal(ctx, id)
	if err != nil {
		return err
	}

	if withdrawal.Status != string(shared.StatusApproved) {
		return errors.NewAppError(400, "INVALID_STATUS", "Withdrawal is not in approved status")
	}

	now := time.Now()
	withdrawal.Status = string(shared.StatusCompleted)
	withdrawal.CompletedAt = &now
	withdrawal.PaystackTransferID = paystackTransferID

	if err := r.db.Save(withdrawal).Error; err != nil {
		return fmt.Errorf("failed to complete withdrawal: %v", err)
	}

	// Update wallet balance after completion
	storeRepo := store.NewRepository()
	if err := storeRepo.UpdateWalletBalance(ctx, withdrawal.StoreID); err != nil {
		log.Printf("Warning: failed to update wallet balance after completion: %v", err)
	}

	// Send notification to store
	store, err := storeRepo.GetStore(ctx, withdrawal.StoreID)
	if err != nil {
		return fmt.Errorf("failed to get store: %v", err)
	}

	utils.SendEmail(store.Email, "Withdrawal Completed", []string{}, map[string]string{
		"amount": fmt.Sprintf("%.2f", withdrawal.Amount),
	})

	return nil
}

// Helper function to get minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
