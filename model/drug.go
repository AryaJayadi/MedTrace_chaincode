package model

type Drug struct {
	ID             string `json:"ID"`             // Unique drug ID
	Name           string `json:"Name"`           // Drug name
	Description    string `json:"Description"`    // Description of the drug
	Manufacturer   string `json:"Manufacturer"`   // Manufacturer name
	BatchID        string `json:"BatchID"`        // Batch ID
	ManufacturedAt string `json:"ManufacturedAt"` // Manufacture timestamp, e.g., "2025-04-07T10:00:00Z"
	ExpiryDate     string `json:"ExpiryDate"`     // Drug expiry date
	Owner          string `json:"Owner"`          // Current owner (useful in chain of custody)
	Location       string `json:"Location"`       // physical or logical location
	Status         string `json:"Status"`         // e.g., Manufactured, InTransit, Delivered, Expired
}
