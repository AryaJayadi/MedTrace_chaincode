package dto

import (
	"time"
)

type BatchCreate struct {
	ID                  string    `json:"ID"`                  // Unique batch ID
	ManufactureLocation string    `json:"ManufactureLocation"` // Manufacture location
	ExpiryDate          time.Time `json:"ExpiryDate"`          // Expiry date for all drugs in the batch
	DrugName            string    `json:"DrugName"`            // Drug name
}
