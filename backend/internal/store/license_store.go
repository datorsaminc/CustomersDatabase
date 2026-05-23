package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"clientsweb-backend/internal/model"
)

const defaultLicenseFile = "data/licenses.json"

// LicenseStore provides thread-safe file-based storage for licenses.
type LicenseStore struct {
	mu   sync.RWMutex
	file string
	data []model.License
}

// NewLicenseStore creates a new store, loading existing data from the specified file.
func NewLicenseStore(filePath string) (*LicenseStore, error) {
	s := &LicenseStore{
		file: filePath,
		data: make([]model.License, 0),
	}

	if err := s.load(); err != nil {
		return nil, fmt.Errorf("failed to load license store: %w", err)
	}

	return s, nil
}

// load reads licenses from the JSON file into memory.
func (s *LicenseStore) load() error {
	data, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading license file: %w", err)
	}

	var licenses []model.License
	if err := json.Unmarshal(data, &licenses); err != nil {
		return fmt.Errorf("unmarshaling license JSON: %w", err)
	}

	s.data = licenses
	return nil
}

// save writes the current in-memory data to the JSON file.
func (s *LicenseStore) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling license JSON: %w", err)
	}

	if err := os.WriteFile(s.file, data, 0644); err != nil {
		return fmt.Errorf("writing license file: %w", err)
	}

	return nil
}

// GetByCustomerID returns the license for a given customer ID.
func (s *LicenseStore) GetByCustomerID(customerID string) (*model.License, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, l := range s.data {
		if l.ID == customerID {
			result := l
			return &result, nil
		}
	}

	return nil, fmt.Errorf("license not found for customer: %s", customerID)
}

// GetAll returns all licenses.
func (s *LicenseStore) GetAll() []model.License {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.License, len(s.data))
	copy(result, s.data)
	return result
}

// Upsert creates or updates a license for the given customer ID.
func (s *LicenseStore) Upsert(customerID string, licences int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, l := range s.data {
		if l.ID == customerID {
			s.data[i].Licences = licences
			return s.save()
		}
	}

	// Not found, create new entry
	newLicense := model.License{
		ID:       customerID,
		Licences: licences,
	}
	s.data = append(s.data, newLicense)
	return s.save()
}

// Delete removes a license by customer ID.
func (s *LicenseStore) Delete(customerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, l := range s.data {
		if l.ID == customerID {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return s.save()
		}
	}

	return fmt.Errorf("license not found for customer: %s", customerID)
}

// TotalLicenses returns the sum of all licences across all records.
func (s *LicenseStore) TotalLicenses() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := 0
	for _, l := range s.data {
		total += l.Licences
	}
	return total
}
