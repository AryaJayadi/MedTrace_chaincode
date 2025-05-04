package model

import "time"

type Transfer struct {
	ID           string    `json:"ID"`           // Unique transfer ID
	TransferDate time.Time `json:"TransferDate"` // Transfer date
	ReceiveDate  time.Time `json:"ReceiveDate"`  // Receive date
	IsAccepted   bool      `json:"isAccepted"`   // null, true, false
}
