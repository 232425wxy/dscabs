package main

import (
	"log"

	"github.com/232425wxy/dscabs/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	gateKeeper, err := contractapi.NewChaincode(&chaincode.DSCABS{})
	if err != nil {
		log.Panicf("Error creating DSCABS chaincode: %v", err)
	}

	if err := gateKeeper.Start(); err != nil {
		log.Panicf("Error starting DSCABS chaincode: %v", err)
	}
}
