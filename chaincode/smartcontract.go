package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/AryaJayadi/SupplyChain_chaincode/model"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	drugs := []model.Drug{
		{
			ID:             uuid.New().String(),
			Name:           "Aspirin",
			Description:    "Pain reliever",
			Manufacturer:   "ABC Pharma",
			BatchID:        "B12345",
			ManufacturedAt: "2025-04-07T10:00:00Z",
			ExpiryDate:     "2027-04-07T10:00:00Z",
			Owner:          "PharmaCorp",
			Location:       "Warehouse 1",
			Status:         "Manufactured",
		},
		{
			ID:             uuid.New().String(),
			Name:           "Paracetamol",
			Description:    "Fever and pain reducer",
			Manufacturer:   "XYZ Healthcare",
			BatchID:        "B67890",
			ManufacturedAt: "2025-03-15T10:00:00Z",
			ExpiryDate:     "2027-03-15T10:00:00Z",
			Owner:          "MediSupply",
			Location:       "Warehouse 2",
			Status:         "InTransit",
		},
		{
			ID:             uuid.New().String(),
			Name:           "Ibuprofen",
			Description:    "Anti-inflammatory",
			Manufacturer:   "HealthCorp",
			BatchID:        "B13579",
			ManufacturedAt: "2025-02-25T10:00:00Z",
			ExpiryDate:     "2027-02-25T10:00:00Z",
			Owner:          "MediSupply",
			Location:       "Warehouse 3",
			Status:         "Delivered",
		},
	}

	for _, drug := range drugs {
		productJSON, err := json.Marshal(drug)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(drug.ID, productJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) CreateDrug(ctx contractapi.TransactionContextInterface, createdAt string, description string, id string, manufacturedPlace string, name string) error {
	exists, err := s.DrugExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	product := model.Drug{
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

func (s *SmartContract) ReadDrug(ctx contractapi.TransactionContextInterface, id string) (*model.Drug, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if productJSON == nil {
		return nil, fmt.Errorf("the product %s does not exist", id)
	}

	var product model.Drug
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *SmartContract) UpdateDrug(ctx contractapi.TransactionContextInterface, description string, id string, manufacturedPlace string, name string) error {
	exists, err := s.DrugExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the product %s does not exist", id)
	}

	product := model.Drug{
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

func (s *SmartContract) DeleteDrug(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.DrugExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the product %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) DrugExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return productJSON != nil, nil
}

func (s *SmartContract) GetAllDrugs(ctx contractapi.TransactionContextInterface) ([]*model.Drug, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err := resultsIterator.Close()
		if err != nil {
			fmt.Printf("failed to close iterator: %v\n", err)
		}
	}(resultsIterator)

	var drugs []*model.Drug
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product model.Drug
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		drugs = append(drugs, &product)
	}

	return drugs, nil
}
