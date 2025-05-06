package dto

import "time"

type UpdateBatch struct {
	DrugName       string    `json:"DrugName"`       // Drug name
	ExpiryDate     string    `json:"ExpiryDate"`     // Expiry date for all drugs in the batch
	ID             string    `json:"ID"`             // Unique batch ID
	ProductionDate time.Time `json:"ProductionDate"` // Production date
}
