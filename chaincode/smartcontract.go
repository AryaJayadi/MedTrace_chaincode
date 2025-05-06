package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/AryaJayadi/MedTrace_chaincode/dto"
	"github.com/AryaJayadi/MedTrace_chaincode/model"

	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	organizations := []model.Organization{
		{
			ID:       "Org1",
			Location: "Switzerland",
			Name:     "PharmaCorp",
			Type:     "Manufacturer",
		},
		{
			ID:       "Org2",
			Location: "Indonesia",
			Name:     "SehatDistribusi",
			Type:     "Distributor",
		},
		{
			ID:       "Org3",
			Location: "Indonesia",
			Name:     "ApotekSehat",
			Type:     "Pharmacy",
		},
	}

	for _, org := range organizations {
		orgJSON, err := json.Marshal(org)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(org.ID, orgJSON)
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
	for range param.Amount {
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

	key := string(model.PrefixBatch) + batchID
	return ctx.GetStub().PutState(key, batchJSON)
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

	key := string(model.PrefixDrug) + DrugID
	err = ctx.GetStub().PutState(key, drugJSON)
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

func (s *SmartContract) GetAllBatches(ctx contractapi.TransactionContextInterface) ([]*model.Batch, error) {
	startKey := string(model.PrefixBatch)
	endKey := string(model.PrefixBatch) + "~"
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err := resultsIterator.Close()
		if err != nil {
			fmt.Printf("failed to close iterator: %v\n", err)
		}
	}(resultsIterator)

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
