package dto

type DrugCreate struct {
	Name        string `json:"Name"`        // Drug name
	Description string `json:"Description"` // Drug description
	BatchID     string `json:"BatchID"`     // Reference to Batch.ID
	Owner       string `json:"Owner"`       // Current owner
	Location    string `json:"Location"`    // Current location
	Status      string `json:"Status"`      // Drug status
}
