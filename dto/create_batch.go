package dto

import (
	"time"
)

type CreateBatch struct {
	Amount         int       `json:"Amount"`         // Amount of drugs in the batch
	DrugName       string    `json:"DrugName"`       // Drug name
	ExpiryDate     time.Time `json:"ExpiryDate"`     // Expiry date for all drugs in the batch
	ID             string    `json:"ID"`             // Unique batch ID
	ProductionDate time.Time `json:"ProductionDate"` // Production date for all drugs in the batch
}
