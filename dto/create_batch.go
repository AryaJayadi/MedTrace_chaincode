package dto

type BatchCreate struct {
	Description  string `json:"Description"`  // Drug description
	Manufacturer string `json:"Manufacturer"` // Manufacturer name
	Name         string `json:"Name"`         // Drug name
	ExpiryDate   string `json:"ExpiryDate"`   // Drug expiry date
	Location     string `json:"Location"`     // Location of the drug
	Status       string `json:"Status"`       // Drug status
}
