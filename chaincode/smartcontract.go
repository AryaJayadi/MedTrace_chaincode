package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/AryaJayadi/MedTrace_chaincode/dto"
	"github.com/AryaJayadi/MedTrace_chaincode/model"

	"github.com/hyperledger/fabric-chaincode-go/v2/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const (
	batchDrugIndex = "batch~drug"
)

const (
	batchKey    = "B"
	transferKey = "T"
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

func (s *SmartContract) CreateBatch(ctx contractapi.TransactionContextInterface, req string) (*model.Batch, error) {
	org, err := s.getOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization ID: %v", err)
	}
	if org.Type != "Manufacturer" {
		return nil, fmt.Errorf("only manufacturers can create batches")
	}

	var createBatch dto.CreateBatch
	err = json.Unmarshal([]byte(req), &createBatch)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %v", err)
	}

	exist, err := s.BatchExists(ctx, createBatch.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if batch exists: %v", err)
	}
	if exist {
		return nil, fmt.Errorf("batch with ID %s already exists", createBatch.ID)
	}

	ID, err := s.generateModelId(ctx, batchKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate batch ID: %v", err)
	}

	batch := model.Batch{
		DrugName:            createBatch.DrugName,
		ExpiryDate:          createBatch.ExpiryDate,
		ID:                  ID,
		ManufacturerName:    org.Name,
		ManufactureLocation: org.Location,
		ProductionDate:      time.Now(),
	}
	batchJSON, err := json.Marshal(batch)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal batch: %v", err)
	}

	err = ctx.GetStub().PutState(batch.ID, batchJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to put batch to world state: %v", err)
	}

	var drugsIDs []string
	for range createBatch.Amount {
		drugID, err := s.CreateDrug(ctx, org.ID, batch.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create drug: %v", err)
		}
		drugsIDs = append(drugsIDs, drugID)
	}
	fmt.Printf("Drugs created: %v\n", drugsIDs)

	return &batch, nil
}

func (s *SmartContract) GetBatch(ctx contractapi.TransactionContextInterface, id string) (*model.Batch, error) {
	batchJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if batchJSON == nil {
		return nil, fmt.Errorf("batch %s does not exist", id)
	}

	var batch model.Batch
	err = json.Unmarshal(batchJSON, &batch)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal batch: %v", err)
	}

	return &batch, nil
}

func (s *SmartContract) UpdateBatch(ctx contractapi.TransactionContextInterface, req string) (*model.Batch, error) {
	org, err := s.getOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization ID: %v", err)
	}
	if org.Type != "Manufacturer" {
		return nil, fmt.Errorf("only manufacturers can update batches")
	}

	var updateBatch dto.UpdateBatch
	err = json.Unmarshal([]byte(req), &updateBatch)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %v", err)
	}

	batch, err := s.GetBatch(ctx, updateBatch.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch: %v", err)
	}

	batch.DrugName = updateBatch.DrugName
	batch.ExpiryDate = updateBatch.ExpiryDate
	batch.ProductionDate = time.Now()

	batchJSON, err := json.Marshal(batch)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal batch: %v", err)
	}

	err = ctx.GetStub().PutState(batch.ID, batchJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to put batch to world state: %v", err)
	}

	return batch, nil
}

func (s *SmartContract) getOrg(ctx contractapi.TransactionContextInterface) (*model.Organization, error) {
	orgID, ok, err := cid.GetAttributeValue(ctx.GetStub(), "org.id")
	if err != nil {
		return nil, fmt.Errorf("failed to get org.id attribute: %v", err)
	}
	if !ok {
		return nil, fmt.Errorf("org.id attribute not found")
	}

	org, err := s.GetOrganization(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %v", err)
	}
	return org, nil
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

func (s *SmartContract) BatchExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	batchJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return batchJSON != nil, nil
}

func (s *SmartContract) generateModelId(ctx contractapi.TransactionContextInterface, modelKey string) (string, error) {
	// Create a unique key to store the latest ID number
	latestIDKey := fmt.Sprintf("LatestID_%s", modelKey)

	// Get the latest ID from the ledger
	latestIDBytes, err := ctx.GetStub().GetState(latestIDKey)
	if err != nil {
		return "", fmt.Errorf("failed to get latest ID: %v", err)
	}

	latestNum := 0
	if latestIDBytes != nil {
		latestNum, err = strconv.Atoi(string(latestIDBytes))
		if err != nil {
			return "", fmt.Errorf("failed to parse latest ID number: %v", err)
		}
	}

	// Increment the number
	newIDNum := latestNum + 1

	// Save the updated number back to the ledger
	err = ctx.GetStub().PutState(latestIDKey, []byte(strconv.Itoa(newIDNum)))
	if err != nil {
		return "", fmt.Errorf("failed to store new latest ID: %v", err)
	}

	// Format the final ID, e.g., B001, T002, etc.
	formattedID := fmt.Sprintf("%s%03d", modelKey, newIDNum)
	return formattedID, nil
}
