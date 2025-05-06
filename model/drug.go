package model

type Drug struct {
	BatchID    string `json:"BatchID"`    // Reference to Batch.ID
	ID         string `json:"ID"`         // Unique drug ID
	OwnerID    string `json:"OwnerID"`    // Current owner
	TransferID string `json:"TransferID"` // ID of the transfer transaction
}
