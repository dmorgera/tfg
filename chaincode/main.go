package main

import (
	"fmt"
	"tfg/identity"
	"tfg/registry"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	identityContract := new(identity.Contract)
	identityContract.TransactionContextHandler = new(identity.IdentityTxContext)
	identityContract.Name = "org.net.identity"
	identityContract.Info.Version = "0.0.1"

	registryContract := new(registry.Contract)
	registryContract.TransactionContextHandler = new(registry.RegistryTxContext)
	registryContract.Name = "org.net.registry"
	registryContract.Info.Version = "0.0.1"

	chaincode, err := contractapi.NewChaincode(identityContract, registryContract)
	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
	}

	chaincode.Info.Title = "AlastriaIDChaincode"
	chaincode.Info.Version = "0.0.1"

	err = chaincode.Start()
	if err != nil {
		panic(fmt.Sprintf("Error starting chaincode. %s", err.Error()))
	}

}
