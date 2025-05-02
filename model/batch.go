package model

type Batch struct {
	ID             string   `json:"ID"`             // Unique batch ID
	Manufacturer   string   `json:"Manufacturer"`   // Manufacturer name
	ManufacturedAt string   `json:"ManufacturedAt"` // Manufacture timestamp
	ExpiryDate     string   `json:"ExpiryDate"`     // Expiry date for all drugs in the batch
	Drugs          []string `json:"Drugs"`          // List of Drug IDs in this batch
}
