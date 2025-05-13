package dto

type CreateTransfer struct {
	SenderOrgID   string   `json:"SenderOrgID"`   // Sender organization ID
	ReceiverOrgID string   `json:"ReceiverOrgID"` // Receiver organization ID
	DrugsID       []string `json:"DrugsID"`       // List of drug IDs
}
