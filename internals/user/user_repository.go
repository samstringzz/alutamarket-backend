package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Chrisentech/aluta-market-api/utils"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRespository(db DBTX) Respository {
	return &repository{db: db}
}

func (r *repository) GetUserByEmailOrPhone(ctx context.Context, identifier string) (*User, error) {
	u := User{}
	query := "SELECT id, email, password, campus, phone, usertype, fullname FROM users WHERE email = $1 OR phone = $2"
	err := r.db.QueryRowContext(ctx, query, identifier).Scan(&u.ID, &u.Email, &u.Password, &u.Campus, &u.Phone, &u.Usertype, &u.Fullname)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *repository) CreateUser(ctx context.Context, req *User) (*User, error) {
	var lastInsertedId int
	u := User{}
	otpCode := utils.GenerateOTP()
	_, err2 := utils.SendOtpMessage(otpCode, req.Phone)
	if err2 != nil {
		return nil, err2
	}

	isUserEmail, err3 := r.GetUserByEmailOrPhone(ctx, req.Email)
	if err3 != nil {
		return nil, err3
	}
	if *isUserEmail != u {
		return nil, errors.New("User already exist")
	}
	isUserPhone, err4 := r.GetUserByEmailOrPhone(ctx, req.Phone)
	if err4 != nil {
		return nil, err4
	}
	if *isUserPhone != u {
		return nil, errors.New("User already exist")
	}
	query := "INSERT INTO users(campus, email, password, fullname, phone, usertype, active,twofa, wallet,code)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id"
	err := r.db.QueryRowContext(ctx, query, req.Campus, req.Email, req.Password, req.Fullname, req.Phone, req.Usertype, false, false, 0, otpCode).Scan(&lastInsertedId)
	if err != nil {
		return nil, err
	}

	user := &User{
		Phone: req.Phone,
	}
	return user, nil
}
