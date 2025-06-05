package withdrawal

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Chrisentech/aluta-market-api/internals/shared"
	"github.com/Chrisentech/aluta-market-api/internals/user"
)

type Service interface {
	CreateWithdrawal(ctx context.Context, req *shared.NewWithdrawal) (*shared.Withdrawal, error)
	GetWithdrawal(ctx context.Context, id uint32) (*shared.Withdrawal, error)
	GetStoreWithdrawals(ctx context.Context, storeID uint32) ([]*shared.Withdrawal, error)
	GetPendingWithdrawals(ctx context.Context) ([]*shared.Withdrawal, error)
	ProcessPendingWithdrawals(ctx context.Context) error
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
