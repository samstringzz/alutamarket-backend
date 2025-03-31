package user

import "encoding/json"

type PaystackCustomer struct {
	ID           json.Number `json:"id"`
	FirstName    string      `json:"first_name"`
	LastName     string      `json:"last_name"`
	Email        string      `json:"email"`
	CustomerCode string      `json:"customer_code"`
	Phone        string      `json:"phone"`
	RiskAction   string      `json:"risk_action"`
}

type PaystackDVAResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Customer      PaystackCustomer `json:"customer"`
		AccountName   string           `json:"account_name"`
		AccountNumber string           `json:"account_number"`
		Bank          struct {
			Name string `json:"name"`
			ID   string `json:"id"`
			Slug string `json:"slug"`
		} `json:"bank"`
	} `json:"data"`
}
