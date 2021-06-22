package registry

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the credential
type RegistryTxContextInterface interface {
	contractapi.TransactionContextInterface
	GetCredentialList() CredentialListInterface
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
