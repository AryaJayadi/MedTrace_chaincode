package model

import (
	"time"
)

type Batch struct {
	DrugName            string    `json:"DrugName"`            // Drug name
	ExpiryDate          time.Time `json:"ExpiryDate"`          // Expiry date for all drugs in the batch
	ID                  string    `json:"ID"`                  // Unique batch ID
	ManufacturerName    string    `json:"ManufacturerName"`    // Manufacturer name
	ManufactureLocation string    `json:"ManufactureLocation"` // Manufacture timestamp
	ProductionDate      time.Time `json:"ProductionDate"`      // Production date
}
