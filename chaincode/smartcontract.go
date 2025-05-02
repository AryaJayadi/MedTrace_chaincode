package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/AryaJayadi/SupplyChain_chaincode/dto"
	"github.com/AryaJayadi/SupplyChain_chaincode/model"

	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	drugs := []model.Drug{}

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

func (s *SmartContract) CreateBatch(ctx contractapi.TransactionContextInterface, param dto.BatchCreate) error {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSP ID: %v", err)
	}

	batchID := uuid.NewString()

	drugCreate := dto.DrugCreate{
		Name:        param.Name,
		Description: param.Description,
		BatchID:     batchID,
		Owner:       mspID,
		Location:    param.Location,
		Status:      param.Status,
	}

	var drugs []string

	for i := 0; i < param.Amount; i++ {
		drugID, err := s.CreateDrug(ctx, drugCreate)
		if err != nil {
			return fmt.Errorf("failed to create drug: %v", err)
		}
		drugs = append(drugs, drugID)
	}

	batch := model.Batch{
		ID:             batchID,
		Manufacturer:   param.Manufacturer,
		ManufacturedAt: time.Now(),
		ExpiryDate:     param.ExpiryDate,
		Drugs:          drugs,
	}

	batchJSON, err := json.Marshal(batch)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(batchID, batchJSON)
}

func (s *SmartContract) CreateDrug(ctx contractapi.TransactionContextInterface, param dto.DrugCreate) (string, error) {
	DrugID := uuid.NewString()

	drug := model.Drug{
		ID:          DrugID,
		Name:        param.Name,
		Description: param.Description,
		BatchID:     param.BatchID,
		Owner:       param.Owner,
		Location:    param.Location,
		Status:      param.Status,
	}

	drugJSON, err := json.Marshal(drug)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(DrugID, drugJSON)
	if err != nil {
		return "", err
	}

	return DrugID, nil
}

func (s *SmartContract) BatchExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	batchJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return batchJSON != nil, nil
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
