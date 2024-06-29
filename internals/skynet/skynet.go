package skynet

import (
	"context"

	"gorm.io/gorm"
)

type Data struct {
	Amount        int64  `json:"amount" db:"amount"`
	UserID        uint32 `json:"user_id" db:"user_id"`
	RequestID     string `json:"request_id" db:"request_id"`
	BillersCode   string `json:"biller_code" db:"biller_code"` //The Person number the sub should go to
	VariationCode string `json:"variation_code" db:"variation_code"`
	ServiceID     string `json:"serviceID" db:"serviceID"`
	Phone         string `json:"phone" db:"phone"`
}

type Airtime struct {
	Amount    int64  `json:"amount" db:"amount"`
	UserID    uint32 `json:"user_id" db:"user_id"`
	RequestID string `json:"request_id" db:"request_id"`
	ServiceID string `json:"serviceID" db:"serviceID"`
	Phone     string `json:"phone" db:"phone"`
}

type Electricity struct {
	Amount      int64  `json:"amount" db:"amount"`
	UserID      string `json:"user_id" db:"user_id"`
	RequestID   string `json:"request_id" db:"request_id"`
	BillersCode string `json:"biller_code" db:"biller_code"` //The Person number the sub should go to
	ServiceID   string `json:"serviceID" db:"serviceID"`
	Phone       string `json:"phone" db:"phone"`
}

type TVSubscription struct {
	Amount           int64  `json:"amount" db:"amount"`
	UserID           string `json:"user_id" db:"user_id"`
	RequestID        string `json:"request_id" db:"request_id"`
	BillersCode      string `json:"biller_code" db:"biller_code"` //The Person number the sub should go to
	VariationCode    string `json:"variation_code" db:"variation_code"`
	ServiceID        string `json:"serviceID" db:"serviceID"`
	Phone            string `json:"phone" db:"phone"`
	SubscriptionType string `json:"subscription_type" db:"subscription_type"`
}

type EducationPayment struct {
	Amount        int64  `json:"amount" db:"amount"`
	UserID        uint32 `json:"user_id" db:"user_id"`
	RequestID     string `json:"request_id" db:"request_id"`
	BillersCode   string `json:"biller_code" db:"biller_code"` //The Person number the sub should go to
	VariationCode string `json:"variation_code" db:"variation_code"`
	ServiceID     string `json:"serviceID" db:"serviceID"`
	Phone         string `json:"phone" db:"phone"`
}

type Skynet struct {
	gorm.Model
	ID            string `gorm:"primaryKey;uniqueIndex;not null"`
	UserID        uint32 `json:"user_id" db:"user_id"`
	Status        string `json:"status" db:"status"`
	RequestID     string `json:"request_id" db:"request_id"`
	TransactionID string `json:"transactiont_id" db:"transactiont_id"`
	Type          string `json:"type" db:"type"`
	Receiver      string `json:"receiever" db:"receiever"`
}
type BundleVariation struct {
	VariationCode   string `json:"variation_code"`
	Name            string `json:"name"`
	VariationAmount string `json:"variation_amount"`
	FixedPrice      string `json:"fixedPrice"`
}

// Content represents the content part of the response
type DataBundle struct {
	ServiceName    string            `json:"ServiceName"`
	ServiceID      string            `json:"serviceID"`
	ConvinienceFee string            `json:"convinience_fee"`
	Variations     []BundleVariation `json:"variations"`
}
type SmartcardVerificationResponse struct {
	Code    string  `json:"code"`
	Content Content `json:"content"`
}

type Content struct {
	CustomerName       string  `json:"Customer_Name"`
	Status             string  `json:"Status"`
	DueDate            string  `json:"DUE_DATE"`
	CustomerNumber     int     `json:"Customer_Number"`
	CustomerType       string  `json:"Customer_Type"`
	CurrentBouquet     string  `json:"Current_Bouquet"`
	CurrentBouquetCode string  `json:"Current_Bouquet_Code"`
	RenewalAmount      float64 `json:"Renewal_Amount"`
}
type Repository interface {
	BuyAirtime(ctx context.Context, airtime *Airtime) (*string, error)
	BuyData(ctx context.Context, data *Data) (*string, error)
	GetSubscriptionsBundles(ctx context.Context, serviceId string) (*DataBundle, error)
	VerifySmartCard(ctx context.Context, serviceId, billersCode string) (*SmartcardVerificationResponse, error)
	// BuyElectricity(ctx context.Context, electricity *Electricity) (*string, error)
	// SubscribeToTV(ctx context.Context, subscription *TVSubscription) (*string, error)
	// PayForEducation(ctx context.Context, payment *EducationPayment) (*string, error)
	// GetBalance(ctx context.Context, filter string) (*Data, error)
}

type Service interface {
	BuyAirtime(ctx context.Context, airtime *Airtime) (*string, error)
	BuyData(ctx context.Context, data *Data) (*string, error)
	GetSubscriptionsBundles(ctx context.Context, serviceId string) (*DataBundle, error)
	VerifySmartCard(ctx context.Context, serviceId, billersCode string) (*SmartcardVerificationResponse, error)
	// BuyElectricity(ctx context.Context, electricity *Electricity) (*string, error)
	// SubscribeToTV(ctx context.Context, subscription *TVSubscription) (*string, error)
	// PayForEducation(ctx context.Context, payment *EducationPayment) (*string, error)
	// GetBalance(ctx context.Context, filter string) (*Data, error)
}
