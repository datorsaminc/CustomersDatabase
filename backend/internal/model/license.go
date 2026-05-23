package model

// License represents a license record linked to a customer by ID.
type License struct {
	ID       string `json:"id"`
	Licences int    `json:"licences"`
}
