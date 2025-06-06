package withdrawal

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/samstringzz/alutamarket-backend/errors"
	"github.com/samstringzz/alutamarket-backend/internals/shared"
	"github.com/samstringzz/alutamarket-backend/internals/user"
)

type Service interface {
	CreateWithdrawal(ctx context.Context, req *shared.NewWithdrawal) (*shared.Withdrawal, error)
	GetWithdrawal(ctx context.Context, id uint32) (*shared.Withdrawal, error)
	GetStoreWithdrawals(ctx context.Context, storeID uint32) ([]*shared.Withdrawal, error)
	GetPendingWithdrawals(ctx context.Context) ([]*shared.Withdrawal, error)
	ProcessPendingWithdrawals(ctx context.Context) error
	GetWithdrawals(ctx context.Context, status *string) ([]*shared.Withdrawal, error)
	ProcessWithdrawal(ctx context.Context, withdrawalID uint32, action string) error
}

type service struct {
	repo           Repository
	paystackClient user.PaystackClient
}

func NewService(repo Repository) Service {
	return &service{
		repo:           repo,
		paystackClient: user.NewPaystackClient(os.Getenv("PAYSTACK_SECRET_KEY")),
	}
}

func (s *service) CreateWithdrawal(ctx context.Context, req *shared.NewWithdrawal) (*shared.Withdrawal, error) {
	return s.repo.CreateWithdrawal(ctx, req)
}

func (s *service) GetWithdrawal(ctx context.Context, id uint32) (*shared.Withdrawal, error) {
	return s.repo.GetWithdrawal(ctx, id)
}

func (s *service) GetStoreWithdrawals(ctx context.Context, storeID uint32) ([]*shared.Withdrawal, error) {
	return s.repo.GetStoreWithdrawals(ctx, storeID)
}

func (s *service) GetPendingWithdrawals(ctx context.Context) ([]*shared.Withdrawal, error) {
	return s.repo.GetPendingWithdrawals(ctx)
}

func (s *service) ProcessPendingWithdrawals(ctx context.Context) error {
	withdrawals, err := s.repo.GetPendingWithdrawals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending withdrawals: %v", err)
	}

	for _, withdrawal := range withdrawals {
		if time.Since(withdrawal.CreatedAt) < 24*time.Hour {
			continue
		}

		if err := s.repo.ApproveWithdrawal(ctx, withdrawal.ID); err != nil {
			log.Printf("Failed to approve withdrawal: %v", err)
			continue
		}

		transfer, err := s.paystackClient.InitiateTransfer(ctx, &user.TransferRequest{
			Amount:    withdrawal.Amount,
			Recipient: withdrawal.AccountNumber,
			Reason:    "Withdrawal from Aluta Market",
		})
		if err != nil {
			if rejectErr := s.repo.RejectWithdrawal(ctx, withdrawal.ID,
				fmt.Sprintf("Failed to process transfer: %v", err)); rejectErr != nil {
				log.Printf("Failed to reject withdrawal after transfer failure: %v", rejectErr)
			}
			continue
		}

		if err := s.repo.CompleteWithdrawal(ctx, withdrawal.ID, transfer.Data.TransferCode); err != nil {
			log.Printf("Failed to complete withdrawal: %v", err)
			continue
		}
	}

	return nil
}

func (s *service) GetWithdrawals(ctx context.Context, status *string) ([]*shared.Withdrawal, error) {
	return s.repo.GetWithdrawals(ctx, status)
}

// ProcessWithdrawal handles admin actions on withdrawals (approve/reject)
func (s *service) ProcessWithdrawal(ctx context.Context, withdrawalID uint32, action string) error {
	// Get the withdrawal details
	w, err := s.repo.GetWithdrawal(ctx, withdrawalID)
	if err != nil {
		return fmt.Errorf("withdrawal not found: %v", err)
	}

	// Validate the action based on current status
	switch action {
	case "approve":
		if w.Status != string(shared.StatusPending) {
			return errors.NewAppError(400, "INVALID_STATUS", "Withdrawal must be pending to be approved")
		}

		// Get list of supported banks from Paystack
		banks, err := s.paystackClient.GetBanks(ctx)
		if err != nil {
			return fmt.Errorf("failed to get supported banks: %v", err)
		}

		// Find the bank code for the withdrawal's bank
		var bankCode string
		for _, bank := range banks.Data {
			if strings.EqualFold(bank.Name, w.BankName) {
				bankCode = bank.Code
				break
			}
		}

		if bankCode == "" {
			// If bank code not found, reject the withdrawal
			if rejectErr := s.repo.RejectWithdrawal(ctx, withdrawalID,
				fmt.Sprintf("Unsupported bank: %s. Please contact support for assistance.", w.BankName)); rejectErr != nil {
				log.Printf("Failed to reject withdrawal after bank code lookup failure: %v", rejectErr)
			}
			return fmt.Errorf("unsupported bank: %s", w.BankName)
		}

		// First, create a Paystack recipient
		recipient, err := s.paystackClient.CreateTransferRecipient(ctx, &user.RecipientRequest{
			Type:          "nuban",
			Name:          w.AccountName,
			AccountNumber: w.AccountNumber,
			BankCode:      bankCode,
		})
		if err != nil {
			// If recipient creation fails, reject the withdrawal
			if rejectErr := s.repo.RejectWithdrawal(ctx, withdrawalID,
				fmt.Sprintf("Failed to create transfer recipient: %v", err)); rejectErr != nil {
				log.Printf("Failed to reject withdrawal after recipient creation failure: %v", rejectErr)
			}
			return fmt.Errorf("failed to create transfer recipient: %v", err)
		}

		// Approve the withdrawal (updates status to 'approved')
		if err := s.repo.ApproveWithdrawal(ctx, withdrawalID); err != nil {
			return fmt.Errorf("failed to approve withdrawal: %v", err)
		}

		// Initiate Paystack transfer using the recipient code
		transfer, err := s.paystackClient.InitiateTransfer(ctx, &user.TransferRequest{
			Amount:    w.Amount,
			Recipient: recipient.Data.RecipientCode, // Use the recipient code from Paystack
			Reason:    fmt.Sprintf("Withdrawal #%d from Aluta Market", w.ID),
		})

		if err != nil {
			// If transfer fails, reject the withdrawal and refund
			if rejectErr := s.repo.RejectWithdrawal(ctx, withdrawalID,
				fmt.Sprintf("Automated transfer failed: %v", err)); rejectErr != nil {
				log.Printf("Failed to reject withdrawal after transfer failure: %v", rejectErr)
			}
			return fmt.Errorf("paystack transfer failed: %v", err)
		}

		// If transfer is successful, complete the withdrawal
		if err := s.repo.CompleteWithdrawal(ctx, withdrawalID, transfer.Data.TransferCode); err != nil {
			log.Printf("Warning: Failed to complete withdrawal after successful transfer: %v", err)
			// TODO: Consider alerting/manual intervention if completion fails after transfer
		}

	case "reject":
		if w.Status != string(shared.StatusPending) {
			return errors.NewAppError(400, "INVALID_STATUS", "Withdrawal must be pending to be rejected")
		}
		// Reject the withdrawal (refunds and updates status to 'rejected')
		if err := s.repo.RejectWithdrawal(ctx, withdrawalID, "Rejected by admin"); err != nil {
			return fmt.Errorf("failed to reject withdrawal: %v", err)
		}

	default:
		return errors.NewAppError(400, "INVALID_ACTION", "Invalid withdrawal action. Use 'approve' or 'reject'.")
	}

	return nil
}
