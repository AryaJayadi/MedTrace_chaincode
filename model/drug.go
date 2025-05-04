package model

import "time"

type Drug struct {
	ID           string    `json:"ID"`           // Unique drug ID
	BatchID      string    `json:"BatchID"`      // Reference to Batch.ID
	OwnerID      string    `json:"OwnerID"`      // Current owner
	TransferDate time.Time `json:"TransferDate"` // Transfer date
	ReceiveDate  time.Time `json:"ReceiveDate"`  // Receive date
}
