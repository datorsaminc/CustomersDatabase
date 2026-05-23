package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"clientsweb-backend/internal/model"
	"clientsweb-backend/internal/service"
)

// CustomerHandler handles HTTP requests for customer operations.
type CustomerHandler struct {
	service *service.CustomerService
}

// NewCustomerHandler creates a new handler with the given service.
func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// Routes returns a router with all customer routes registered.
func (h *CustomerHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListCustomers)
	r.Post("/", h.CreateCustomer)
	r.Get("/{id}", h.GetCustomer)
	r.Put("/{id}", h.UpdateCustomer)
	r.Delete("/{id}", h.DeleteCustomer)
	return r
}

// LicenseStats handles GET /api/licenses/stats.
func (h *CustomerHandler) LicenseStats(w http.ResponseWriter, r *http.Request) {
	total := h.service.TotalLicenses()
	writeJSON(w, map[string]int{"totalLicenses": total})
}

// ListCustomers handles GET /api/customers?search=query&page=1&limit=50.
func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	response, err := h.service.ListCustomers(query, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, response)
}

// GetCustomer handles GET /api/customers/:id.
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "customer ID is required"}`, http.StatusBadRequest)
		return
	}

	customer, err := h.service.GetCustomer(id)
	if err != nil {
		http.Error(w, `{"error": "customer not found"}`, http.StatusNotFound)
		return
	}

	writeJSON(w, customer)
}

// CreateCustomer handles POST /api/customers.
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var input model.Customer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Company) == "" || strings.TrimSpace(input.Name1) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "company and name1 are required fields",
		})
		return
	}

	customer, err := h.service.CreateCustomer(input)
	if err != nil {
		http.Error(w, `{"error": "failed to create customer"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, customer)
}

// UpdateCustomer handles PUT /api/customers/:id.
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "customer ID is required"}`, http.StatusBadRequest)
		return
	}

	var input model.Customer
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	customer, err := h.service.UpdateCustomer(id, input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "customer not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "failed to update customer"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, customer)
}

// DeleteCustomer handles DELETE /api/customers/:id.
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "customer ID is required"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCustomer(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, `{"error": "customer not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "failed to delete customer"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "customer deleted successfully"})
}

// writeJSON writes a JSON response with Content-Type header.
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
