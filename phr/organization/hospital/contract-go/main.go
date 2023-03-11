/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/Ha-youngPark/phr-trading-system/phr/organization/hospital/contract-go/phr"
)

func main() {

	contract := new(phr.Contract)
	contract.TransactionContextHandler = new(phr.TransactionContext)
	contract.Name = "org.phrnet.phrlist"
	contract.Info.Version = "0.0.1"

	chaincode, err := contractapi.NewChaincode(contract)

	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
	}

	chaincode.Info.Title = "PHRChaincode"
	chaincode.Info.Version = "0.0.1"

	err = chaincode.Start()

	if err != nil {
		panic(fmt.Sprintf("Error starting chaincode. %s", err.Error()))
	}
}
