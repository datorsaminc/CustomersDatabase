package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"clientsweb-backend/internal/model"
)

const defaultDataFile = "data/customers.json"

// JSONStore provides thread-safe file-based storage for customers.
type JSONStore struct {
	mu   sync.RWMutex
	file string
	data []model.Customer
}

// NewJSONStore creates a new store, loading existing data from the specified file.
func NewJSONStore(filePath string) (*JSONStore, error) {
	s := &JSONStore{
		file: filePath,
		data: make([]model.Customer, 0),
	}

	if err := s.load(); err != nil {
		return nil, fmt.Errorf("failed to load store: %w", err)
	}

	return s, nil
}

// load reads customers from the JSON file into memory.
func (s *JSONStore) load() error {
	data, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, start with empty data
			return nil
		}
		return fmt.Errorf("reading file: %w", err)
	}

	var customers []model.Customer
	if err := json.Unmarshal(data, &customers); err != nil {
		return fmt.Errorf("unmarshaling JSON: %w", err)
	}

	s.data = customers
	return nil
}

// save writes the current in-memory data to the JSON file.
func (s *JSONStore) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	if err := os.WriteFile(s.file, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// GetAll returns all customers.
func (s *JSONStore) GetAll() []model.Customer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Customer, len(s.data))
	copy(result, s.data)
	return result
}

// GetByID returns a customer by ID, or an error if not found.
func (s *JSONStore) GetByID(id string) (*model.Customer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.data {
		if c.ID == id {
			result := c
			return &result, nil
		}
	}

	return nil, fmt.Errorf("customer not found: %s", id)
}

// Create adds a new customer and persists the change.
func (s *JSONStore) Create(customer model.Customer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, customer)
	return s.save()
}

// Update replaces an existing customer by ID and persists the change.
func (s *JSONStore) Update(id string, updated model.Customer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, c := range s.data {
		if c.ID == id {
			s.data[i] = updated
			return s.save()
		}
	}

	return fmt.Errorf("customer not found: %s", id)
}

// Delete removes a customer by ID and persists the change.
func (s *JSONStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, c := range s.data {
		if c.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return s.save()
		}
	}

	return fmt.Errorf("customer not found: %s", id)
}

// GetAllUnlocked returns all customers without locking (for use when lock is already held).
func (s *JSONStore) GetAllUnlocked() []model.Customer {
	result := make([]model.Customer, len(s.data))
	copy(result, s.data)
	return result
}
