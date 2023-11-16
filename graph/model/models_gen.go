// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

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

type Follower struct {
	FollowerID    int    `json:"follower_id"`
	FollowerName  string `json:"follower_name"`
	StoreID       int    `json:"store_id"`
	FollowerImage string `json:"follower_image"`
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

type ModifyCartItemInput struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	User      int    `json:"user"`
}

type NewCategory struct {
	Name string `json:"name"`
}

type NewHandleProductInput struct {
	User    int `json:"user"`
	Product int `json:"product"`
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
	Phone string  `json:"phone"`
	Code  string  `json:"code"`
	Email *string `json:"email,omitempty"`
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

type Product struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Discount    float64    `json:"discount"`
	Status      bool       `json:"status"`
	Quantity    int        `json:"quantity"`
	Thumbnail   string     `json:"thumbnail"`
	Image       []string   `json:"image"`
	Variant     []*Variant `json:"variant,omitempty"`
	Store       string     `json:"store"`
	Category    string     `json:"category"`
	Subcategory string     `json:"subcategory"`
}

type ProductInput struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Discount    float64       `json:"discount"`
	Thumbnail   string        `json:"thumbnail"`
	Image       []string      `json:"image"`
	Quantity    int           `json:"quantity"`
	Variant     []*NewVariant `json:"variant,omitempty"`
	Store       string        `json:"store"`
	Category    int           `json:"category"`
	Subcategory int           `json:"subcategory"`
}

type ProductPaginationData struct {
	Data        []*Product `json:"data"`
	CurrentPage int        `json:"current_page"`
	PerPage     int        `json:"per_page"`
	Total       int        `json:"total"`
}

type Review struct {
	Username  string  `json:"username"`
	Image     string  `json:"image"`
	Message   string  `json:"message"`
	Rating    float64 `json:"rating"`
	ProductID int     `json:"productId"`
}

type ReviewInput struct {
	Username  string  `json:"username"`
	Image     string  `json:"image"`
	Message   string  `json:"message"`
	Rating    float64 `json:"rating"`
	ProductID int     `json:"productId"`
}

type Store struct {
	ID                 string      `json:"id"`
	Link               string      `json:"link"`
	Name               string      `json:"name"`
	Wallet             float64     `json:"wallet"`
	User               int         `json:"user"`
	Description        string      `json:"description"`
	Followers          []*Follower `json:"followers,omitempty"`
	Product            []*Product  `json:"product,omitempty"`
	Order              []*Order    `json:"order,omitempty"`
	Address            string      `json:"address"`
	Status             bool        `json:"status"`
	Thumbnail          string      `json:"thumbnail"`
	Phone              string      `json:"phone"`
	Background         string      `json:"background"`
	HasPhysicalAddress bool        `json:"has_physical_address"`
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
	Email              *string `json:"email,omitempty"`
	Thumbnail          *string `json:"thumbnail,omitempty"`
	Background         *string `json:"background,omitempty"`
}

type StorePaginationData struct {
	Data        []*Store `json:"data"`
	CurrentPage int      `json:"current_page"`
	PerPage     int      `json:"per_page"`
	Total       int      `json:"total"`
}

type SubCategory struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Category int    `json:"category"`
}

type Transaction struct {
	StoreID        []int    `json:"storeID,omitempty"`
	CartID         int      `json:"cartID"`
	Status         *string  `json:"status,omitempty"`
	UserID         string   `json:"userID"`
	Amount         *float64 `json:"amount,omitempty"`
	UUID           *string  `json:"UUID,omitempty"`
	PaymentGateway *string  `json:"paymentGateway,omitempty"`
}

type User struct {
	ID           string   `json:"id"`
	Fullname     string   `json:"fullname"`
	Email        string   `json:"email"`
	Campus       string   `json:"campus"`
	Avatar       *string  `json:"avatar,omitempty"`
	Password     string   `json:"password"`
	Phone        string   `json:"phone"`
	Usertype     string   `json:"usertype"`
	Stores       []*Store `json:"stores,omitempty"`
	Active       bool     `json:"active"`
	AccessToken  *string  `json:"access_token,omitempty"`
	RefreshToken *string  `json:"refresh_token,omitempty"`
	Twofa        bool     `json:"twofa"`
	Code         string   `json:"code"`
	Codeexpiry   string   `json:"codeexpiry"`
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
