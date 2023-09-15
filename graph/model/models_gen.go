// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type AddToCartItemInput struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	User      int    `json:"user"`
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
	SubCategories []*SubCategory `json:"SubCategories,omitempty"`
}

type Follower struct {
	FollowerID    int    `json:"follower_id"`
	FollowerName  string `json:"follower_name"`
	StoreID       int    `json:"store_id"`
	FollowerImage int    `json:"follower_image"`
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

type NewCategory struct {
	Name string `json:"name"`
}

type NewProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	Quantity    int     `json:"quantity"`
	Campus      string  `json:"campus"`
	Variant     string  `json:"variant"`
	Condition   string  `json:"condition"`
	Store       int     `json:"store"`
	Category    int     `json:"category"`
	Subcategory int     `json:"subcategory"`
}

type NewSubCategory struct {
	Name     string `json:"name"`
	Category int    `json:"category"`
}

type NewUser struct {
	Fullname   string     `json:"fullname"`
	Email      string     `json:"email"`
	Campus     string     `json:"campus"`
	Password   string     `json:"password"`
	Stores     []int      `json:"stores,omitempty"`
	Phone      string     `json:"phone"`
	Usertype   string     `json:"usertype"`
	Code       *string    `json:"code,omitempty"`
	Codeexpiry *time.Time `json:"codeexpiry,omitempty"`
	Store      *string    `json:"store,omitempty"`
	Link       *string    `json:"link,omitempty"`
}

type NewVerifyOtp struct {
	Phone string  `json:"phone"`
	Code  string  `json:"code"`
	Email *string `json:"email,omitempty"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Status      bool    `json:"status"`
	Quantity    int     `json:"quantity"`
	Campus      string  `json:"campus"`
	Image       string  `json:"image"`
	Variant     string  `json:"variant"`
	Condition   string  `json:"condition"`
	Store       int     `json:"store"`
	Category    int     `json:"category"`
	Subcategory int     `json:"subcategory"`
}

type Store struct {
	ID                 string      `json:"id"`
	Link               string      `json:"link"`
	Name               string      `json:"name"`
	User               int         `json:"user"`
	Products           []*Product  `json:"products,omitempty"`
	Description        string      `json:"description"`
	Followers          []*Follower `json:"followers,omitempty"`
	HasPhysicalAddress bool        `json:"has_physical_address"`
}

type SubCategory struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Category int    `json:"category"`
}

type User struct {
	ID           string   `json:"id"`
	Fullname     string   `json:"fullname"`
	Email        string   `json:"email"`
	Campus       string   `json:"campus"`
	Password     string   `json:"password"`
	Phone        string   `json:"phone"`
	Usertype     string   `json:"usertype"`
	Stores       []*Store `json:"stores,omitempty"`
	Wallet       *int     `json:"wallet,omitempty"`
	Active       bool     `json:"active"`
	AccessToken  *string  `json:"access_token,omitempty"`
	RefreshToken *string  `json:"refresh_token,omitempty"`
	Twofa        bool     `json:"twofa"`
	Code         string   `json:"code"`
	Codeexpiry   string   `json:"codeexpiry"`
}

type VerifyOtp struct {
	Phone string  `json:"phone"`
	Code  string  `json:"code"`
	Email *string `json:"email,omitempty"`
}
