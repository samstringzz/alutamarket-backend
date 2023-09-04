package user

import (
	"context"
	"time"
)

type User struct {
	ID            int64  `json:"id" db:"id"`                   // Unique identifier for the user
	Campus        string `json:"campus" db:"campus"`           // Campus of the user
	Email         string `json:"email" db:"email"`             // Email address of the user
	Password      string `json:"password" db:"password"`       // Password of the user
	Fullname      string `json:"fullname" db:"fullname"`       // Full name of the user
	Phone         string `json:"phone" db:"phone"`             // Phone number of the user
	Usertype      string `json:"usertype" db:"usertype"`       // Type of user (e.g., seller,buyer,admin)
	Active        bool 	`json:"active" db:"active"`  		   // For accesibility of user,
	Twofa		  bool  `json:"twofa" db:"twofa"`              // Two factor authentication
	Wallet		  int64  `json:"wallet,omitempty" db:"wallet"` // Balance of the user's wallet (only for seller)
	Code		  string `json:"code" db:"code"`			   // otp code for verifications
	CodeExpiry    time.Time `json:"codeexpiry" db:"codeexpiry"`	// Expiry time for otpCode
}

type CreateUserReq struct {
	Campus        string `json:"campus" db:"campus"`           // Campus of the user
	Email         string `json:"email" db:"email"`             // Email address of the user
	Password      string `json:"password" db:"password"`       // Password of the user
	Fullname      string `json:"fullname" db:"fullname"`       // Full name of the user
	Phone         string `json:"phone" db:"phone"`             // Phone number of the user
	Usertype      string `json:"usertype" db:"usertype"`       // Type of user (e.g., seller,buyer,admin)
	Active        bool 	`json:"active" db:"active"`  		   // For accesibility of user,
	Twofa		  bool  `json:"twofa" db:"twofa"`              // Two factor authentication
	Wallet		  int64  `json:"wallet,omitempty" db:"wallet"` // Balance of the user's wallet (only for seller)
	Code		  string `json:"code,omitempty" db:"code"`			   // otp code for verifications
	CodeExpiry    time.Time `json:"codeexpiry,omitempty" db:"codeexpiry"`	// Expiry time for otpCode
}

type CreateUserRes struct {
	Message string
	Status  int
	Data    interface {}
}

type LoginUserReq struct {
	Password string `json:"password" db:"password"` // Password of the user for login
	Email    string `json:"email" db:"email"`       // Email address of the user for login
}

type LoginUserRes struct {
	accessToken string // Access token generated for the logged-in user
	ID          string `json:"id" db:"id"` // Unique identifier for the logged-in user
}

type Respository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)              // Create a new user
	GetUserByEmailOrPhone(ctx context.Context, email string) (*User, error) // Get user information by email
}

type Service interface {
	CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) // Create a new user
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)        // Perform user login
}
