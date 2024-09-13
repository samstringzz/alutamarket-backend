// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Account struct {
	Customer      *Customer    `json:"customer"`
	Bank          *Bank        `json:"bank"`
	ID            int          `json:"id"`
	AccountNumber int          `json:"account_number"`
	AccountName   string       `json:"account_name"`
	CreatedAt     string       `json:"created_at"`
	UpdatedAt     string       `json:"updated_at"`
	SplitConfig   *SplitConfig `json:"split_config"`
	Active        bool         `json:"active"`
	Assigned      bool         `json:"assigned"`
}

type Bank struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

type BundleVariation struct {
	VariationCode   string `json:"variationCode"`
	Name            string `json:"name"`
	VariationAmount string `json:"variationAmount"`
	FixedPrice      string `json:"fixedPrice"`
}

type Cart struct {
	Items  []*CartItem `json:"items"`
	Total  float64     `json:"total"`
	Active bool        `json:"active"`
	User   int         `json:"user"`
	ID     *string     `json:"id,omitempty"`
}

type CartItem struct {
	Product  *Product `json:"product"`
	Quantity int      `json:"quantity"`
}

type Category struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Slug          string         `json:"slug"`
	Subcategories []*SubCategory `json:"subcategories,omitempty"`
}

type Chat struct {
	ID            *int           `json:"id,omitempty"`
	LatestMessage *Message       `json:"latest_message,omitempty"`
	Messages      []*Message     `json:"messages,omitempty"`
	UnreadCount   int            `json:"unread_count"`
	Users         []*MessageUser `json:"users"`
	Time          time.Time      `json:"time"`
}

type ChatInput struct {
	Users []*MessageUserInput `json:"users,omitempty"`
}

type Customer struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	CustomerCode string `json:"customer_code"`
	Phone        string `json:"phone"`
	RiskAction   string `json:"risk_action"`
}

type DVAAccountInput struct {
	UserID    string `json:"user_id"`
	StoreName string `json:"store_name"`
}

type DVADetails struct {
	Surname       string `json:"surname"`
	Othername     string `json:"othername"`
	Bvn           string `json:"bvn"`
	Country       string `json:"country"`
	BankCode      string `json:"bank_code"`
	AccountNumber string `json:"account_number"`
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	StoreName     string `json:"store_name"`
}

type HandledProducts struct {
	UserID           int      `json:"userId"`
	ProductID        int      `json:"productId"`
	ProductName      *string  `json:"productName,omitempty"`
	ProductThumbnail *string  `json:"productThumbnail,omitempty"`
	ProductPrice     *float64 `json:"productPrice,omitempty"`
	ProductDiscount  *float64 `json:"productDiscount,omitempty"`
	ProductStatus    *bool    `json:"productStatus,omitempty"`
	ProductQuantity  *int     `json:"productQuantity,omitempty"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	ID           int    `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Message struct {
	ID        int            `json:"id"`
	ChatID    int            `json:"chat_id"`
	Sender    int            `json:"sender"`
	Content   string         `json:"content"`
	Users     []*MessageUser `json:"users"`
	Media     MediaType      `json:"media"`
	IsRead    bool           `json:"is_read"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
}

type MessageInput struct {
	ID      *int       `json:"id,omitempty"`
	ChatID  int        `json:"chat_id"`
	Content string     `json:"content"`
	Sender  int        `json:"sender"`
	Media   *MediaType `json:"media,omitempty"`
	IsRead  bool       `json:"is_read"`
}

type MessageUser struct {
	ID     int     `json:"id"`
	Avatar *string `json:"avatar,omitempty"`
	Online bool    `json:"online"`
	Name   string  `json:"name"`
	Status string  `json:"status"`
}

type MessageUserInput struct {
	ID     int     `json:"id"`
	Avatar *string `json:"avatar,omitempty"`
	Name   string  `json:"name"`
}

type ModifyCartItemInput struct {
	ProductID   *string `json:"productId,omitempty"`
	ProductName *string `json:"productName,omitempty"`
	Quantity    int     `json:"quantity"`
	User        int     `json:"user"`
}

type Mutation struct {
}

type NewCategory struct {
	Name string `json:"name"`
}

type NewReview struct {
	Message   string  `json:"message"`
	Rating    float64 `json:"rating"`
	ProductID string  `json:"product_id"`
	Image     string  `json:"image"`
	Username  string  `json:"username"`
}

type NewSubCategory struct {
	Name     string `json:"name"`
	Category int    `json:"category"`
}

type NewUser struct {
	Fullname   string      `json:"fullname"`
	Email      string      `json:"email"`
	Campus     string      `json:"campus"`
	Password   string      `json:"password"`
	Stores     *StoreInput `json:"stores,omitempty"`
	Phone      string      `json:"phone"`
	Usertype   string      `json:"usertype"`
	Code       *string     `json:"code,omitempty"`
	Codeexpiry *time.Time  `json:"codeexpiry,omitempty"`
}

type NewVariant struct {
	Name  string             `json:"name"`
	Value []*NewVariantValue `json:"value"`
}

type NewVariantValue struct {
	Value  string   `json:"value"`
	Price  *float64 `json:"price,omitempty"`
	Images []string `json:"images,omitempty"`
}

type NewVerifyOtp struct {
	Phone    string  `json:"phone"`
	Code     string  `json:"code"`
	Email    *string `json:"email,omitempty"`
	Attempts int     `json:"attempts"`
}

type Order struct {
	ID            int     `json:"id"`
	Customer      string  `json:"customer"`
	CustomerEmail string  `json:"customer_email"`
	Price         float64 `json:"price"`
	Status        string  `json:"status"`
	Date          string  `json:"date"`
	StoreID       string  `json:"store_id"`
}

type PaymentData struct {
	StoreID        []int    `json:"storeID,omitempty"`
	Status         *string  `json:"status,omitempty"`
	UserID         string   `json:"userID"`
	Amount         *float64 `json:"amount,omitempty"`
	UUID           *string  `json:"UUID,omitempty"`
	PaymentGateway string   `json:"paymentGateway"`
}

type PaymentDetails struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Info    string `json:"info"`
}

type PaymentDetailsInput struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Info    string `json:"info"`
}

type Product struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	Slug            string     `json:"slug"`
	Description     string     `json:"description"`
	Price           float64    `json:"price"`
	Discount        float64    `json:"discount"`
	Status          bool       `json:"status"`
	Quantity        int        `json:"quantity"`
	Thumbnail       string     `json:"thumbnail"`
	Image           []string   `json:"image"`
	File            *string    `json:"file,omitempty"`
	Variant         []*Variant `json:"variant,omitempty"`
	Review          []*Review  `json:"review,omitempty"`
	Store           string     `json:"store"`
	Category        string     `json:"category"`
	Subcategory     string     `json:"subcategory"`
	AlwaysAvailable bool       `json:"always_available"`
}

type ProductInput struct {
	Name            string        `json:"name"`
	ID              *string       `json:"id,omitempty"`
	Description     string        `json:"description"`
	File            string        `json:"file"`
	Price           float64       `json:"price"`
	Discount        float64       `json:"discount"`
	Thumbnail       string        `json:"thumbnail"`
	Image           []string      `json:"image"`
	Quantity        int           `json:"quantity"`
	Variant         []*NewVariant `json:"variant,omitempty"`
	Review          []*NewReview  `json:"review,omitempty"`
	Store           string        `json:"store"`
	Category        int           `json:"category"`
	Subcategory     int           `json:"subcategory"`
	AlwaysAvailable bool          `json:"always_available"`
}

type ProductPaginationData struct {
	Data        []*Product `json:"data"`
	CurrentPage int        `json:"current_page"`
	PerPage     int        `json:"per_page"`
	Total       int        `json:"total"`
	NextPage    int        `json:"next_page"`
	PrevPage    int        `json:"prev_page"`
}

type PurchasedOrder struct {
	CartID          int               `json:"cart_id"`
	Coupon          string            `json:"coupon"`
	Fee             float64           `json:"fee"`
	Status          string            `json:"status"`
	UserID          string            `json:"user_id"`
	Amount          float64           `json:"amount"`
	UUID            string            `json:"uuid"`
	PaymentGateway  string            `json:"paymentGateway"`
	PaymentMethod   string            `json:"paymentMethod"`
	TransRef        string            `json:"transRef"`
	TransStatus     string            `json:"transStatus"`
	Products        []*TrackedProduct `json:"products"`
	DeliveryDetails *DeliveryDetails  `json:"deliveryDetails"`
	TextRef         string            `json:"textRef"`
}

type Query struct {
}

type Review struct {
	Rating    float64 `json:"rating"`
	Message   string  `json:"message"`
	Image     string  `json:"image"`
	ProductID int     `json:"product_id"`
	Username  string  `json:"username"`
	ID        *string `json:"id,omitempty"`
}

type ReviewInput struct {
	Username  string  `json:"username"`
	Image     string  `json:"image"`
	Message   string  `json:"message"`
	Rating    float64 `json:"rating"`
	ProductID int     `json:"productId"`
}

type Skynet struct {
	ID            string  `json:"id"`
	UserID        *string `json:"user_id,omitempty"`
	Status        *string `json:"status,omitempty"`
	RequestID     string  `json:"request_id"`
	TransactionID *string `json:"transaction_id,omitempty"`
	Type          *string `json:"type,omitempty"`
	Receiever     *string `json:"receiever,omitempty"`
}

type SkynetInput struct {
	Amount           int     `json:"amount"`
	UserID           int     `json:"user_id"`
	BillersCode      *string `json:"billers_code,omitempty"`
	VariantCode      *string `json:"variant_code,omitempty"`
	ServiceID        string  `json:"service_id"`
	PhoneNumber      *string `json:"phone_number,omitempty"`
	Quantity         *string `json:"quantity,omitempty"`
	SubscriptionType *string `json:"subscription_type,omitempty"`
	Type             string  `json:"type"`
}

type SmartCardInput struct {
	ServiceID   string  `json:"service_id"`
	BillersCode string  `json:"billers_code"`
	CardType    *string `json:"card_type,omitempty"`
}

type SmartcardContent struct {
	CustomerName       string  `json:"customerName"`
	Status             string  `json:"status"`
	DueDate            string  `json:"dueDate"`
	CustomerNumber     int     `json:"customerNumber"`
	CustomerType       string  `json:"customerType"`
	CurrentBouquet     string  `json:"currentBouquet"`
	CurrentBouquetCode string  `json:"currentBouquetCode"`
	RenewalAmount      float64 `json:"renewalAmount"`
}

type SmartcardVerificationResponse struct {
	Code    string            `json:"code"`
	Content *SmartcardContent `json:"content"`
}

type SplitConfig struct {
	Subaccount string `json:"Subaccount"`
}

type Store struct {
	ID                 string           `json:"id"`
	Link               string           `json:"link"`
	Name               string           `json:"name"`
	Wallet             float64          `json:"wallet"`
	User               int              `json:"user"`
	Email              string           `json:"email"`
	Description        string           `json:"description"`
	Followers          []*StoreFollower `json:"followers,omitempty"`
	Product            []*Product       `json:"product,omitempty"`
	Transactions       []*Transaction   `json:"transactions,omitempty"`
	Orders             []*StoreOrder    `json:"orders,omitempty"`
	Address            string           `json:"address"`
	Status             bool             `json:"status"`
	Thumbnail          string           `json:"thumbnail"`
	Phone              string           `json:"phone"`
	Background         string           `json:"background"`
	HasPhysicalAddress bool             `json:"has_physical_address"`
	Visitors           []string         `json:"visitors"`
}

type StoreCustomer struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type StoreFollower struct {
	FollowerID    int    `json:"follower_id"`
	FollowerName  string `json:"follower_name"`
	StoreID       int    `json:"store_id"`
	FollowerImage string `json:"follower_image"`
}

type StoreFollowerInput struct {
	FollowerID    int    `json:"follower_id"`
	FollowerName  string `json:"follower_name"`
	FollowerImage string `json:"follower_image"`
	StoreID       int    `json:"store_id"`
	Action        string `json:"action"`
}

type StoreInput struct {
	ID                 *string `json:"id,omitempty"`
	Link               string  `json:"link"`
	Name               string  `json:"name"`
	User               int     `json:"user"`
	Description        string  `json:"description"`
	Address            string  `json:"address"`
	Wallet             int     `json:"wallet"`
	HasPhysicalAddress bool    `json:"has_physical_address"`
	Phone              string  `json:"phone"`
	Status             bool    `json:"status"`
	Email              *string `json:"email,omitempty"`
	Thumbnail          *string `json:"thumbnail,omitempty"`
	Background         *string `json:"background,omitempty"`
}

type StoreOrder struct {
	StoreID   string         `json:"store_id"`
	Product   []*Product     `json:"product,omitempty"`
	TrtRef    string         `json:"trtRef"`
	Active    bool           `json:"active"`
	Status    string         `json:"status"`
	Customer  *StoreCustomer `json:"customer"`
	UUID      string         `json:"uuid"`
	CreatedAt time.Time      `json:"createdAt"`
}

type StoreOrderInput struct {
	StoreID  string               `json:"store_id"`
	Product  []*StoreProductInput `json:"product,omitempty"`
	Status   string               `json:"status"`
	Customer *CustomerInput       `json:"customer"`
}

type StorePaginationData struct {
	Data        []*Store `json:"data"`
	CurrentPage int      `json:"current_page"`
	PerPage     int      `json:"per_page"`
	Total       int      `json:"total"`
}

type StoreProductInput struct {
	Name      string  `json:"name"`
	Thumbnail string  `json:"thumbnail"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	ID        *string `json:"id,omitempty"`
}

type SubCategory struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Category int    `json:"category"`
}

type Subscription struct {
}

type SubscriptionBundle struct {
	ServiceName    string             `json:"serviceName"`
	ServiceID      string             `json:"serviceID"`
	ConvinienceFee string             `json:"convinienceFee"`
	Variations     []*BundleVariation `json:"variations"`
}

type TrackedProduct struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Thumbnail string  `json:"thumbnail"`
	Price     float64 `json:"price"`
	Discount  float64 `json:"discount"`
	Status    string  `json:"status"`
}

type Transaction struct {
	StoreID   int       `json:"storeID"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	User      string    `json:"user"`
	Amount    float64   `json:"amount"`
	UUID      string    `json:"UUID"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionInput struct {
	StoreID  string   `json:"store_id"`
	Status   string   `json:"status"`
	User     string   `json:"user"`
	Amount   *float64 `json:"amount,omitempty"`
	Type     string   `json:"type"`
	Category string   `json:"category"`
}

type UpdateStoreInput struct {
	ID                 *string `json:"id,omitempty"`
	Link               *string `json:"link,omitempty"`
	Name               *string `json:"name,omitempty"`
	User               *int    `json:"user,omitempty"`
	Description        *string `json:"description,omitempty"`
	Address            *string `json:"address,omitempty"`
	Wallet             *int    `json:"wallet,omitempty"`
	HasPhysicalAddress *bool   `json:"has_physical_address,omitempty"`
	Status             *bool   `json:"status,omitempty"`
	Phone              *string `json:"phone,omitempty"`
	Email              *string `json:"email,omitempty"`
	Thumbnail          *string `json:"thumbnail,omitempty"`
	Background         *string `json:"background,omitempty"`
	Visitor            *string `json:"visitor,omitempty"`
}

type UpdateStoreOrderInput struct {
	ID      *string `json:"id,omitempty"`
	Status  *string `json:"status,omitempty"`
	StoreID *string `json:"store_id,omitempty"`
}

type UpdateUserInput struct {
	ID             *string              `json:"id,omitempty"`
	Fullname       *string              `json:"fullname,omitempty"`
	Email          *string              `json:"email,omitempty"`
	Campus         *string              `json:"campus,omitempty"`
	Password       *string              `json:"password,omitempty"`
	Stores         *StoreInput          `json:"stores,omitempty"`
	Dob            *string              `json:"dob,omitempty"`
	Phone          *string              `json:"phone,omitempty"`
	Gender         *string              `json:"gender,omitempty"`
	Active         *bool                `json:"active,omitempty"`
	Usertype       *string              `json:"usertype,omitempty"`
	Code           *string              `json:"code,omitempty"`
	Avatar         *string              `json:"avatar,omitempty"`
	PaymnetDetails *PaymentDetailsInput `json:"paymnetDetails,omitempty"`
}

type User struct {
	ID             string          `json:"id"`
	Fullname       string          `json:"fullname"`
	Email          string          `json:"email"`
	Campus         string          `json:"campus"`
	Avatar         *string         `json:"avatar,omitempty"`
	Dob            *string         `json:"dob,omitempty"`
	Gender         *string         `json:"gender,omitempty"`
	Password       string          `json:"password"`
	Phone          string          `json:"phone"`
	Usertype       string          `json:"usertype"`
	Stores         []*Store        `json:"stores,omitempty"`
	Active         bool            `json:"active"`
	AccessToken    *string         `json:"access_token,omitempty"`
	RefreshToken   *string         `json:"refresh_token,omitempty"`
	Twofa          bool            `json:"twofa"`
	Code           string          `json:"code"`
	PaymnetDetails *PaymentDetails `json:"paymnetDetails,omitempty"`
	Codeexpiry     string          `json:"codeexpiry"`
}

type Variant struct {
	Name  string          `json:"name"`
	Value []*VariantValue `json:"value"`
}

type VariantValue struct {
	Value  string   `json:"value"`
	Price  float64  `json:"price"`
	Images []string `json:"images,omitempty"`
}

type VerifyOtp struct {
	Phone string  `json:"phone"`
	Code  string  `json:"code"`
	Email *string `json:"email,omitempty"`
}

type CustomerInput struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type DeliveryDetails struct {
	Method  string  `json:"method"`
	Address string  `json:"address"`
	Fee     float64 `json:"fee"`
}

type PasswordResetInput struct {
	Link  string `json:"link"`
	Email string `json:"email"`
}

type PasswordUpdateInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Verifyotpinput struct {
	Code     string `json:"code"`
	Phone    string `json:"phone"`
	Attempts int    `json:"attempts"`
}

type MediaType string

const (
	MediaTypeImage    MediaType = "IMAGE"
	MediaTypeVideo    MediaType = "VIDEO"
	MediaTypeAudio    MediaType = "AUDIO"
	MediaTypeDocument MediaType = "DOCUMENT"
)

var AllMediaType = []MediaType{
	MediaTypeImage,
	MediaTypeVideo,
	MediaTypeAudio,
	MediaTypeDocument,
}

func (e MediaType) IsValid() bool {
	switch e {
	case MediaTypeImage, MediaTypeVideo, MediaTypeAudio, MediaTypeDocument:
		return true
	}
	return false
}

func (e MediaType) String() string {
	return string(e)
}

func (e *MediaType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MediaType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MediaType", str)
	}
	return nil
}

func (e MediaType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RoleType string

const (
	RoleTypeReciver RoleType = "RECIVER"
	RoleTypeSender  RoleType = "SENDER"
)

var AllRoleType = []RoleType{
	RoleTypeReciver,
	RoleTypeSender,
}

func (e RoleType) IsValid() bool {
	switch e {
	case RoleTypeReciver, RoleTypeSender:
		return true
	}
	return false
}

func (e RoleType) String() string {
	return string(e)
}

func (e *RoleType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RoleType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RoleType", str)
	}
	return nil
}

func (e RoleType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
