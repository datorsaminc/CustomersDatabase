package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"clientsweb-backend/internal/model"
	"clientsweb-backend/internal/search"
	"clientsweb-backend/internal/store"
)

// CustomerService provides business logic for customer operations.
type CustomerService struct {
	store        *store.JSONStore
	licenseStore *store.LicenseStore
}

// NewCustomerService creates a new service with the given stores.
func NewCustomerService(store *store.JSONStore, licenseStore *store.LicenseStore) *CustomerService {
	return &CustomerService{store: store, licenseStore: licenseStore}
}

// ListCustomers returns paginated customers enriched with license data.
func (s *CustomerService) ListCustomers(query string, page int, limit int) (*model.CustomerListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 200 {
		limit = 50
	}

	all := s.store.GetAll()
	filtered := search.Search(all, query)
	total := len(filtered)

	start := (page - 1) * limit
	end := start + limit
	if end > total {
		end = total
	}

	var customers []model.Customer
	if start < total {
		customers = filtered[start:end]
	} else {
		customers = make([]model.Customer, 0)
	}

	// Enrich customers with license data
	for i := range customers {
		license, err := s.licenseStore.GetByCustomerID(customers[i].ID)
		if err == nil && license != nil {
			customers[i].Licenses = &license.Licences
		}
	}

	return &model.CustomerListResponse{
		Customers: customers,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// GetCustomer returns a single customer enriched with license data.
func (s *CustomerService) GetCustomer(id string) (*model.Customer, error) {
	customer, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	license, err := s.licenseStore.GetByCustomerID(customer.ID)
	if err == nil && license != nil {
		customer.Licenses = &license.Licences
	}

	return customer, nil
}

// CreateCustomer creates a new customer with generated ID and timestamps, and optionally saves license data.
func (s *CustomerService) CreateCustomer(input model.Customer) (*model.Customer, error) {
	now := time.Now().UTC()
	customer := model.Customer{
		ID:             uuid.NewString(),
		ProgramVersion: input.ProgramVersion,
		DeliveryDate:   input.DeliveryDate,
		Name1:          input.Name1,
		Name2:          input.Name2,
		Company:        input.Company,
		VisitAddress:   input.VisitAddress,
		MailingAddress: input.MailingAddress,
		PostalCodeCity: input.PostalCodeCity,
		LandlinePhone:  input.LandlinePhone,
		MobilePhone:    input.MobilePhone,
		FaxNumber:      input.FaxNumber,
		Email:          input.Email,
		Comments:       input.Comments,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.store.Create(customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	// Save license data if provided
	if input.Licenses != nil && *input.Licenses > 0 {
		if err := s.licenseStore.Upsert(customer.ID, *input.Licenses); err != nil {
			return nil, fmt.Errorf("failed to save license: %w", err)
		}
		customer.Licenses = input.Licenses
	}

	return &customer, nil
}

// UpdateCustomer updates an existing customer by ID and optionally saves license data.
func (s *CustomerService) UpdateCustomer(id string, input model.Customer) (*model.Customer, error) {
	existing, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	updated := model.Customer{
		ID:             id,
		ProgramVersion: input.ProgramVersion,
		DeliveryDate:   input.DeliveryDate,
		Name1:          input.Name1,
		Name2:          input.Name2,
		Company:        input.Company,
		VisitAddress:   input.VisitAddress,
		MailingAddress: input.MailingAddress,
		PostalCodeCity: input.PostalCodeCity,
		LandlinePhone:  input.LandlinePhone,
		MobilePhone:    input.MobilePhone,
		FaxNumber:      input.FaxNumber,
		Email:          input.Email,
		Comments:       input.Comments,
		CreatedAt:      existing.CreatedAt,
		UpdatedAt:      time.Now().UTC(),
	}

	if err := s.store.Update(id, updated); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	// Save license data if provided
	if input.Licenses != nil && *input.Licenses > 0 {
		if err := s.licenseStore.Upsert(id, *input.Licenses); err != nil {
			return nil, fmt.Errorf("failed to save license: %w", err)
		}
		updated.Licenses = input.Licenses
	}

	return &updated, nil
}

// DeleteCustomer removes a customer by ID and their associated license record.
func (s *CustomerService) DeleteCustomer(id string) error {
	if err := s.store.Delete(id); err != nil {
		return err
	}
	_ = s.licenseStore.Delete(id) // Best effort, ignore errors
	return nil
}

// TotalLicenses returns the total number of licenses sold across all customers.
func (s *CustomerService) TotalLicenses() int {
	return s.licenseStore.TotalLicenses()
}
