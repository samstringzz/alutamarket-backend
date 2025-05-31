package withdrawal

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Chrisentech/aluta-market-api/internals/user"
)

type Service interface {
	CreateWithdrawal(ctx context.Context, req *NewWithdrawal) (*Withdrawal, error)
	GetWithdrawal(ctx context.Context, id uint32) (*Withdrawal, error)
	GetStoreWithdrawals(ctx context.Context, storeID uint32) ([]*Withdrawal, error)
	GetPendingWithdrawals(ctx context.Context) ([]*Withdrawal, error)
	ProcessPendingWithdrawals(ctx context.Context) error
}

type service struct {
	repo           Repository
	paystackClient user.PaystackClient // Changed from paystack to paystackClient
}

func NewService(repo Repository) Service {
	return &service{
		repo:           repo,
		paystackClient: user.NewPaystackClient(os.Getenv("PAYSTACK_SECRET_KEY")),
	}
}

func (s *service) CreateWithdrawal(ctx context.Context, req *NewWithdrawal) (*Withdrawal, error) {
	return s.repo.CreateWithdrawal(ctx, req)
}

func (s *service) GetWithdrawal(ctx context.Context, id uint32) (*Withdrawal, error) {
	return s.repo.GetWithdrawal(ctx, id)
}

func (s *service) GetStoreWithdrawals(ctx context.Context, storeID uint32) ([]*Withdrawal, error) {
	return s.repo.GetStoreWithdrawals(ctx, storeID)
}

func (s *service) GetPendingWithdrawals(ctx context.Context) ([]*Withdrawal, error) {
	return s.repo.GetPendingWithdrawals(ctx)
}

func (s *service) ProcessPendingWithdrawals(ctx context.Context) error {
	// Get all pending withdrawals
	withdrawals, err := s.repo.GetPendingWithdrawals(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending withdrawals: %v", err)
	}

	for _, withdrawal := range withdrawals {
		// Check if withdrawal is older than 24 hours
		if time.Since(withdrawal.CreatedAt) < 24*time.Hour {
			continue
		}

		// Approve the withdrawal
		if err := s.repo.ApproveWithdrawal(ctx, withdrawal.ID); err != nil {
			log.Printf("Failed to approve withdrawal: %v", err)
			continue
		}

		// Process the withdrawal through Paystack
		transfer, err := s.paystackClient.InitiateTransfer(ctx, &user.TransferRequest{
			Amount:    withdrawal.Amount,
			Recipient: withdrawal.PaystackRecipientCode,
			Reason:    "Withdrawal from Aluta Market",
		})
		if err != nil {
			// If transfer fails, reject the withdrawal
			if rejectErr := s.repo.RejectWithdrawal(ctx, withdrawal.ID,
				fmt.Sprintf("Failed to process transfer: %v", err)); rejectErr != nil {
				log.Printf("Failed to reject withdrawal after transfer failure: %v", rejectErr)
			}
			continue
		}

		// Complete the withdrawal
		withdrawal.PaystackTransferID = transfer.Data.TransferCode
		withdrawal.Status = "processing"
		if err := s.repo.CompleteWithdrawal(ctx, withdrawal.ID, transfer.Data.TransferCode); err != nil {
			log.Printf("Failed to complete withdrawal: %v", err)
			continue
		}
	}

	return nil
}
