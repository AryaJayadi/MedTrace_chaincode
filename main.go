package main

import (
	"github.com/AryaJayadi/SupplyChain_chaincode/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"log"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	productChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating product chaincode: %s", err)
	}

	if err := productChaincode.Start(); err != nil {
		log.Panicf("Error starting product chaincode: %s", err)
	}
}
