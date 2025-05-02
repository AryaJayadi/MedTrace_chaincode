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
		drugJSON, err := json.Marshal(drug)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(drug.ID, drugJSON)
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

func (s *SmartContract) ReadDrug(ctx contractapi.TransactionContextInterface, id string) (*model.Drug, error) {
	drugJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if drugJSON == nil {
		return nil, fmt.Errorf("the drug %s does not exist", id)
	}

	var drug model.Drug
	err = json.Unmarshal(drugJSON, &drug)
	if err != nil {
		return nil, err
	}

	return &drug, nil
}

func (s *SmartContract) GetAllBatches(ctx contractapi.TransactionContextInterface) ([]*model.Batch, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var batches []*model.Batch
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var batch model.Batch
		err = json.Unmarshal(queryResponse.Value, &batch)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &batch)
	}

	return batches, nil
}

func (s *SmartContract) DeleteDrug(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.DrugExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the drug %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (s *SmartContract) DrugExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	drugJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return drugJSON != nil, nil
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

		var drug model.Drug
		err = json.Unmarshal(queryResponse.Value, &drug)
		if err != nil {
			return nil, err
		}
		drugs = append(drugs, &drug)
	}

	return drugs, nil
}
