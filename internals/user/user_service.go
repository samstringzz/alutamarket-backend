package user

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/Chrisentech/aluta-market-api/utils"
	"github.com/golang-jwt/jwt/v4"
)

type service struct {
	Respository
	timeout time.Duration
}

func NewService(repository Respository) Service {
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
	u := &User{
		Email:    req.Email,
		Password: hashedPassword,
		Campus:   req.Campus,
		Fullname: req.Fullname,
		Phone:    req.Phone,
		Usertype: req.Usertype,
	}
	r, err := s.Respository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	res := &CreateUserRes{
		Message: fmt.Sprintf("Registration successful.Verify the OTP sent to %s", r.Phone),
		Status:  http.StatusOK,
		Data : r,
	}
	return res, nil
}
func (s *service) GetUsers(c context.Context) ([]*User, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()
	r, err := s.Respository.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.Respository.GetUserByEmailOrPhone(ctx, req.Email)
	if err != nil {
		return &LoginUserRes{}, err
	}
	err = utils.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &LoginUserRes{}, err

	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Fullname: u.Fullname,
		Usertype: u.Usertype,
		Phone:    u.Phone,
		Campus:   u.Campus,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return &LoginUserRes{}, err
	}
	return &LoginUserRes{accessToken: ss, ID: strconv.Itoa(int(u.ID))}, nil
}
