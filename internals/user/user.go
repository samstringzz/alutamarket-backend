package user

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// type Product *product.Product
type User struct {
	gorm.Model
	ID           uint32    `gorm:"primaryKey;uniqueIndex;not null;autoIncrement" json:"id" db:"id"` // Unique identifier for the user
	Campus       string    `json:"campus" db:"campus"`                                              // Campus of the user
	Email        string    `json:"email" db:"email"`                                                // Email address of the user
	Password     string    `json:"password" db:"password"`                                          // Password of the user
	Fullname     string    `json:"fullname" db:"fullname"`                                          // Full name of the user
	Phone        string    `json:"phone" db:"phone"`                                                // Phone number of the user
	Avatar       string    `json:"avatar" db:"avatar"`                                              // Phone number of the user
	Usertype     string    `json:"usertype" db:"usertype"`                                          // Type of user (e.g., seller,buyer,admin)
	Dob          string    `json:"dob" db:"dob"`                                                    // Type of user (e.g., seller,buyer,admin)
	Gender       string    `json:"gender" db:"gender"`                                              // Type of user (e.g., seller,buyer,admin)
	Active       *bool     `json:"active" db:"active"`
	Twofa        *bool     `json:"twofa" db:"twofa"`                           // Two factor authentication
	AccessToken  string    `json:"access_token,omitempty" db:"access_token"`   // Balance of the user's wallet (only for seller)
	RefreshToken string    `json:"refresh_token,omitempty" db:"refresh_token"` // Balance of the user's wallet (only for seller)
	Code         string    `json:"code,omitempty" db:"code"`                   // otp code for verifications
	Codeexpiry   time.Time `json:"codeexpiry,omitempty" db:"codeexpiry"`       // Expiry time for otpCode
	CreatedAt    time.Time // Set to current time if it is zero on creating
}

type CreateUserReq struct {
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
	StoreName          string    `json:"name" db:"name"`
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

type DVADetails struct {
	Surname       string `json:"surname" db:"surname"`
	Othernames    string `json:"othernames" db:"othernames"`
	BVN           string `json:"bvn" db:"bvn"`
	Country       string `json:"country"`
	BankCode      string `json:"bank_code" db:"bank_code"`
	AccountNumber string `json:"account_number" db:"account_number"`
	UserID        string `json:"user_id" db:"user_id"`
	StoreName     string `json:"store_name" db:"store_name"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *CreateUserReq) (*User, error)     // Create a new user
	GetUsers(ctx context.Context) ([]*User, error)                          // Create a new user
	GetUser(ctx context.Context, filter string) (*User, error)              // Create a new user
	GetUserByEmailOrPhone(ctx context.Context, email string) (*User, error) // Get user information by email
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)      // Perform user login
	VerifyOTP(ctx context.Context, req *User) (*User, error)
	ToggleStoreFollowStatus(ctx context.Context, userId, storeId uint32) error
	UpdateUser(ctx context.Context, user *User) (*User, error)
	CreateDVAAccount(ctx context.Context, req *DVADetails) (string, error)
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) // Create a new user
	GetUsers(ctx context.Context) ([]*User, error)
	GetUser(ctx context.Context, filter string) (*User, error) // Create a new user
	VerifyOTP(ctx context.Context, req *User) (*User, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) // Perform user login
	UpdateUser(ctx context.Context, user *User) (*User, error)
	ToggleStoreFollowStatus(ctx context.Context, userId, storeId uint32) error
	CreateDVAAccount(ctx context.Context, req *DVADetails) (string, error)
}
