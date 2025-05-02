package dto

import (
	"time"
)

type BatchCreate struct {
	Amount       int       `json:"Amount"`       // Amount of drugs in the batch
	Description  string    `json:"Description"`  // Drug description
	Manufacturer string    `json:"Manufacturer"` // Manufacturer name
	Name         string    `json:"Name"`         // Drug name
	ExpiryDate   time.Time `json:"ExpiryDate"`   // Drug expiry date
	Location     string    `json:"Location"`     // Location of the drug
	Status       string    `json:"Status"`       // Drug status
}
