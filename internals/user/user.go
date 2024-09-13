package user

import (
	"context"
	"time"

	"github.com/Chrisentech/aluta-market-api/internals/store"
	"gorm.io/gorm"
)

type PaymentDetails struct {
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
	Info    string `json:"info" db:"info"`
}

// type Product *product.Product
type User struct {
	gorm.Model
	ID             uint32         `gorm:"primaryKey;uniqueIndex;not null;autoIncrement" json:"id" db:"id"` // Unique identifier for the user
	Campus         string         `json:"campus" db:"campus"`                                              // Campus of the user
	Email          string         `json:"email" db:"email"`                                                // Email address of the user
	Password       string         `json:"password" db:"password"`                                          // Password of the user
	Fullname       string         `json:"fullname" db:"fullname"`                                          // Full name of the user
	Phone          string         `json:"phone" db:"phone"`                                                // Phone number of the user
	Avatar         string         `json:"avatar" db:"avatar"`                                              // Phone number of the user
	Usertype       string         `json:"usertype" db:"usertype"`                                          // Type of user (e.g., seller,buyer,admin)
	Dob            string         `json:"dob" db:"dob"`                                                    // Type of user (e.g., seller,buyer,admin)
	Gender         string         `json:"gender" db:"gender"`                                              // Type of user (e.g., seller,buyer,admin)
	Active         *bool          `json:"active" db:"active"`
	Twofa          *bool          `json:"twofa" db:"twofa"` // Two factor authentication
	AccessToken    string         `json:"access_token,omitempty" db:"access_token"`
	RefreshToken   string         `json:"refresh_token,omitempty" db:"refresh_token"`
	FollowedStores []string       `gorm:"serializer:json" json:"followed_stores" db:"followed_stores"`
	Code           string         `json:"code,omitempty" db:"code"` // otp code for verifications
	Online         bool           `json:"online,omitempty" db:"online"`
	PaymentDetails PaymentDetails `gorm:"serializer:json"`
	Codeexpiry     time.Time      `json:"codeexpiry,omitempty" db:"codeexpiry"` // Expiry time for otpCode
	CreatedAt      time.Time      // Set to current time if it is zero on creating
}

type CreateUserReq struct {
	// ID                 uint32    `gorm:"primaryKey;uniqueIndex;not null;autoIncrement" json:"id" db:"id"` // Unique identifier for the user
	Campus             string    `json:"campus" db:"campus"`     // Campus of the user
	Email              string    `json:"email" db:"email"`       // Email address of the user
	Password           string    `json:"password" db:"password"` // Password of the user
	Fullname           string    `json:"fullname" db:"fullname"`
	Phone              string    `json:"phone" db:"phone"`                     // Phone number of the user
	Usertype           string    `json:"usertype" db:"usertype"`               // Type of user (e.g., seller,buyer,admin)
	Active             bool      `json:"active" db:"active"`                   // For accesibility of user,
	Twofa              bool      `json:"twofa" db:"twofa"`                     // Two factor authentication
	Code               string    `json:"code,omitempty" db:"code"`             // otp code for verifications
	Codeexpiry         time.Time `json:"codeexpiry,omitempty" db:"codeexpiry"` // Expiry time for otpCode
	StoreName          string    `json:"store_name" db:"name"`
	StoreEmail         string    `json:"store_email" db:"name"`
	StoreUser          uint32    `json:"user" db:"user_id"`
	StoreLink          string    `json:"link" db:"link"`
	StorePhone         string    `json:"store_phone" db:"store_phone"`
	Description        string    `json:"description" db:"description"`
	StoreAddress       string    `json:"address" db:"store_address"`
	HasPhysicalAddress bool      `json:"has_physical_address" db:"has_physical_address"`
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
	User       User   `json:"user"`
	StoreName  string `json:"store_name" db:"store_name"`
	StoreEmail string `json:"store_email" db:"store_email"`
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
	CreateStore(ctx context.Context, req *store.Store) (*store.Store, error)
	GetMyDVA(ctx context.Context, userEmail string) (*Account, error)
	SetPaymentDetais(ctx context.Context, req *PaymentDetails, userId uint32) error
	SendPasswordResetLink(ctx context.Context, req *PasswordReset) error
	UpdatePassword(ctx context.Context, req *PasswordReset) error
	VerifyResetLink(ctx context.Context, token string) error
	GetBalance(ctx context.Context, userId string) error
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
	CreateStore(ctx context.Context, req *store.Store) (*store.Store, error)
	SetPaymentDetais(ctx context.Context, req *PaymentDetails, userId uint32) error
	SendPasswordResetLink(ctx context.Context, req *PasswordReset) error
	UpdatePassword(ctx context.Context, req *PasswordReset) error
	GetBalance(ctx context.Context, userId string) error
	VerifyResetLink(ctx context.Context, token string) error
}
