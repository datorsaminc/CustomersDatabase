package search

import (
	"strings"

	"clientsweb-backend/internal/model"
)

// Search performs a case-insensitive full-text search across all customer fields.
func Search(customers []model.Customer, query string) []model.Customer {
	if strings.TrimSpace(query) == "" {
		return customers
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var results []model.Customer

	for _, c := range customers {
		if matchesQuery(c, query) {
			results = append(results, c)
		}
	}

	return results
}

// matchesQuery checks if any searchable field contains the query string.
func matchesQuery(c model.Customer, query string) bool {
	fields := []string{
		c.ProgramVersion,
		c.DeliveryDate,
		c.Name1,
		c.Company,
		c.VisitAddress,
		c.MailingAddress,
		c.PostalCodeCity,
		c.LandlinePhone,
		c.Email,
		c.Comments,
	}

	if c.Name2 != nil {
		fields = append(fields, *c.Name2)
	}
	if c.MobilePhone != nil {
		fields = append(fields, *c.MobilePhone)
	}
	if c.FaxNumber != nil {
		fields = append(fields, *c.FaxNumber)
	}

	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}

	return false
}
