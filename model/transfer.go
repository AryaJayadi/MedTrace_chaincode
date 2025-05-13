package model

import "time"

type Transfer struct {
	ID           *string    `json:"ID"`           // Unique transfer ID
	IsAccepted   *bool      `json:"isAccepted"`   // null, true, false
	ReceiveDate  *time.Time `json:"ReceiveDate"`  // Receive date
	ReceiverID   *string    `json:"ReceiverID"`   // Receiver ID
	SenderID     *string    `json:"SenderID"`     // Sender ID
	TransferDate *time.Time `json:"TransferDate"` // Transfer date
}
