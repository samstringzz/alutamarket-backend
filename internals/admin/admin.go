package admin

import (
	"context"
)

type Admin struct {
	ID          uint8    `gorm:"primaryKey;uniqueIndex;not null;autoIncrement"  json:"id" db:"id"`
	Fullname    string   `json:"fullname" db:"fullname"`
	Email       string   `json:"email" db:"email"`
	Permissions []string `gorm:"serializer:json" json:"permissions" db:"permissions"`
	Password    string   `gorm:"serializer:json" json:"password" db:"password"`
}
type LoginAdminRes struct {
	AccessToken  string // Access token generated for the logged-in user
	RefreshToken string // RefreshToken token generated for the logged-in user
	ID           uint32 `json:"id" db:"id"` // Unique identifier for the logged-in user
}
type LoginAdminReq struct {
	Password string `json:"password" db:"password"` // Password of the user for login
	Email    string `json:"email" db:"email"`       // Email address of the user for login
}
type Repository interface {
	Login(ctx context.Context, input *LoginAdminReq) (*LoginAdminRes, error)
	CreateAdmin(ctx context.Context, input *Admin) (*LoginAdminRes, error)
	GetAdmin(ctx context.Context, id uint8) (*Admin, error)
	GetAdmins(ctx context.Context) ([]*Admin, error)
	ApproveProduct()
	MessageUser()

}

type Service interface {
	Login(ctx context.Context, input *LoginAdminReq) (*LoginAdminRes, error)
	CreateAdmin(ctx context.Context, input *Admin) (*LoginAdminRes, error)
	GetAdmin(ctx context.Context, id uint8) (*Admin, error)
	GetAdmins(ctx context.Context) ([]*Admin, error)

}
