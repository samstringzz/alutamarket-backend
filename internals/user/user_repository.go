package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
    "errors"
	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
    db *gorm.DB
}

func NewRepository(dbConnStr string) Respository {
    // Initialize the database connection
    db, err := gorm.Open(postgres.Open(dbConnStr), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto-migrate the User model to create the table if not exists
    db.AutoMigrate(&User{})

    return &repository{
        db: db,
    }
}

func (r *repository) GetUserByEmailOrPhone(ctx context.Context, identifier string) (*User, error) {
    u := User{}
    err := r.db.Where("email = ? OR phone = ?", identifier, identifier).First(&u).Error
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *repository) CreateUser(ctx context.Context, req *User) (*User, error) {
    otpCode := utils.GenerateOTP()
    fmt.Printf("The generatedOtp is%s",otpCode)
    // _, err := utils.SendOtpMessage(otpCode, req.Phone)
    // if err != nil {
    //     return nil, err
    // }

    var count int64
    codeExpiry := time.Now().Add(5 * time.Minute)
    r.db.Model(&User{}).Where("email = ? OR phone = ?", req.Email, req.Phone).Count(&count)
    if count > 0 {
        return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "User already exists")
    }

    newUser := &User{
        Campus:   req.Campus,
        Email:    req.Email,
        Password: req.Password,
        Fullname: req.Fullname,
        Phone:    req.Phone,
        Usertype: req.Usertype,
        Active:   false,
        Twofa:    false,
        Wallet:   0,
        Code:     "12345",
        CodeExpiry: codeExpiry,
    }
    if err := r.db.Create(newUser).Error; err != nil {
        return nil, err
    }
    return newUser, nil
}

func (r *repository) GetUsers(ctx context.Context) ([]*User, error) {
    var users []*User
    if err := r.db.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

func (r *repository) GetUser(ctx context.Context, filter string, filterOptions string) (*User, error) {
    var user User
    query := r.db.Where("active = true")

    switch filter {
    case "id":
        query = query.Where("id = ?", filterOptions)
    case "email":
        query = query.Where("email = ?", filterOptions)
    case "phone":
        query = query.Where("phone = ?", filterOptions)
    default:
        return nil, errors.New("Invalid attribute type")
    }

    if err := query.First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *repository) VerifyOTP(ctx context.Context, req *User) (*User, error) {
    foundUser := User{}
    err := r.db.Where("phone = ?", req.Phone).First(&foundUser).Error
    if err != nil {
        return nil, err
    }
    if foundUser.ID == 0 {
        return nil, errors.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "User does not exist in the database")
    }
    return &User{
        Active: req.Active,
    }, nil
}
