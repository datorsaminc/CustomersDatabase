package model

import (
	"time"
)

// Customer represents a customer record in the database.
type Customer struct {
	ID             string    `json:"id"`
	ProgramVersion string    `json:"programVersion"`
	DeliveryDate   string    `json:"deliveryDate"`
	Name1          string    `json:"name1"`
	Name2          *string   `json:"name2,omitempty"`
	Company        string    `json:"company"`
	VisitAddress   string    `json:"visitAddress"`
	MailingAddress string    `json:"mailingAddress"`
	PostalCodeCity string    `json:"postalCodeCity"`
	LandlinePhone  string    `json:"landlinePhone"`
	MobilePhone    *string   `json:"mobilePhone,omitempty"`
	FaxNumber      *string   `json:"faxNumber,omitempty"`
	Email          string    `json:"email"`
	Comments       string    `json:"comments"`
	Licenses       *int      `json:"licenses,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// CustomerListResponse is the response for listing customers with pagination.
type CustomerListResponse struct {
	Customers []Customer `json:"customers"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}
