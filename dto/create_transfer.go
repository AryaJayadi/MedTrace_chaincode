package dto

type CreateTransfer struct {
	OrganizationID string         `json:"organizationId"` // Organization ID
	BatchAmounts   map[string]int `json:"batchAmounts"`   // Amount of drugs in the batch [ID, Amount]
}
