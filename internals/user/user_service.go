package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Chrisentech/aluta-market-api/utils"
	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(5) * time.Second,
	}
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Fullname string `json:"fullname"`
	Campus   string `json:"campus"`
	Phone    string `json:"phone"`
	Usertype string `json:"usertype"`
	jwt.RegisteredClaims
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := utils.HashPasswword(req.Password)
	if err != nil {
		return nil, err
	}
	u := &CreateUserReq{
		Email:              req.Email,
		Password:           hashedPassword,
		Campus:             req.Campus,
		Fullname:           req.Fullname,
		Phone:              req.Phone,
		Usertype:           req.Usertype,
		StoreName:          req.StoreName,
		StoreAddress:       req.StoreAddress,
		StoreLink:          req.StoreLink,
		Description:        req.Description,
		HasPhysicalAddress: req.HasPhysicalAddress,
	}
	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	res := &CreateUserRes{
		Message: fmt.Sprintf("Registration successful.Verify the OTP sent to %s", r.Phone),
		Status:  http.StatusOK,
		Data:    r,
	}
	return res, nil
}
func (s *service) VerifyOTP(c context.Context, req *User) (*User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	u := &User{
		Phone: req.Phone,
		Email: req.Email,
		Code:  req.Code,
	}
	r, err := s.Repository.VerifyOTP(ctx, u)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (s *service) GetUsers(c context.Context) ([]*User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, err := s.Repository.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *service) GetUser(c context.Context, filter string) (*User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, err := s.Repository.GetUser(ctx, filter)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Repository.Login(ctx, req)
	if err != nil {
		return &LoginUserRes{}, err
	}
	return u, nil
}

func (s *service) ToggleStoreFollowStatus(c context.Context, userId, storeId uint32) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	err := s.Repository.ToggleStoreFollowStatus(ctx, userId, storeId)
	if err != nil {
		// return &LoginUserRes{}, err
		return err
	}
	return nil
}
func (s *service) UpdateUser(ctx context.Context, user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	usr, err := s.Repository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
