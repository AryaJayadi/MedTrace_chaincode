package model

import (
	"time"
)

type Batch struct {
	ID                  string    `json:"ID"`                  // Unique batch ID
	ManufacturerName    string    `json:"ManufacturerName"`    // Manufacturer name
	ManufactureLocation string    `json:"ManufactureLocation"` // Manufacture timestamp
	ProductionDate      time.Time `json:"ProductionDate"`      // Production date
	ExpiryDate          time.Time `json:"ExpiryDate"`          // Expiry date for all drugs in the batch
	DrugName            string    `json:"DrugName"`            // Drug name
}
