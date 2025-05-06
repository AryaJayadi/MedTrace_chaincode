package model

type Organization struct {
	ID       string `json:"ID"`       // Unique organization ID
	Location string `json:"Location"` // Organization location
	Name     string `json:"Name"`     // Organization name
	Type     string `json:"Type"`     // Organization type (e.g., Manufacturer, Distributor, Pharmacy)
}
