package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Product struct {
	CreatedAt         string `json:"CreatedAt"`         // e.g., "2025-04-07T10:00:00Z"
	Description       string `json:"Description"`       // Product description
	ID                string `json:"ID"`                // Unique ID
	ManufacturedPlace string `json:"ManufacturedPlace"` // Where it was manufactured
	Name              string `json:"Name"`              // Product name
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	products := []Product{
		{
			Name:              "Wireless Mouse",
			CreatedAt:         "2025-04-01T09:00:00Z",
			Description:       "Ergonomic wireless mouse with USB receiver",
			ManufacturedPlace: "Shenzhen, China",
		},
		{
			Name:              "Mechanical Keyboard",
			CreatedAt:         "2025-04-01T09:15:00Z",
			Description:       "RGB backlit mechanical keyboard with blue switches",
			ManufacturedPlace: "Taipei, Taiwan",
		},
		{
			Name:              "Smartphone Stand",
			CreatedAt:         "2025-04-01T09:30:00Z",
			Description:       "Adjustable aluminum stand for mobile phones",
			ManufacturedPlace: "Jakarta, Indonesia",
		},
		{
			Name:              "USB-C Hub",
			CreatedAt:         "2025-04-01T09:45:00Z",
			Description:       "6-in-1 USB-C hub with HDMI and Ethernet ports",
			ManufacturedPlace: "Seoul, South Korea",
		},
		{
			Name:              "Noise-Cancelling Headphones",
			CreatedAt:         "2025-04-01T10:00:00Z",
			Description:       "Wireless over-ear headphones with ANC",
			ManufacturedPlace: "Tokyo, Japan",
		},
		{
			Name:              "4K Webcam",
			CreatedAt:         "2025-04-01T10:15:00Z",
			Description:       "Ultra HD webcam with built-in microphone",
			ManufacturedPlace: "Hanoi, Vietnam",
		},
		{
			Name:              "External SSD",
			CreatedAt:         "2025-04-01T10:30:00Z",
			Description:       "1TB portable SSD with USB 3.2",
			ManufacturedPlace: "Bangkok, Thailand",
		},
		{
			Name:              "Gaming Chair",
			CreatedAt:         "2025-04-01T10:45:00Z",
			Description:       "Ergonomic chair with lumbar support and tilt lock",
			ManufacturedPlace: "Kuala Lumpur, Malaysia",
		},
		{
			Name:              "Smartwatch",
			CreatedAt:         "2025-04-01T11:00:00Z",
			Description:       "Fitness-focused smartwatch with GPS and heart rate",
			ManufacturedPlace: "New Delhi, India",
		},
		{
			Name:              "Portable Projector",
			CreatedAt:         "2025-04-01T11:15:00Z",
			Description:       "Mini HD projector with Wi-Fi and Bluetooth",
			ManufacturedPlace: "Manila, Philippines",
		},
	}

	for i := range products {
		products[i].ID = uuid.New().String()

		productJSON, err := json.Marshal(products[i])
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(products[i].ID, productJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, createdAt string, description string, id string, manufacturedPlace string, name string) error {
	exists, err := s.ProductExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	product := Product{
		CreatedAt:         createdAt,
		Description:       description,
		ID:                id,
		ManufacturedPlace: manufacturedPlace,
		Name:              name,
	}
	productJson, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, productJson)
}

func (s *SmartContract) ReadProduct(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if productJSON == nil {
		return nil, fmt.Errorf("the product %s does not exist", id)
	}

	var product Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *SmartContract) UpdateProduct(ctx contractapi.TransactionContextInterface, description string, id string, manufacturedPlace string, name string) error {
	exists, err := s.ProductExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the product %s does not exist", id)
	}

	product := Product{
		Description:       description,
		ID:                id,
		ManufacturedPlace: manufacturedPlace,
		Name:              name,
	}
	productJson, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, productJson)
}

func (s *SmartContract) DeleteProduct(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.ProductExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the product %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) ProductExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return productJSON != nil, nil
}
