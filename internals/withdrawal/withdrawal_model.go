package withdrawal

import (
	"time"
)

type Withdrawal struct {
	ID                    uint32     `json:"id" gorm:"primaryKey"`
	StoreID               uint32     `json:"store_id"`
	Amount                float64    `json:"amount"`
	Status                string     `json:"status"`
	PaystackTransferID    string     `json:"paystack_transfer_id"`
	PaystackRecipientCode string     `json:"paystack_recipient_code"`
	BankName              string     `json:"bank_name"`
	AccountNumber         string     `json:"account_number"`
	AccountName           string     `json:"account_name"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	ApprovedAt            *time.Time `json:"approved_at"`
	CompletedAt           *time.Time `json:"completed_at"`
	RejectionReason       string     `json:"rejection_reason"`
}

type NewWithdrawal struct {
	StoreID       uint32  `json:"store_id"`
	Amount        float64 `json:"amount"`
	BankName      string  `json:"bank_name"`
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
}

type WithdrawalStatus string

const (
	StatusPending   WithdrawalStatus = "pending"
	StatusApproved  WithdrawalStatus = "approved"
	StatusRejected  WithdrawalStatus = "rejected"
	StatusCompleted WithdrawalStatus = "completed"
)
