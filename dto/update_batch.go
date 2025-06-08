package dto

import "time"

type UpdateBatch struct {
	DrugName       string    `json:"DrugName"`       // Drug name
	ExpiryDate     time.Time `json:"ExpiryDate"`     // Expiry date for all drugs in the batch
	ProductionDate time.Time `json:"ProductionDate"` // Production date
}
