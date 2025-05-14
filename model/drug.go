package model

type Drug struct {
	BatchID       string `json:"BatchID"`       // Reference to Batch.ID
	ID            string `json:"ID"`            // Unique drug ID
	IsTransferred bool   `json:"isTransferred"` // Indicates if the drug has been transferred
	OwnerID       string `json:"OwnerID"`       // Current owner
	TransferID    string `json:"TransferID"`    // ID of the transfer transaction
}
