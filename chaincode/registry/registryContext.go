package registry

import (
	"encoding/hex"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the credential
type RegistryTxContextInterface interface {
	contractapi.TransactionContextInterface
	GetCredentialList() CredentialListInterface
	Keccak256(input string) string
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// credential contract
type RegistryTxContext struct {
	contractapi.TransactionContext
	credentialList *list
}

// GetCredentialList return credential list
func (tc *RegistryTxContext) GetCredentialList() CredentialListInterface {
	if tc.credentialList == nil {
		tc.credentialList = newCredentialList(tc)
	}

	return tc.credentialList
}

func (tc *RegistryTxContext) Keccak256(input string) string {
	hash := solsha3.SoliditySHA3(
		// types
		[]string{"string"},
		// values
		[]interface{}{
			input,
		},
	)
	return hex.EncodeToString(hash)
}
