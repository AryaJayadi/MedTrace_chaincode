package model

type Drug struct {
	ID         string `json:"ID"`         // Unique drug ID
	BatchID    string `json:"BatchID"`    // Reference to Batch.ID
	OwnerID    string `json:"OwnerID"`    // Current owner
	TransferID string `json:"TransferID"` // ID of the transfer transaction
}
