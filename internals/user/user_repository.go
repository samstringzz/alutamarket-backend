package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Chrisentech/aluta-market-api/errors"
	"github.com/Chrisentech/aluta-market-api/internals/store"
	"github.com/Chrisentech/aluta-market-api/services"
	"github.com/Chrisentech/aluta-market-api/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
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

func (r *repository) resendOTP(ctx context.Context, phone string) error {
	otpCode := utils.GenerateOTP()
	codeExpiry := time.Now().Add(5 * time.Minute) //An expiry time of 5min
	user, _ := r.GetUserByEmailOrPhone(ctx, phone)

	user.Code = otpCode
	user.Codeexpiry = codeExpiry
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	// Define the template string with a placeholder for the passcode
	messageTemplate := "Hello Comrade, your Alutamarket VIP passcode is: %s. Make haste, the party is waiting."
	// Remove the plus sign from the phone number
	phoneWithoutPlus := strings.TrimPrefix(phone, "+")
	// Insert the dynamic data into the template string
	message := fmt.Sprintf(messageTemplate, otpCode)
	_, err := services.SendSMS(phoneWithoutPlus, "N-Alert", message)
	if err != nil {
		return err
	}
	return nil
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
func (r *repository) CreateStore(ctx context.Context, req *store.Store) (*store.Store, error) {

	resp := &store.Store{
		Name:               req.Name,
		Email:              req.Email,
		Link:               req.Link,
		UserID:             req.UserID,
		Description:        req.Description,
		HasPhysicalAddress: req.HasPhysicalAddress,
		Address:            req.Address,
		Wallet:             0,
		Status:             true,
		Phone:              req.Phone,
	}
	if err := r.db.Create(resp).Error; err != nil {
		r.db.Rollback()
		return nil, err
	}
	// user: &User{
	// 	// Fullname: req.,
	// }
	//Create DVA for seller link for user
	// _, err := r.CreateDVAAccount(ctx, &DVADetails{UserEmail: req.Email, StoreName: req.Name, StoreEmail: req.Email})
	// if err != nil {
	// 	return nil, err
	// }
	return resp, nil
}

func boolPtr(b bool) *bool {
	return &b
}

func (r *repository) CreateUser(ctx context.Context, req *CreateUserReq) (*User, error) {
	// Start a new database transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		log.Printf("Failed to start transaction: %v", tx.Error)
		return nil, tx.Error
	}

	// Defer a function to handle transaction rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
			tx.Rollback()
		}
	}()

	otpCode := utils.GenerateOTP()
	fmt.Printf("The generated OTP is %s\n", otpCode)

	var count int64
	codeExpiry := time.Now().Add(5 * time.Minute) // An expiry time of 5 minutes
	tx.Model(&User{}).Where("email = ? OR phone = ?", req.Email, req.Phone).Count(&count)
	if count > 0 {
		log.Printf("User already exists: email=%s, phone=%s", req.Email, req.Phone)
		tx.Rollback()
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "User already exists")
	}

	err := utils.AddEmailSubscriber(req.Email)
	if err != nil {
		log.Printf("Failed to add email subscriber: %v", err)
		return nil, err
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
		Code:       otpCode,
		Codeexpiry: codeExpiry,
		Avatar:     "https://icon-library.com/images/anonymous-avatar-icon/anonymous-avatar-icon-25.jpg",
	}
	if err := tx.Create(newUser).Error; err != nil {
		log.Printf("Failed to create new user: %v", err)
		tx.Rollback()
		return nil, err
	}

	// Email credentials
	to := []string{req.Email}
	contents := map[string]string{
		"otp_code": otpCode,
	}

	// Define the template string with a placeholder for the passcode
	messageTemplate := "Hello Comrade, your Alutamarket VIP passcode is: %s. Make haste, the party is waiting."
	// Remove the plus sign from the phone number
	phoneWithoutPlus := strings.TrimPrefix(req.Phone, "+")
	// Insert the dynamic data into the template string
	message := fmt.Sprintf(messageTemplate, otpCode)
	_, err = services.SendSMS(phoneWithoutPlus, "N-Alert", message)
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
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

		user := &User{
			Fullname: req.Fullname,
			Email:    req.Email,
			Phone:    req.Phone,
		}

		// Create DVA for seller link for user
		_, err := r.CreateDVAAccount(ctx, &DVADetails{User: *user, StoreName: req.StoreName, StoreEmail: req.StoreEmail})
		if err != nil {
			log.Printf("Failed to create DVA account: %v", err)
			return nil, err
		}

		if err := tx.Create(createdStore).Error; err != nil {
			log.Printf("Failed to create store: %v", err)
			tx.Rollback()
			return nil, err
		}

		templateID := "7178d0b2-a957-410d-b24d-e812252451da"
		err = services.SendEmail(templateID, "Welcome to Alutamarket🎉", to, contents)
		if err != nil {
			log.Printf("Failed to send welcome email to seller: %v", err)
			return nil, err
		}
	} else {
		templateID := "633d65f7-0545-4550-9983-8b309afa3d03"
		err := services.SendEmail(templateID, "Welcome to Alutamarket🎉", to, contents)
		if err != nil {
			log.Printf("Failed to send welcome email: %v", err)
			return nil, err
		}
	}

	// Commit the transaction if everything succeeded
	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		tx.Rollback()
		return nil, err
	}

	log.Printf("Successfully created user: ID=%d", newUser.ID)
	return newUser, nil
}

func getEmail(storeEmail, userEmail string) string {
	if storeEmail != "" {
		return storeEmail
	}
	return userEmail
}

func (r *repository) CreateDVAAccount(ctx context.Context, req *DVADetails) (string, error) {

	// Create dedicated account
	dedicatedAccountURL := "https://api.paystack.co/dedicated_account/assign"
	method := "POST"
	names := strings.Split(req.User.Fullname, " ")
	if len(names) < 2 {
		return "", fmt.Errorf("invalid user name")
	}

	payload := map[string]interface{}{
		"email":          getEmail(req.StoreEmail, req.User.Email),
		"first_name":     names[0],
		"middle_name":    names[1],
		"last_name":      req.StoreName,
		"phone":          req.User.Phone,
		"preferred_bank": "wema",
		"country":        "NG",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	newReq, err := http.NewRequest(method, dedicatedAccountURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	newReq.Header.Add("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	newReq.Header.Add("Content-Type", "application/json")

	res, err := client.Do(newReq)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("error: paystack dedicated account creation failed with status %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	var dedicatedAccountResp map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&dedicatedAccountResp)
	if err != nil {
		return "", fmt.Errorf("error decoding paystack dedicated account response: %w", err)
	}

	fmt.Println("Paystack Dedicated Account Response:", dedicatedAccountResp)

	// Return the response as a JSON string
	jsonString, err := json.Marshal(dedicatedAccountResp)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}

func (r *repository) VerifyOTP(ctx context.Context, req *VerifyOTPReq) (*LoginUserRes, error) {
	foundUser, _ := r.GetUserByEmailOrPhone(ctx, req.Phone)
	err := r.db.Where("phone = ?", req.Phone).First(foundUser).Error
	fmt.Print(req.Attempts)
	if err != nil {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist")
	}

	// If the code is incorrect, increment the req.Attempts and send a new code if the req.Attempts is greater than 3.
	if req.Code != foundUser.Code {
		if req.Attempts > 3 {
			// Send a new code here.
			r.resendOTP(ctx, foundUser.Phone)
			return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "New code has been sent")
		} else {
			return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "Invalid code!!")
		}
	}
	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(foundUser.ID)),
		Fullname: foundUser.Fullname,
		Usertype: foundUser.Usertype,
		Campus:   foundUser.Campus,
		Phone:    foundUser.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(foundUser.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Example: Refresh token expires in 7 days
		},
	})
	refreshSS, err := refreshClaims.SignedString([]byte(os.Getenv("REFRESH_SECRET_KEY")))
	if err != nil {
		return nil, err
	}

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(foundUser.ID)),
		Fullname: foundUser.Fullname,
		Campus:   foundUser.Campus,
		Usertype: foundUser.Usertype,
		Phone:    foundUser.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(foundUser.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	accessSS, err := accessClaims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}
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

	if foundUser.Codeexpiry.Before(time.Now()) {
		return nil, errors.NewAppError(http.StatusConflict, "BAD REQUEST", "OTP Expired!!")
	}
	// If the req.Attempts is less than or equal to 3 and the code is correct, verify the user.
	if req.Attempts <= 3 && req.Code == foundUser.Code {
		if foundUser.ID == 0 {
			return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist in the database")
		}
		if *foundUser.Active {
			return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "User account is already verified!")
		}
		if err := r.db.Model(foundUser).Update("active", true).Error; err != nil {
			return nil, err
		}
		return &LoginUserRes{AccessToken: accessSS, RefreshToken: refreshSS, ID: foundUser.ID}, nil
	}

	// If the code is incorrect and the counter is less than or equal to 3, return an error.
	return &LoginUserRes{AccessToken: accessSS, RefreshToken: refreshSS, ID: foundUser.ID}, nil

}
func (r *repository) Login(ctx context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	var user User

	if err := r.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist")
	}
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid Credentials")
	}
	if err := r.db.Where("active = ?", true).First(&user).Error; err != nil {
		r.resendOTP(ctx, user.Phone)
		return nil, errors.NewAppError(http.StatusExpectationFailed, user.Phone, "Your account is suspended/not verified")
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
		return nil, errors.NewAppError(http.StatusProxyAuthRequired, "ACTION REQUIRED", "This account is 2-FA protected,enter Otp to continue")
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
	if req.PaymentDetails.Address != "" {
		existingUser.PaymentDetails = req.PaymentDetails
	}

	// Update the User in the repository
	err = r.db.Save(existingUser).Error
	if err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (r *repository) GetMyDVA(ctx context.Context, userEmail string) (*Account, error) {

	dedicatedAccountURL := "https://api.paystack.co/dedicated_account"
	method := "GET"
	client := &http.Client{}
	newReq, err := http.NewRequest(method, dedicatedAccountURL, nil)
	if err != nil {
		return nil, err
	}
	newReq.Header.Add("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	newReq.Header.Add("Content-Type", "application/json")
	res, err := client.Do(newReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("error: paystack dedicated account creation failed with status %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("error: paystack dedicated account retrieval failed with status %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	var dedicatedAccountResp struct {
		Status  bool       `json:"status"`
		Message string     `json:"message"`
		Data    []*Account `json:"data"`
	}

	err = json.NewDecoder(res.Body).Decode(&dedicatedAccountResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding paystack dedicated account response: %w", err)
	}

	if !dedicatedAccountResp.Status {
		return nil, fmt.Errorf("error: paystack dedicated account retrieval failed with message: %s", dedicatedAccountResp.Message)
	}

	for _, account := range dedicatedAccountResp.Data {
		if account.Customer.Email == userEmail {
			return account, nil
		}
	}

	return nil, fmt.Errorf("error: no account found for user with email %s", userEmail)
}

func (r *repository) SetPaymentDetais(ctx context.Context, req *PaymentDetails, userId uint32) error {
	existingUser, err := r.GetUser(ctx, strconv.FormatUint(uint64(userId), 10))
	if err != nil {
		return err
	}
	detail := &PaymentDetails{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Info:    req.Info,
	}
	existingUser.PaymentDetails = *detail

	// Update the Store in the repository
	err = r.db.Save(existingUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SendPasswordResetLink(ctx context.Context, req *PasswordReset) error {

	panic("hello")
}
