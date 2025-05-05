package main

import (
	"log"

	"github.com/AryaJayadi/MedTrace_chaincode/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func main() {
	drugChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating drug chaincode: %s", err)
	}

	if err := drugChaincode.Start(); err != nil {
		log.Panicf("Error starting drug chaincode: %s", err)
	}
}
