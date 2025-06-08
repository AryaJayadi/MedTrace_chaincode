package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AryaJayadi/MedTrace_chaincode/dto"
	"github.com/AryaJayadi/MedTrace_chaincode/model"

	"github.com/hyperledger/fabric-chaincode-go/v2/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const (
	ownerDrugIndex        = "owner~drug"
	batchDrugIndex        = "batch~drug"
	senderTransferIndex   = "sender~transfer"
	receiverTransferIndex = "receiver~transfer"
	transferDrugIndex     = "transfer~drug"
)

const (
	batchKey    = "B"
	transferKey = "T"
	drugKey     = "D"
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
		{
			ID:       "Org4",
			Location: "Indonesia",
			Name:     "Pasien",
			Type:     "Patient",
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

func (s *SmartContract) updateDrugOwner(ctx contractapi.TransactionContextInterface, drug *model.Drug, newOwnerID string) (*string, error) {
	ownerIndexKey, err := ctx.GetStub().CreateCompositeKey(ownerDrugIndex, []string{drug.OwnerID, drug.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %w", err)
	}

	if err := ctx.GetStub().DelState(ownerIndexKey); err != nil {
		return nil, fmt.Errorf("failed to delete old owner-drug index from world state: %w", err)
	}

	drug.OwnerID = newOwnerID

	value := []byte{0x00}
	ownerDrugIndexKey, err := ctx.GetStub().CreateCompositeKey(ownerDrugIndex, []string{newOwnerID, drug.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %w", err)
	}
	if err := ctx.GetStub().PutState(ownerDrugIndexKey, value); err != nil {
		return nil, fmt.Errorf("failed to put owner-drug index to world state: %w", err)
	}

	return &drug.ID, nil
}

func (s *SmartContract) updateDrugTransfer(ctx contractapi.TransactionContextInterface, drug *model.Drug, transferID string) (*string, error) {
	transferDrugIndexKey, err := ctx.GetStub().CreateCompositeKey(transferDrugIndex, []string{drug.TransferID, drug.ID})
	if err != nil {
		return nil, err
	}

	if err := ctx.GetStub().DelState(transferDrugIndexKey); err != nil {
		return nil, fmt.Errorf("failed to delete old transfer-drug index from world state: %w", err)
	}

	if transferID == "" {
		return &drug.ID, nil
	}

	drug.TransferID = transferID

	value := []byte{0x00}
	transferDrugIndexKey, err = ctx.GetStub().CreateCompositeKey(transferDrugIndex, []string{transferID, drug.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key: %w", err)
	}
	if err := ctx.GetStub().PutState(transferDrugIndexKey, value); err != nil {
		return nil, fmt.Errorf("failed to put transfer-drug index to world state: %w", err)
	}

	return &drug.ID, nil
}

func (s *SmartContract) CreateDrug(ctx contractapi.TransactionContextInterface, org *model.Organization, batchID string, drugID string) (string, error) {
	drug := model.Drug{
		BatchID:  batchID,
		ID:       drugID,
		Location: org.Location,
		OwnerID:  org.ID,
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
	ownderDrugIndexKey, err := ctx.GetStub().CreateCompositeKey(ownerDrugIndex, []string{org.ID, drug.ID})
	if err != nil {
		return "", fmt.Errorf("failed to create composite key: %w", err)
	}

	value := []byte{0x00}
	err = ctx.GetStub().PutState(batchDrugIndexKey, value)
	if err != nil {
		return "", fmt.Errorf("failed to put batch-drug index to world state: %v", err)
	}
	if err := ctx.GetStub().PutState(ownderDrugIndexKey, value); err != nil {
		return "", fmt.Errorf("failed to put owner-drug index to world state: %w", err)
	}

	return drugID, nil
}

func (s *SmartContract) GetDrug(ctx contractapi.TransactionContextInterface, drugID string) (*model.Drug, error) {
	drugJSON, err := ctx.GetStub().GetState(drugID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if drugJSON == nil {
		return nil, fmt.Errorf("drug %s does not exist", drugID)
	}

	var drug model.Drug
	err = json.Unmarshal(drugJSON, &drug)
	if err != nil {
		return nil, err
	}

	return &drug, nil
}

type drugFilter func(drug *model.Drug, org *model.Organization) bool

func (s *SmartContract) getFilteredDrugs(ctx contractapi.TransactionContextInterface, filter drugFilter) ([]*model.Drug, error) {
	org, err := s.getOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization ID: %w", err)
	}

	drugsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(ownerDrugIndex, []string{org.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get drugs: %w", err)
	}
	defer drugsIterator.Close()

	drugs := make([]*model.Drug, 0)
	for drugsIterator.HasNext() {
		responseRange, err := drugsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate drugs: %w", err)
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %w", err)
		}

		if len(compositeKeyParts) > 1 {
			returnedDrugID := compositeKeyParts[1]
			drug, err := s.GetDrug(ctx, returnedDrugID)
			if err != nil {
				return nil, fmt.Errorf("failed to get drug: %w", err)
			}

			if filter(drug, org) {
				drugs = append(drugs, drug)
			}
		}
	}

	return drugs, nil
}

func (s *SmartContract) GetMyDrug(ctx contractapi.TransactionContextInterface) ([]*model.Drug, error) {
	return s.getFilteredDrugs(ctx, func(drug *model.Drug, org *model.Organization) bool {
		return true
	})
}

func (s *SmartContract) GetMyAvailDrugs(ctx contractapi.TransactionContextInterface) ([]*model.Drug, error) {
	return s.getFilteredDrugs(ctx, func(drug *model.Drug, org *model.Organization) bool {
		return !drug.IsTransferred
	})
}

func (s *SmartContract) GetDrugByBatch(ctx contractapi.TransactionContextInterface, batchID string) ([]*model.Drug, error) {
	drugsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(batchDrugIndex, []string{batchID})
	if err != nil {
		return nil, fmt.Errorf("failed to get drugs: %w", err)
	}
	defer drugsIterator.Close()

	drugs := make([]*model.Drug, 0)
	for drugsIterator.HasNext() {
		responseRange, err := drugsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate drugs: %w", err)
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %w", err)
		}

		if len(compositeKeyParts) > 1 {
			returnedDrugID := compositeKeyParts[1]
			drug, err := s.GetDrug(ctx, returnedDrugID)
			if err != nil {
				return nil, fmt.Errorf("failed to get drug: %w", err)
			}

			drugs = append(drugs, drug)
		}
	}

	return drugs, nil
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

func (s *SmartContract) CreateTransfer(ctx contractapi.TransactionContextInterface, req string) (*model.Transfer, error) {
	org, err := s.getOrg(ctx)
	if err != nil {
		return nil, err
	}

	var createTransfer dto.CreateTransfer
	err = json.Unmarshal([]byte(req), &createTransfer)
	if err != nil {
		return nil, err
	}

	transferID, _, err := s.generateModelId(ctx, transferKey)
	if err != nil {
		return nil, err
	}

	isAccepted := false
	transfer := model.Transfer{
		ID:         transferID,
		IsAccepted: isAccepted,
		// ReceiveDate:  nil,
		ReceiverID:   *createTransfer.ReceiverID,
		SenderID:     org.ID,
		TransferDate: *createTransfer.TransferDate,
	}
	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return nil, err
	}

	if err := ctx.GetStub().PutState(transferID, transferJSON); err != nil {
		return nil, err
	}

	value := []byte{0x00}
	senderTransferIndexKey, err := ctx.GetStub().CreateCompositeKey(senderTransferIndex, []string{org.ID, transferID})
	if err != nil {
		return nil, err
	}
	receiverTransferIndexKey, err := ctx.GetStub().CreateCompositeKey(receiverTransferIndex, []string{*createTransfer.ReceiverID, transferID})
	if err != nil {
		return nil, err
	}
	if err := ctx.GetStub().PutState(senderTransferIndexKey, value); err != nil {
		return nil, err
	}
	if err := ctx.GetStub().PutState(receiverTransferIndexKey, value); err != nil {
		return nil, err
	}

	for _, drugID := range createTransfer.DrugsID {
		drug, err := s.GetDrug(ctx, *drugID)
		if err != nil {
			return nil, err
		}

		if drug.IsTransferred {
			return nil, fmt.Errorf("drug %s has already been transferred", *drugID)
		}

		if drug.OwnerID != org.ID {
			return nil, fmt.Errorf("drug %s does not belong to the sender", *drugID)
		}

		drug.IsTransferred = true

		drugJSON, err := json.Marshal(drug)
		if err != nil {
			return nil, err
		}
		if err := ctx.GetStub().PutState(*drugID, drugJSON); err != nil {
			return nil, err
		}

		transferDrugIndexKey, err := ctx.GetStub().CreateCompositeKey(transferDrugIndex, []string{transferID, *drugID})
		if err != nil {
			return nil, err
		}
		if err := ctx.GetStub().PutState(transferDrugIndexKey, value); err != nil {
			return nil, err
		}
	}
	log.Printf("Drugs transferred: %v\n", createTransfer.DrugsID)

	return &transfer, nil
}

func (s *SmartContract) GetTransfer(ctx contractapi.TransactionContextInterface, id string) (*model.Transfer, error) {
	transferJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if transferJSON == nil {
		return nil, fmt.Errorf("transfer %s does not exist", id)
	}

	var transfer model.Transfer
	if err := json.Unmarshal(transferJSON, &transfer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transfer: %w", err)
	}

	return &transfer, nil
}

func (s *SmartContract) GetMyOutTransfer(ctx contractapi.TransactionContextInterface) ([]*model.Transfer, error) {
	return s.getMyTransfer(ctx, false)
}

func (s *SmartContract) GetMyInTransfer(ctx contractapi.TransactionContextInterface) ([]*model.Transfer, error) {
	return s.getMyTransfer(ctx, true)
}

func (s *SmartContract) GetMyTransfers(ctx contractapi.TransactionContextInterface) ([]*model.Transfer, error) {
	outTransfers, err := s.getMyTransfer(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get out transfers: %w", err)
	}

	inTransfers, err := s.getMyTransfer(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get in transfers: %w", err)
	}

	return append(outTransfers, inTransfers...), nil
}

func (s *SmartContract) getMyTransfer(ctx contractapi.TransactionContextInterface, isIn bool) ([]*model.Transfer, error) {
	org, err := s.getOrg(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization ID: %w", err)
	}

	var transferIndex string
	if isIn {
		transferIndex = receiverTransferIndex
	} else {
		transferIndex = senderTransferIndex
	}

	transferredDrugsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(transferIndex, []string{org.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get transferred drugs: %w", err)
	}
	defer transferredDrugsIterator.Close()

	transfers := make([]*model.Transfer, 0)
	for transferredDrugsIterator.HasNext() {
		responseRange, err := transferredDrugsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate transferred drugs: %w", err)
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %w", err)
		}

		if len(compositeKeyParts) > 1 {
			returnedTransferID := compositeKeyParts[1]
			transfer, err := s.GetTransfer(ctx, returnedTransferID)
			if err != nil {
				return nil, fmt.Errorf("failed to get transfer: %w", err)
			}

			transfers = append(transfers, transfer)
		}
	}

	return transfers, nil
}

func (s *SmartContract) validateProcessTransfer(ctx contractapi.TransactionContextInterface, req string) (*model.Transfer, *model.Organization, *dto.ProcessTransfer, error) {
	org, err := s.getOrg(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get organization ID: %w", err)
	}

	var processTransfer dto.ProcessTransfer
	if err := json.Unmarshal([]byte(req), &processTransfer); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	transfer, err := s.GetTransfer(ctx, processTransfer.TransferID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get transfer: %w", err)
	}

	if org.ID != transfer.ReceiverID {
		return nil, nil, nil, fmt.Errorf("only the receiver can accept the transfer")
	}

	return transfer, org, &processTransfer, nil
}

func (s *SmartContract) AcceptTransfer(ctx contractapi.TransactionContextInterface, req string) (*model.Transfer, error) {
	transfer, org, processTransfer, err := s.validateProcessTransfer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate process transfer: %w", err)
	}

	isAccepted := true
	transfer.IsAccepted = isAccepted
	transfer.ReceiveDate = *processTransfer.ReceiveDate

	transferDrugsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(transferDrugIndex, []string{transfer.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get transferred drugs: %w", err)
	}
	defer transferDrugsIterator.Close()

	var drugsIDs []string
	for transferDrugsIterator.HasNext() {
		responseRange, err := transferDrugsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate transferred drugs: %w", err)
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %w", err)
		}

		if len(compositeKeyParts) > 1 {
			returnedDrugID := compositeKeyParts[1]
			drug, err := s.GetDrug(ctx, returnedDrugID)
			if err != nil {
				return nil, fmt.Errorf("failed to get drug: %w", err)
			}

			drug.IsTransferred = false
			drug.Location = org.Location

			_, err = s.updateDrugOwner(ctx, drug, org.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to set drug owner: %w", err)
			}

			_, err = s.updateDrugTransfer(ctx, drug, transfer.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to set drug transfer ID: %w", err)
			}

			drugJSOn, err := json.Marshal(drug)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal drug: %w", err)
			}

			if err := ctx.GetStub().PutState(drug.ID, drugJSOn); err != nil {
				return nil, fmt.Errorf("failed to put drug to world state: %w", err)
			}

			drugsIDs = append(drugsIDs, drug.ID)
		}
	}
	log.Printf("Drugs accepted: %v\n", drugsIDs)

	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer: %w", err)
	}

	if err := ctx.GetStub().PutState(transfer.ID, transferJSON); err != nil {
		return nil, fmt.Errorf("failed to put transfer to world state: %w", err)
	}

	return transfer, nil
}

func (s *SmartContract) RejectTransfer(ctx contractapi.TransactionContextInterface, req string) (*model.Transfer, error) {
	transfer, _, processTransfer, err := s.validateProcessTransfer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate process transfer: %w", err)
	}

	isAccepted := false
	transfer.IsAccepted = isAccepted
	transfer.ReceiveDate = *processTransfer.ReceiveDate

	transferDrugsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(transferDrugIndex, []string{transfer.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to get transferred drugs: %w", err)
	}
	defer transferDrugsIterator.Close()

	var drugsIDs []string
	for transferDrugsIterator.HasNext() {
		responseRange, err := transferDrugsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate transferred drugs: %w", err)
		}

		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to split composite key: %w", err)
		}

		if len(compositeKeyParts) > 1 {
			returnedDrugID := compositeKeyParts[1]
			drug, err := s.GetDrug(ctx, returnedDrugID)
			if err != nil {
				return nil, fmt.Errorf("failed to get drug: %w", err)
			}

			drug.IsTransferred = false

			_, err = s.updateDrugTransfer(ctx, drug, "")
			if err != nil {
				return nil, fmt.Errorf("failed to remove drug transfer ID: %w", err)
			}

			drugJSON, err := json.Marshal(drug)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal drug: %w", err)
			}

			if err := ctx.GetStub().PutState(drug.ID, drugJSON); err != nil {
				return nil, fmt.Errorf("failed to put drug to world state: %w", err)
			}

			drugsIDs = append(drugsIDs, drug.ID)
		}
	}
	log.Printf("Drugs rejected: %v\n", drugsIDs)

	transferJSON, err := json.Marshal(transfer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer: %w", err)
	}

	if err := ctx.GetStub().PutState(transfer.ID, transferJSON); err != nil {
		return nil, fmt.Errorf("failed to put transfer to world state: %w", err)
	}

	return transfer, nil
}

func (s *SmartContract) CreateBatch(ctx contractapi.TransactionContextInterface, req string) (*model.Batch, error) {
	org, err := s.getOrg(ctx)
	if err != nil {
		fmt.Printf("error: failed to get organization ID: %v\n", err)
		return nil, fmt.Errorf("failed to get organization ID: %v", err)
	}
	if org.Type != "Manufacturer" {
		err := fmt.Errorf("only manufacturers can create batches")
		fmt.Printf("error: %v\n", err)
		return nil, err
	}

	var createBatch dto.CreateBatch
	err = json.Unmarshal([]byte(req), &createBatch)
	if err != nil {
		fmt.Printf("error: failed to unmarshal request: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal request: %v", err)
	}

	batchID, _, err := s.generateModelId(ctx, batchKey)
	if err != nil {
		fmt.Printf("error: failed to generate batch ID: %v\n", err)
		return nil, fmt.Errorf("failed to generate batch ID: %v", err)
	}

	batch := model.Batch{
		DrugName:            createBatch.DrugName,
		ExpiryDate:          createBatch.ExpiryDate,
		ID:                  batchID,
		ManufacturerName:    org.Name,
		ManufactureLocation: org.Location,
		ProductionDate:      createBatch.ProductionDate,
	}
	batchJSON, err := json.Marshal(batch)
	if err != nil {
		fmt.Printf("error: failed to marshal batch: %v\n", err)
		return nil, fmt.Errorf("failed to marshal batch: %v", err)
	}

	err = ctx.GetStub().PutState(batch.ID, batchJSON)
	if err != nil {
		fmt.Printf("error: failed to put batch to world state: %v\n", err)
		return nil, fmt.Errorf("failed to put batch to world state: %v", err)
	}

	_, drugInt, err := s.generateModelId(ctx, drugKey)
	if err != nil {
		fmt.Printf("error: failed to generate drug ID: %v\n", err)
		return nil, fmt.Errorf("failed to generate drug ID: %v", err)
	}

	var drugsIDs []string
	for i := range createBatch.Amount {
		currDrugInt := drugInt + i
		drugID := s.formatModelId(drugKey, currDrugInt)

		drugID, err = s.CreateDrug(ctx, org, batch.ID, drugID)
		if err != nil {
			fmt.Printf("error: failed to create drug: %v\n", err)
			return nil, fmt.Errorf("failed to create drug: %v", err)
		}
		drugsIDs = append(drugsIDs, drugID)
	}
	fmt.Printf("Drugs created: %v\n", drugsIDs)

	err = s.saveModelId(ctx, drugKey, (drugInt-1)+createBatch.Amount)
	if err != nil {
		fmt.Printf("error: failed to save drug ID: %v\n", err)
		return nil, fmt.Errorf("failed to save drug ID: %v", err)
	}

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

func (s *SmartContract) UpdateBatch(ctx contractapi.TransactionContextInterface, batchID string, req string) (*model.Batch, error) {
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

	batch, err := s.GetBatch(ctx, batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch: %v", err)
	}

	batch.DrugName = updateBatch.DrugName
	batch.ExpiryDate = updateBatch.ExpiryDate
	batch.ProductionDate = updateBatch.ProductionDate

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

func (s *SmartContract) GetAllBatches(ctx contractapi.TransactionContextInterface) ([]*model.Batch, error) {
	resIterator, err := ctx.GetStub().GetStateByRange(batchKey, batchKey+"~")
	if err != nil {
		return nil, fmt.Errorf("failed to get batches: %v", err)
	}
	defer resIterator.Close()

	batches := make([]*model.Batch, 0)
	for resIterator.HasNext() {
		res, err := resIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate batches: %v", err)
		}

		var batch model.Batch
		err = json.Unmarshal(res.Value, &batch)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &batch)
	}

	return batches, nil
}

func (s *SmartContract) getOrg(ctx contractapi.TransactionContextInterface) (*model.Organization, error) {
	mspID, err := cid.GetMSPID(ctx.GetStub())
	if err != nil {
		return nil, fmt.Errorf("failed to get MSP ID: %v", err)
	}

	orgID := strings.TrimSuffix(mspID, "MSP")

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

	orgs := make([]*model.Organization, 0)
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

func (s *SmartContract) generateModelId(ctx contractapi.TransactionContextInterface, modelKey string) (string, int, error) {
	latestIDKey := fmt.Sprintf("LatestID_%s", modelKey)

	latestIDBytes, err := ctx.GetStub().GetState(latestIDKey)
	if err != nil {
		return "", -1, fmt.Errorf("failed to get latest ID: %v", err)
	}

	latestNum := 0
	if latestIDBytes != nil {
		latestNum, err = strconv.Atoi(string(latestIDBytes))
		if err != nil {
			return "", -1, fmt.Errorf("failed to parse latest ID number: %v", err)
		}
	}

	newIDNum := latestNum + 1

	err = ctx.GetStub().PutState(latestIDKey, []byte(strconv.Itoa(newIDNum)))
	if err != nil {
		return "", -1, fmt.Errorf("failed to store new latest ID: %v", err)
	}

	formattedID := s.formatModelId(modelKey, newIDNum)
	return formattedID, newIDNum, err
}

func (s *SmartContract) saveModelId(ctx contractapi.TransactionContextInterface, modelKey string, id int) error {
	latestIDKey := fmt.Sprintf("LatestID_%s", modelKey)

	return ctx.GetStub().PutState(latestIDKey, []byte(strconv.Itoa(id)))
}

func (s *SmartContract) formatModelId(modelKey string, id int) string {
	formattedID := fmt.Sprintf("%s%016d", modelKey, id)
	return formattedID
}

func (s *SmartContract) GetHistoryDrug(ctx contractapi.TransactionContextInterface, drugID string) ([]*model.HistoryDrug, error) {
	log.Printf("Getting history for drug: %s\n", drugID)

	resultIterator, err := ctx.GetStub().GetHistoryForKey(drugID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for drug %s: %v", drugID, err)
	}
	defer resultIterator.Close()

	records := make([]*model.HistoryDrug, 0)
	for resultIterator.HasNext() {
		response, err := resultIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate history for drug %s: %v", drugID, err)
		}

		var drug model.Drug
		if len(response.Value) > 0 {
			err := json.Unmarshal(response.Value, &drug)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal drug: %v", err)
			}
		} else {
			drug = model.Drug{
				ID: drugID,
			}
		}

		record := model.HistoryDrug{
			TxId:      response.TxId,
			Timestamp: response.Timestamp.AsTime(),
			Drug:      &drug,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}
