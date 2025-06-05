package shared

import "time"

type WithdrawalStatus string

const (
	StatusPending   WithdrawalStatus = "pending"
	StatusApproved  WithdrawalStatus = "approved"
	StatusRejected  WithdrawalStatus = "rejected"
	StatusCompleted WithdrawalStatus = "completed"
)

type Withdrawal struct {
	ID                 uint32     `gorm:"primaryKey" json:"id"`
	StoreID            uint32     `json:"store_id"`
	Amount             float64    `json:"amount"`
	Status             string     `json:"status"`
	BankName           string     `json:"bank_name"`
	AccountNumber      string     `json:"account_number"`
	AccountName        string     `json:"account_name"`
	RejectionReason    string     `json:"rejection_reason,omitempty"`
	PaystackTransferID string     `json:"paystack_transfer_id,omitempty"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type NewWithdrawal struct {
	StoreID       uint32  `json:"store_id"`
	Amount        float64 `json:"amount"`
	BankName      string  `json:"bank_name"`
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
}
