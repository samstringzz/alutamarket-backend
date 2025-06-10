package user

import (
	"context"
	"time"

	"github.com/samstringzz/alutamarket-backend/graph/model"
	"github.com/samstringzz/alutamarket-backend/internals/models"
	"gorm.io/gorm"
)

type PaymentDetails struct {
	Name    string `json:"Name" db:"name"`
	Phone   string `json:"Phone" db:"phone"`
	Address string `json:"Address" db:"address"`
	Info    string `json:"Info" db:"info"`
}

type followedStores struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Thumbnail   string `json:"thumbnail" db:"thumbnail"`
	Background  string `json:"background" db:"background"`
	Link        string `json:"link" db:"link"`
}

// Add this struct
type DVAAccount struct {
	ID            string      `json:"id" gorm:"type:varchar(100);primary_key"`
	AccountName   string      `json:"account_name" gorm:"type:varchar(100)"`
	AccountNumber string      `json:"account_number" gorm:"type:varchar(20)"`
	CustomerID    string      `json:"customer_id" gorm:"type:varchar(100)"`
	BankID        string      `json:"bank_id" gorm:"type:varchar(100)"`
	Customer      DVACustomer `json:"customer" gorm:"foreignKey:CustomerID"`
	Bank          DVABank     `json:"bank" gorm:"foreignKey:BankID"`
}

type DVACustomer struct {
	ID    string `json:"id" gorm:"type:varchar(100);primary_key"`
	Email string `json:"email" gorm:"type:varchar(100)"`
}

type DVABank struct {
	ID   string `json:"id" gorm:"type:varchar(100);primary_key"`
	Name string `json:"name" gorm:"type:varchar(100)"`
	Slug string `json:"slug" gorm:"type:varchar(50)"`
}

// type Product *product.Product
type User struct {
	gorm.Model
	ID             uint32           `gorm:"primaryKey;uniqueIndex;not null;autoIncrement" json:"id" db:"id"` // Unique identifier for the user
	Campus         string           `json:"campus" db:"campus"`                                              // Campus of the user
	Email          string           `json:"email" db:"email"`                                                // Email address of the user
	Password       string           `json:"password" db:"password"`                                          // Password of the user
	Fullname       string           `json:"fullname" db:"fullname"`                                          // Full name of the user
	UUID           string           `json:"uuid" db:"uuid"`                                                  // Full name of the user
	Phone          string           `json:"phone" db:"phone"`                                                // Phone number of the user
	Avatar         string           `json:"avatar" db:"avatar"`                                              // Phone number of the user
	Usertype       string           `json:"usertype" db:"usertype"`                                          // Type of user (e.g., seller,buyer,admin)
	Dob            string           `json:"dob" db:"dob"`                                                    // Type of user (e.g., seller,buyer,admin)
	Gender         string           `json:"gender" db:"gender"`                                              // Type of user (e.g., seller,buyer,admin)
	Active         *bool            `json:"active" db:"active"`
	Twofa          *bool            `json:"twofa" db:"twofa"` // Two factor authentication
	AccessToken    string           `json:"access_token,omitempty" db:"access_token"`
	RefreshToken   string           `json:"refresh_token,omitempty" db:"refresh_token"`
	FollowedStores []followedStores `gorm:"serializer:json" json:"followed_stores" db:"followed_stores"`
	Code           string           `json:"code,omitempty" db:"code"` // otp code for verifications
	Online         bool             `json:"online,omitempty" db:"online"`
	PaymentDetails PaymentDetails   `gorm:"serializer:json"`
	Codeexpiry     time.Time        `json:"codeexpiry,omitempty" db:"codeexpiry"` // Expiry time for otpCode
	CreatedAt      time.Time        // Set to current time if it is zero on creating
}

// Add this type definition after the User struct
type UpdateUserReq struct {
	ID                 uint32         `json:"id"`
	Fullname           string         `json:"fullname,omitempty"`
	Email              string         `json:"email,omitempty"`
	Campus             string         `json:"campus,omitempty"`
	Phone              string         `json:"phone,omitempty"`
	Avatar             string         `json:"avatar,omitempty"`
	Dob                string         `json:"dob,omitempty"`
	PaymentDetails     PaymentDetails `json:"payment_details,omitempty"`
	StoreName          string         `json:"store_name,omitempty"`
	StoreEmail         string         `json:"store_email,omitempty"`
	HasPhysicalAddress bool           `json:"has_physical_address,omitempty"`
}

type CreateUserReq struct {
	// ID                 uint32    `gorm:"primaryKey;uniqueIndex;not null;autoIncrement" json:"id" db:"id"` // Unique identifier for the user
	Campus             string         `json:"campus" db:"campus"`     // Campus of the user
	Email              string         `json:"email" db:"email"`       // Email address of the user
	Password           string         `json:"password" db:"password"` // Password of the user
	Fullname           string         `json:"fullname" db:"fullname"`
	Phone              string         `json:"phone" db:"phone"`                     // Phone number of the user
	Usertype           string         `json:"usertype" db:"usertype"`               // Type of user (e.g., seller,buyer,admin)
	Active             bool           `json:"active" db:"active"`                   // For accesibility of user,
	Twofa              bool           `json:"twofa" db:"twofa"`                     // Two factor authentication
	UUID               string         `json:"uuid" db:"uuid"`                       // Add UUID field
	Avatar             string         `json:"avatar" db:"avatar"`                   // Add Avatar field
	Online             bool           `json:"online" db:"online"`                   // Add Online field
	FollowedStores     []string       `json:"followed_stores" db:"followed_stores"` // Add FollowedStores field
	Code               string         `json:"code,omitempty" db:"code"`             // otp code for verifications
	Codeexpiry         time.Time      `json:"codeexpiry,omitempty" db:"codeexpiry"` // Expiry time for otpCode
	StoreName          string         `json:"store_name" db:"name"`
	StoreEmail         string         `json:"store_email" db:"name"`
	StoreUser          uint32         `json:"user" db:"user_id"`
	StoreLink          string         `json:"link" db:"link"`
	StorePhone         string         `json:"store_phone" db:"store_phone"`
	Description        string         `json:"description" db:"description"`
	StoreAddress       string         `json:"address" db:"store_address"`
	HasPhysicalAddress bool           `json:"has_physical_address" db:"has_physical_address"`
	Dob                *string        `json:"dob" db:"dob"`       // Add date of birth field
	Gender             string         `json:"gender" db:"gender"` // Add gender field
	PaymentDetails     PaymentDetails `json:"payment_details"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"` // Add created_at field
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"` // Add updated_at field
}

type CreateUserRes struct {
	Message string
	Status  int
	Data    interface{}
}

type FilterOption struct {
	Email string `json:"email" db:"email"`
	Phone string `json:"phone" db:"phone"`
	ID    string `json:"id" db:"id"`
}
type LoginUserReq struct {
	Password string `json:"password" db:"password"` // Password of the user for login
	Email    string `json:"email" db:"email"`       // Email address of the user for login
}

type LoginUserRes struct {
	AccessToken  string // Access token generated for the logged-in user
	RefreshToken string // RefreshToken token generated for the logged-in user
	ID           uint32 `json:"id" db:"id"` // Unique identifier for the logged-in user
}

type Customer struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	CustomerCode string `json:"customer_code"`
	Phone        string `json:"phone"`
	RiskAction   string `json:"risk_action"`
	// InternationalFormatPhone interface{} `json:"international_format_phone"`
}

type Bank struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

type VerifyOTPReq struct {
	Code     string `json:"code," db:"code"`  // otp code for verifications
	Email    string `json:"email" db:"email"` // Email address of the user for login
	Phone    string `json:"phone" db:"phone"` // Email address of the user for login
	Attempts int    `json:"attempts" db:"attempts"`
}

type SplitConfig struct {
	Subaccount string `json:"subaccount"`
}

type Account struct {
	ID            int         `json:"id"`
	AccountNumber string      `json:"account_number"`
	AccountName   string      `json:"account_name"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	Active        bool        `json:"active"`
	Assigned      bool        `json:"assigned"`
	Customer      *Customer   `json:"customer"`
	Bank          *Bank       `json:"bank"`
	SplitConfig   SplitConfig `json:"split_config"`
}

type DVADetails struct {
	User          User   `json:"user"`
	StoreName     string `json:"store_name" db:"store_name"`
	StoreEmail    string `json:"store_email" db:"store_email"`
	AccountNumber string `json:"account_number" db:"account_number"`
	BankName      string `json:"bank_name" db:"bank_name"`
	CustomerCode  string `json:"customer_code" db:"customer_code"`
	AccountName   string `json:"account_name" db:"account_name"`
	PaystackData  string `json:"paystack_data" db:"paystack_data"`
}

// Transaction model
type Transaction struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	Reference string `json:"reference"`
	Amount    int    `json:"amount"`
	Customer  struct {
		ID           int    `json:"id"`
		CustomerCode string `json:"customer_code"`
		Email        string `json:"email"`
	} `json:"customer"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionsResponse struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"data"`
}

type PasswordReset struct {
	gorm.Model
	Link      string    `json:"link" db:"link"`
	ExpiresAt time.Time `json:"codeexpiry,omitempty" db:"codeexpiry"` // Expiry time for otpCode
	Token     string    `json:"token,omitempty" db:"token"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password,omitempty" db:"password"`
}
type Repository interface {
	CreateUser(ctx context.Context, user *CreateUserReq) (*User, error)     // Create a new user
	GetUsers(ctx context.Context) ([]*User, error)                          // Create a new user
	GetUser(ctx context.Context, filter string) (*User, error)              // Create a new user
	GetUserByEmailOrPhone(ctx context.Context, email string) (*User, error) // Get user information by email
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)      // Perform user login
	VerifyOTP(ctx context.Context, req *VerifyOTPReq) (*LoginUserRes, error)
	ToggleStoreFollowStatus(ctx context.Context, userId, storeId uint32) error
	UpdateUser(ctx context.Context, user *User) (*User, error)
	CreateDVAAccount(ctx context.Context, req *DVADetails) (string, error)
	CreateStore(ctx context.Context, req *models.Store) (*models.Store, error)
	GetMyDVA(ctx context.Context, userEmail string) (*Account, error)
	SetPaymentDetais(ctx context.Context, req *PaymentDetails, userId uint32) error
	SendPasswordResetLink(ctx context.Context, req *PasswordReset) error
	UpdatePassword(ctx context.Context, req *PasswordReset) error
	VerifyResetLink(ctx context.Context, token string) error
	GetBalance(ctx context.Context, userId string) error
	ConfirmPassword(ctx context.Context, password, userId string) error
	GetMyDownloads(ctx context.Context, userId string) ([]*models.Downloads, error)
	SendMaintenanceMail(ctx context.Context, userId string, active bool) error
	GetDB() *gorm.DB
	GetTransactionsByCustomerID(customerID string) ([]Transaction, error)
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) // Create a new user
	GetUsers(ctx context.Context) ([]*User, error)
	GetUser(ctx context.Context, filter string) (*User, error) // Create a new user
	VerifyOTP(ctx context.Context, req *VerifyOTPReq) (*LoginUserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) // Perform user login
	UpdateUser(ctx context.Context, user *User) (*User, error)
	ToggleStoreFollowStatus(ctx context.Context, userId, storeId uint32) error
	CreateDVAAccount(ctx context.Context, req *DVADetails) (string, error)
	GetMyDVA(ctx context.Context, userEmail string) (*Account, error)
	CreateStore(ctx context.Context, req *models.Store) (*models.Store, error)
	SetPaymentDetais(ctx context.Context, req *PaymentDetails, userId uint32) error
	SendPasswordResetLink(ctx context.Context, req *PasswordReset) error
	UpdatePassword(ctx context.Context, req *PasswordReset) error
	GetBalance(ctx context.Context, userId string) error
	ConfirmPassword(ctx context.Context, password, userId string) error
	VerifyResetLink(ctx context.Context, token string) error
	GetMyDownloads(ctx context.Context, userId string) ([]*models.Downloads, error)
	SendMaintenanceMail(ctx context.Context, userId string, active bool) error
	GetDB() *gorm.DB
	GetPaystackDepositTransactions(ctx context.Context, storeEmail string) ([]*model.DepositTransaction, error)
}
