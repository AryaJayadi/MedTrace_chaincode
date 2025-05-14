package dto

import "time"

type ProcessTransfer struct {
	ReceiveDate *time.Time `json:"ReceiveDate"` // Receive date
	TransferID  string     `json:"transferID"`  // ID of Transfer to be processed
}
