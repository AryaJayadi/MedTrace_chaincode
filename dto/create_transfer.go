package dto

import "time"

type CreateTransfer struct {
	DrugsID      []*string  `json:"DrugsID"`      // List of drug IDs
	ReceiverID   *string    `json:"ReceiverID"`   // Receiver ID
	SenderID     *string    `json:"SenderID"`     // Sender ID
	TransferDate *time.Time `json:"TransferDate"` // Transfer date
}
