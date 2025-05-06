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

const batchDrugIndex = "batch~drug"

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

func (s *SmartContract) CreateDrug(ctx contractapi.TransactionContextInterface, ownerID string, batchID string) (string, error) {
	drugID := uuid.NewString()

	drug := model.Drug{
		BatchID: batchID,
		ID:      drugID,
		OwnerID: ownerID,
	}

	drugJSON, err := json.Marshal(drug)
	if err != nil {
		return "", fmt.Errorf("failed to marshal drug: %v", err)
	}

	err = ctx.GetStub().PutState(drugID, drugJSON)
	if err != nil {
		return "", fmt.Errorf("failed to put drug to world state: %v", err)
	}

	batchDrugIndexKey, err := ctx.GetStub().CreateCompositeKey(batchDrugIndex, []string{batchID, drugID})
	if err != nil {
		return "", fmt.Errorf("failed to create composite key: %v", err)
	}

	value := []byte{0x00}
	err = ctx.GetStub().PutState(batchDrugIndexKey, value)
	if err != nil {
		return "", fmt.Errorf("failed to put batch-drug index to world state: %v", err)
	}

	return drugID, nil
}

func (s *SmartContract) GetOrganization(ctx contractapi.TransactionContextInterface, id string) (*model.Organization, error) {
	orgJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if orgJSON == nil {
		return nil, fmt.Errorf("organization %s does not exist", id)
	}

	var org model.Organization
	err = json.Unmarshal(orgJSON, &org)
	if err != nil {
		return nil, err
	}

	return &org, nil
}

func (s *SmartContract) GetAllOrganizations(ctx contractapi.TransactionContextInterface) ([]*model.Organization, error) {
	resIterator, err := ctx.GetStub().GetStateByRange("Org", "Org~")
	if err != nil {
		return nil, err
	}
	defer resIterator.Close()

	var orgs []*model.Organization
	for resIterator.HasNext() {
		res, err := resIterator.Next()
		if err != nil {
			return nil, err
		}

		var org model.Organization
		err = json.Unmarshal(res.Value, &org)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, &org)
	}

	return orgs, nil
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
