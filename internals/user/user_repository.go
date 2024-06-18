package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db         *gorm.DB
	otpCounter uint8
}
type accessTokenCookieKey struct{}

// SetAccessTokenCookie sets the access token cookie in the context.
func SetAccessTokenCookie(ctx context.Context, cookie *http.Cookie) context.Context {
	return context.WithValue(ctx, accessTokenCookieKey{}, cookie)
}

// GetAccessTokenCookie retrieves the access token cookie from the context.
func GetAccessTokenCookie(ctx context.Context) *http.Cookie {
	cookie, _ := ctx.Value(accessTokenCookieKey{}).(*http.Cookie)
	return cookie
}
func NewRepository() Repository {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbURI := os.Getenv("DB_URI")

	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

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
func boolPtr(b bool) *bool {
	return &b
}
func (r *repository) CreateUser(ctx context.Context, req *CreateUserReq) (*User, error) {
	// Start a new database transaction

	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Defer a function to handle transaction rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	otpCode := utils.GenerateOTP()
	fmt.Printf("The generatedOtp is %s", otpCode)

	// utils.SendOTPMessage(req.Phone,otpCode)

	var count int64
	codeExpiry := time.Now().Add(5 * time.Minute) //An expiry time of 5min
	tx.Model(&User{}).Where("email = ? OR phone = ?", req.Email, req.Phone).Count(&count)
	if count > 0 {
		tx.Rollback()
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "User already exists")
	}

	newUser := &User{
		Campus:     req.Campus,
		Email:      req.Email,
		Password:   req.Password,
		Fullname:   req.Fullname,
		Phone:      req.Phone,
		Usertype:   req.Usertype,
		Active:     boolPtr(false),
		Twofa:      boolPtr(false),
		Code:       "12345",
		Codeexpiry: codeExpiry,
		Avatar:     "https://icon-library.com/images/anonymous-avatar-icon/anonymous-avatar-icon-25.jpg",
	}

	if err := tx.Create(newUser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if req.Usertype == "seller" {
		createdStore := &store.Store{
			Name:               req.StoreName,
			Link:               req.StoreLink,
			UserID:             newUser.ID,
			Description:        req.Description,
			HasPhysicalAddress: req.HasPhysicalAddress,
			Address:            req.StoreAddress,
			Wallet:             0,
			Status:             true,
			Phone:              req.StorePhone,
		}

		if err := tx.Create(createdStore).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit the transaction if everything succeeded
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return newUser, nil
}

func (r *repository) CreateDVAAccount(ctx context.Context, req *DVADetails) (string, error) {
	user, _ := r.GetUser(ctx, req.UserID)

	//create a user for paystack account
	paystackCustomerURL := "https://api.paystack.co/customer"
	method := "POST"
	names := strings.Split(user.Fullname, " ")
	if len(names) <= 0 {
		return "", errors.NewAppError(505, "Internal Server Error", "There was an error creating paystack account")
	}
	payload := map[string]interface{}{
		"email":      user.Email,
		"first_name": names[0],
		"last_name":  req.StoreName,
		"phone":      user.Phone,
	}
	jsonPayload, _ := json.Marshal(payload)

	client := &http.Client{}
	newReq, err := http.NewRequest(method, paystackCustomerURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println(err)
		// return
	}
	newReq.Header.Add("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	newReq.Header.Add("Content-Type", "application/json")

	res, err := client.Do(newReq)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	var customerResp map[string]interface{}
	json.NewDecoder(res.Body).Decode(&customerResp)
	fmt.Println(customerResp)

	//We should get a customerResp.data.code f

	// verify the customer before creating his DVA account
	// Accessing nested fields within customerResp
	customerCode, ok := customerResp["data"].(map[string]interface{})["customer_code"].(string)
	if !ok {
		fmt.Println("Error: unable to extract customer code")

	}

	// verify the customer before creating his DVA account
	paystackCustomerValidationURL := "https://api.paystack.co/customer/" + customerCode + "/identification"
	method = "POST"
	payload = map[string]interface{}{
		"country":        "NG",
		"type":           "bank_account",
		"account_number": req.AccountNumber,
		"bvn":            req.BVN,
		"bank_code":      "007",
		"first_name":     "Asta",
		"last_name":      "Lavista",
	}
	jsonPayload, _ = json.Marshal(payload)

	client = &http.Client{}

	newReq2, err := http.NewRequest(method, paystackCustomerValidationURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println(err)
		// return
	}
	newReq2.Header.Add("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	newReq2.Header.Add("Content-Type", "application/json")

	res2, err := client.Do(newReq2)
	if err != nil {
		fmt.Println(err)
		// return
	}
	defer res2.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)
	fmt.Println(result)
	return "Successfully Created an account", nil
}

func (r *repository) VerifyOTP(ctx context.Context, req *User) (*User, error) {
	foundUser := &User{}
	err := r.db.Where("phone = ?", req.Phone).First(foundUser).Error

	if err != nil {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist")
	}

	// If the code is incorrect, increment the r.otpCounter and send a new code if the r.otpCounter is greater than 3.
	if req.Code != "12345" {
		r.otpCounter++
		if r.otpCounter > 3 {
			// Send a new code here.
			return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "New code has been sent")
		}
	}
	if foundUser.Codeexpiry.Before(time.Now()) {
		return nil, errors.NewAppError(http.StatusConflict, "BAD REQUEST", "OTP Expired!!")
	}
	// If the r.otpCounter is less than or equal to 3 and the code is correct, verify the user.
	if r.otpCounter <= 3 && req.Code == "12345" {
		if foundUser.ID == 0 {
			return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist in the database")
		}
		if *foundUser.Active {
			return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "User account is already verified!")
		}
		if err := r.db.Model(foundUser).Update("active", true).Error; err != nil {
			return nil, err
		}
		return foundUser, nil
	}

	// If the code is incorrect and the counter is less than or equal to 3, return an error.
	return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Incorrect Otp Provided")
}
func (r *repository) Login(ctx context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	var user User

	if err := r.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist")
	}
	if err := r.db.Where("active = ?", true).First(&user).Error; err != nil {
		return nil, errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Your account is suspended/not verified")
	}
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid Credentials")
	}

	// Generate a new refresh token
	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Fullname: user.Fullname,
		Usertype: user.Usertype,
		Campus:   user.Campus,
		Phone:    user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Example: Refresh token expires in 7 days
		},
	})
	refreshSS, err := refreshClaims.SignedString([]byte(os.Getenv("REFRESH_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	// Store the refresh token in the database (you may need to add a field for this)
	if err := r.db.Model(&user).Update("refresh_token", refreshSS).Error; err != nil {
		return nil, err
	}

	// Generate the access token
	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Fullname: user.Fullname,
		Campus:   user.Campus,
		Usertype: user.Usertype,
		Phone:    user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	accessSS, err := accessClaims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	r.db.Model(&user).Updates(User{RefreshToken: refreshSS, AccessToken: accessSS})
	if *user.Twofa {
		//send otp
		otpCode := utils.GenerateOTP()
		r.db.Model(&user).Update("code", otpCode)
		return nil, errors.NewAppError(http.StatusCreated, "ACTION REQUIRED", "This account is 2-FA protected,enter Otp to continue")
	}

	// if err != nil {
	// 	return nil, err
	// }

	// Set a cookie for the access token with an expiration time matching the token's expiration
	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessSS,
		Expires:  time.Now().Add(24 * time.Hour), // Set the expiration to match the token's expiration
		HttpOnly: true,
		Secure:   false, // Set to true if your server uses HTTPS
		SameSite: http.SameSiteStrictMode,
	}

	// Add the access token cookie to the context
	SetAccessTokenCookie(ctx, &accessCookie)

	return &LoginUserRes{AccessToken: accessSS, RefreshToken: refreshSS, ID: user.ID}, nil
}

func (r *repository) GetUsers(ctx context.Context) ([]*User, error) {
	var users []*User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetUser(ctx context.Context, filter string) (*User, error) {
	var user User
	// query := r.db.Where("active = true")
	query := r.db.Where("id = ?", filter)

	if err := query.First(&user).Error; err != nil {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist")
	}
	return &user, nil
}

func (r *repository) TwoFa(ctx context.Context, req *User) (bool, error) {
	var user User
	if err := r.db.Where("email = ? OR phone = ?", req.Email, req.Phone).First(&user).Error; err != nil {
		return false, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist")
	}
	r.db.Model(&user).Update("two_fa", true)
	return *user.Twofa, nil
}

func (r *repository) ToggleStoreFollowStatus(ctx context.Context, userId, storeId uint32) error {
	// Retrieve the store with the given storeId
	foundStore := &store.Store{}
	if err := r.db.First(foundStore, storeId).Error; err != nil {
		return err
	}
	// Convert userId to a string
	userIdStr := strconv.FormatUint(uint64(userId), 10)
	// Retrieve the user who wants to follow/unfollow the store

	// Retrieve the user using the string representation of userId
	foundUser, err := r.GetUser(ctx, userIdStr)
	if err != nil {
		return err
	}

	// Check if the user is already following the store
	isFollowing := false
	for _, follower := range foundStore.Followers {
		if follower.FollowerID == userId {
			isFollowing = true
			break
		}
	}
	// Toggle the follow status
	if isFollowing {
		// If already following, unfollow
		newFollowers := make([]store.Follower, 0)
		for _, follower := range foundStore.Followers {
			if follower.FollowerID != userId {
				newFollowers = append(newFollowers, follower)
			}
		}
		foundStore.Followers = newFollowers
	} else {
		// If not following, follow
		foundStore.Followers = append(foundStore.Followers, store.Follower{
			FollowerID:    userId,
			FollowerName:  foundUser.Fullname,
			FollowerImage: foundUser.Avatar,
		})
	}

	// Save the updated store with GORM
	if err := r.db.Save(foundStore).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateUser(ctx context.Context, req *User) (*User, error) {

	// First, check if the Store exists by its ID or another unique identifier
	existingUser, err := r.GetUser(ctx, strconv.FormatUint(uint64(req.ID), 10))
	if err != nil {
		return nil, err
	}

	// Update only the fields that are present in the req
	if req.Fullname != "" {
		existingUser.Fullname = req.Fullname
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Avatar != "" {
		existingUser.Avatar = req.Avatar
	}
	if req.Phone != "" {
		existingUser.Phone = req.Phone
	}
	if req.Usertype != "" {
		existingUser.Usertype = req.Usertype
	}
	if req.Gender != "" {
		existingUser.Gender = req.Gender
	}
	if req.Dob != "" {
		existingUser.Dob = req.Dob
	}
	if req.AccessToken != "" {
		existingUser.AccessToken = req.AccessToken
	}
	if req.RefreshToken != "" {
		existingUser.RefreshToken = req.RefreshToken
	}
	if req.Active != existingUser.Active && req.Active != nil {
		existingUser.Active = req.Active
	}

	if req.Twofa != existingUser.Twofa && req.Twofa != nil {
		existingUser.Twofa = req.Twofa
	}

	// Update the Store in the repository
	err = r.db.Save(existingUser).Error
	if err != nil {
		return nil, err
	}

	return existingUser, nil
}
