package user

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/samstringzz/alutamarket-backend/database"
	"github.com/samstringzz/alutamarket-backend/errors"
	"github.com/samstringzz/alutamarket-backend/internals/models"
	"github.com/samstringzz/alutamarket-backend/utils"
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
	_, err := utils.SendSMS(phoneWithoutPlus, "N-Alert", message)
	if err != nil {
		return err
	}
	return nil
}
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(), // Use the database manager
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
func (r *repository) CreateStore(ctx context.Context, req *models.Store) (*models.Store, error) {
	// Convert userId to a string
	userIdStr := strconv.FormatUint(uint64(req.UserID), 10)
	foundUser, _ := r.GetUser(ctx, userIdStr)
	resp := &models.Store{
		Name:               req.Name,
		Email:              req.Email,
		Link:               req.Link,
		UserID:             req.UserID,
		Description:        req.Description,
		HasPhysicalAddress: req.HasPhysicalAddress,
		Address:            req.Address,
		Wallet:             0,
		Status:             true,
		Background:         "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQQLbvWGTFQh6OGWPfkLx2xBS_OP3oZJzQubA&s",
		Phone:              req.Phone,
	}
	if err := r.db.Create(resp).Error; err != nil {
		r.db.Rollback()
		return nil, err
	}

	// Create DVA for seller link for user
	_, err := r.CreateDVAAccount(ctx, &DVADetails{User: *foundUser, StoreName: req.Name, StoreEmail: req.Email})
	if err != nil {
		return nil, err
	}
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
	codeExpiry := time.Now().Add(5 * time.Minute)
	tx.Model(&User{}).Where("email = ? OR phone = ?", req.Email, req.Phone).Count(&count)
	if count > 0 {
		log.Printf("User already exists: email=%s, phone=%s", req.Email, req.Phone)
		tx.Rollback()
		return nil, errors.NewAppError(http.StatusConflict, "CONFLICT", "User already exists")
	}

	// Make email subscription optional
	_ = utils.AddEmailSubscriber(req.Email) // Ignore any subscription errors

	// Generate UUID for the new user
	uuid := utils.GenerateUUID()

	defaultDob := time.Now().Format("2006-01-02")
	newUser := &User{
		UUID:       uuid,
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
		Dob:        defaultDob,
		Gender:     "unspecified",
	}
	// Create user first
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

	// Try to send SMS, but don't fail if SMS fails
	messageTemplate := "Hello Comrade, your Alutamarket VIP passcode is: %s. Make haste, the party is waiting."
	phoneWithoutPlus := strings.TrimPrefix(req.Phone, "+")
	message := fmt.Sprintf(messageTemplate, otpCode)

	if _, err := utils.SendSMS(phoneWithoutPlus, "N-Alert", message); err != nil {
		// Log the error but continue with user creation
		log.Printf("Warning: SMS sending failed: %v. Continuing with registration...", err)
	}

	// Continue with email sending regardless of SMS status
	if req.Usertype == "seller" {
		createdStore := &models.Store{
			Name:               req.StoreName,
			Link:               req.StoreLink,
			UserID:             newUser.ID,
			Description:        req.Description,
			HasPhysicalAddress: req.HasPhysicalAddress,
			Address:            req.StoreAddress,
			Wallet:             0,
			Status:             true,
			Phone:              req.StorePhone,
			Background:         "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQQLbvWGTFQh6OGWPfkLx2xBS_OP3oZJzQubA&s",
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
		err = utils.SendEmail(templateID, "Welcome to Alutamarket🎉", to, contents)
		if err != nil {
			log.Printf("Failed to send welcome email to seller: %v", err)
			return nil, err
		}
	} else {
		templateID := "633d65f7-0545-4550-9983-8b309afa3d03"
		if err := utils.SendEmail(templateID, "Welcome to Alutamarket🎉", to, contents); err != nil {
			log.Printf("Warning: Email sending failed: %v", err)
			// Continue with registration even if email fails
		}
	}

	// Commit the transaction
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
	parts := strings.Fields(req.StoreName)
	firstName := parts[0]
	lastName := ""
	if len(parts) > 1 {
		lastName = strings.Join(parts[1:], " ")
	}
	payload := map[string]interface{}{
		"email":          getEmail(req.StoreEmail, req.User.Email),
		"first_name":     firstName,
		"last_name":      lastName,
		"phone":          req.User.Phone,
		"preferred_bank": "wema-bank",
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

	// First find the user by email
	if err := r.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist")
	}

	// Check password
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, errors.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid Credentials")
	}

	// Check if user is active - Remove the additional query and just check the user's active status
	if !*user.Active {
		r.resendOTP(ctx, user.Phone)
		return nil, errors.NewAppError(http.StatusExpectationFailed, user.Phone, "Your account is suspended/not verified")
	}

	// Generate refresh token
	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(user.ID)),
		Fullname: user.Fullname,
		Usertype: user.Usertype,
		Campus:   user.Campus,
		Phone:    user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	})
	refreshSS, err := refreshClaims.SignedString([]byte(os.Getenv("REFRESH_SECRET_KEY")))
	if err != nil {
		return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "Failed to generate refresh token")
	}

	// Generate access token
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
		return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "Failed to generate access token")
	}

	// Update user tokens in a single query
	if err := r.db.Model(&user).Updates(map[string]interface{}{
		"refresh_token": refreshSS,
		"access_token":  accessSS,
	}).Error; err != nil {
		return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "Failed to update user tokens")
	}

	// Handle 2FA if enabled
	if user.Twofa != nil && *user.Twofa {
		otpCode := utils.GenerateOTP()
		if err := r.db.Model(&user).Update("code", otpCode).Error; err != nil {
			return nil, errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "Failed to update OTP code")
		}
		return nil, errors.NewAppError(http.StatusProxyAuthRequired, "ACTION REQUIRED", "This account is 2-FA protected, enter OTP to continue")
	}

	// Set access token cookie
	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessSS,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
	SetAccessTokenCookie(ctx, &accessCookie)

	return &LoginUserRes{
		AccessToken:  accessSS,
		RefreshToken: refreshSS,
		ID:           user.ID,
	}, nil
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

	// Create a base query
	query := r.db.Model(&User{})

	// Try to find user by ID first if the filter looks like a number
	if _, err := strconv.ParseUint(filter, 10, 64); err == nil {
		if err := query.Where("id = ?", filter).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist")
			}
			return nil, err
		}
	} else {
		// If not found by ID or filter is not a number, try email or phone
		if err := query.Where("email = ? OR phone = ?", filter, filter).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "User does not exist")
			}
			return nil, err
		}
	}

	// Set default values for nil fields
	if user.Active == nil {
		user.Active = boolPtr(false)
	}
	if user.Twofa == nil {
		user.Twofa = boolPtr(false)
	}
	if user.PaymentDetails == (PaymentDetails{}) {
		user.PaymentDetails = PaymentDetails{}
	}

	return &user, nil
}

// Helper function to initialize user fields
func initializeUserFields(user *User) (*User, error) {
	if user == nil {
		return nil, errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Invalid user data")
	}

	if user.Active == nil {
		user.Active = boolPtr(false)
	}
	if user.Twofa == nil {
		user.Twofa = boolPtr(false)
	}

	// Ensure PaymentDetails is initialized
	if user.PaymentDetails == (PaymentDetails{}) {
		user.PaymentDetails = PaymentDetails{}
	}

	return user, nil
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
	foundStore := &models.Store{}
	if err := r.db.First(foundStore, storeId).Error; err != nil {
		return err
	}

	// Retrieve the user
	userIdStr := strconv.FormatUint(uint64(userId), 10)
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
		newFollowers := make([]*models.Follower, 0)
		for _, follower := range foundStore.Followers {
			if follower.FollowerID != userId {
				newFollowers = append(newFollowers, follower)
			}
		}
		foundStore.Followers = newFollowers

		// Remove the store from user's followedStores
		newFollowedStores := make([]followedStores, 0)
		for _, store := range foundUser.FollowedStores {
			if store.Name != foundStore.Name { // Assuming `ID` is the identifier for the store in followedStores
				newFollowedStores = append(newFollowedStores, store)
			}
		}
		foundUser.FollowedStores = newFollowedStores
	} else {
		// If not following, follow
		foundStore.Followers = append(foundStore.Followers, &models.Follower{
			FollowerID:    userId,
			FollowerName:  foundUser.Fullname,
			FollowerImage: foundUser.Avatar,
		})

		// Add the store to user's followedStores
		foundUser.FollowedStores = append(foundUser.FollowedStores, followedStores{
			Name:        foundStore.Name,
			Description: foundStore.Description,
			Thumbnail:   foundStore.Thumbnail,
			Background:  foundStore.Background,
			Link:        foundStore.Link,
		})
	}

	// Save the  user with GORM
	// if err := r.db.Save(foundStore).Error; err != nil {
	// 	return err
	// }
	if err := r.db.Save(foundUser).Error; err != nil {
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
	if req.UUID != "" {
		existingUser.UUID = req.UUID
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
	if req.Online != existingUser.Online {
		existingUser.Online = req.Online
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
		bodyBytes, _ := io.ReadAll(res.Body)
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
		// Check for the email match and return the account if found
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

	foundUser, err := r.GetUserByEmailOrPhone(ctx, req.Email)
	// If user not found, terminate operation
	if err != nil {
		return err
	}

	// Generate a 32-byte secure random token
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		return err
	}

	// Encode the byte slice as a base64 string for safe URL usage
	t := base64.URLEncoding.EncodeToString(token)

	// Construct reset password link
	resetLink := fmt.Sprintf("%s/reset_password?token=%s&email=%s", req.Link, t, foundUser.Email)

	// Send email
	to := []string{foundUser.Email}
	contents := map[string]string{
		"new_link": resetLink,
	}
	templateID := "7ee50170-1af2-44b4-a819-ab638593f08d"
	err = utils.SendEmail(templateID, "Reset Password Link", to, contents)
	if err != nil {
		log.Printf("Failed to send reset password link email: %v", err)
		return err
	}

	// Store the password reset token in the database
	pwdReset := &PasswordReset{
		Link:      fmt.Sprintf("%s/reset_password?token=%s&email=%s", req.Link, t, foundUser.Email),
		Token:     t,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := r.db.Create(pwdReset).Error; err != nil {
		log.Printf("Failed to store password reset token: %v", err)
		return err
	}

	return nil
}

func (r *repository) VerifyResetLink(ctx context.Context, token string) error {
	pwdReset := &PasswordReset{}

	// Decode the base64 token
	// decodedToken, err := base64.URLEncoding.DecodeString(token)
	// if err != nil {
	// 	return errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "invalid token")
	// }
	// fmt.Print(token)
	// Find the reset request by token
	err := r.db.Where("token = ?", token).First(pwdReset).Error
	if err != nil {
		return errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "invalid or expired token")
	}

	// Check if the token has expired
	if time.Now().After(pwdReset.ExpiresAt) {
		// Delete the expired token
		r.db.Unscoped().Delete(&pwdReset)
		return errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "token has expired")
	}

	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, req *PasswordReset) error {
	pwdReset := &PasswordReset{}
	// Find the user by Email
	user, err := r.GetUserByEmailOrPhone(ctx, req.Email)
	if err != nil {
		return err
	}
	err = r.db.Where("token = ?", req.Token).First(pwdReset).Error
	if err != nil {
		return errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "invalid or expired token")
	}
	// Check if the token has expired
	if time.Now().After(pwdReset.ExpiresAt) {
		// Delete the expired token
		r.db.Unscoped().Delete(&pwdReset)
		return errors.NewAppError(http.StatusInternalServerError, "INTERNAL SERVER ERROR", "token has expired")
	}
	// Update user's password (assuming `generateHash` hashes the password)
	user.Password, _ = utils.HashPasswword(req.Password)
	err = r.db.Save(user).Error
	if err != nil {
		return err
	}

	// After successful password update, delete the password reset token
	err = r.db.Delete(&pwdReset).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetTransactionsByCustomerID(customerID string) ([]Transaction, error) {
	transactionURL := fmt.Sprintf("https://api.paystack.co/transaction?customer=%s", customerID)
	method := "GET"
	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest(method, transactionURL, nil)
	if err != nil {
		return nil, err
	}

	// Set the request headers
	req.Header.Add("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))
	req.Header.Add("Content-Type", "application/json")

	// Make the request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check if the response is successful
	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("error: failed to retrieve transactions with status %d, response: %s", res.StatusCode, string(bodyBytes))
	}

	// Decode the response body
	var transactionsResp TransactionsResponse
	err = json.NewDecoder(res.Body).Decode(&transactionsResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding transactions response: %w", err)
	}

	// Check if Paystack API response was successful
	if !transactionsResp.Status {
		return nil, fmt.Errorf("error: transactions retrieval failed with message: %s", transactionsResp.Message)
	}

	return transactionsResp.Data, nil
}

func (r *repository) GetBalance(ctx context.Context, userId string) error {
	// Get the transactions by customer ID
	transactions, err := r.GetTransactionsByCustomerID(userId)
	if err != nil {
		return err
	}

	// Convert the transactions to a JSON string for logging
	transactionsJSON, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling transactions: %w", err)
	}

	// Print the transactions in JSON format
	fmt.Println(string(transactionsJSON))

	return nil
}

func (r *repository) ConfirmPassword(ctx context.Context, password, userId string) error {
	// Get the transactions by customer ID
	foundUser, err := r.GetUser(ctx, userId)
	if err != nil {
		return err
	}
	if err := utils.CheckPassword(password, foundUser.Password); err != nil {
		return errors.NewAppError(http.StatusBadRequest, "BAD REQUEST", "Password Mismatch!!!")
	}

	return nil
}

func (r *repository) GetMyDownloads(ctx context.Context, userId string) ([]*models.Downloads, error) {
	var downloads []*models.Downloads

	// Use the LIKE operator to check if the userId is contained in the users string
	err := r.db.Where("users LIKE ?", fmt.Sprintf(`%%"%s"%%`, userId)).Find(&downloads).Error
	if err != nil {
		fmt.Printf("Error querying downloads: %v\n", err)
		return nil, err
	}

	return downloads, nil
}

func PayFund(amount float32, accountNumber, bankCode string) error {
	err := utils.PayFund(amount, accountNumber, bankCode)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) SendMaintenanceMail(ctx context.Context, userId string, active bool) error {
	seller, err := r.GetUser(ctx, userId)
	// fmt.Print(seller.Email)

	to := []string{seller.Email}
	contents := map[string]string{
		"seller_name": seller.Fullname,
	}
	if err != nil {
		return err
	}
	if active {
		templateID := "37ec9d6c-cfbf-481e-92ba-782fe1ccd4d1"
		err = utils.SendEmail(templateID, "Your Store is On Hold! 🚫", to, contents)
		if err != nil {
			return err
		}
	} else {
		templateID := "39371fe0-e830-455e-8f06-717630d3d4b9"
		err = utils.SendEmail(templateID, "Your Store is Back in Action! 🔥", to, contents)
		if err != nil {
			return err
		}
	}

	return nil
}

// Add this method to your repository struct
func (r *repository) GetDB() *gorm.DB {
	return r.db
}
