package model

type Drug struct {
	ID          string `json:"ID"`          // Unique drug ID
	Name        string `json:"Name"`        // Drug name
	Description string `json:"Description"` // Description
	BatchID     string `json:"BatchID"`     // Reference to Batch.ID
	Owner       string `json:"Owner"`       // Current owner
	Location    string `json:"Location"`    // Current location
	Status      string `json:"Status"`      // e.g., Manufactured, InTransit, Delivered, Expired
}
