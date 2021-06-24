package identity

import (
	"encoding/hex"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the identity
type IdentityTxContextInterface interface {
	contractapi.TransactionContextInterface
	GetIdentityList() IdentityListInterface
	Keccak256(input string) string
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// identity contract
type IdentityTxContext struct {
	contractapi.TransactionContext
	identityList *list
}

// GetIdentityList return identity list
func (tc *IdentityTxContext) GetIdentityList() IdentityListInterface {
	if tc.identityList == nil {
		tc.identityList = newIdentityList(tc)
	}

	return tc.identityList
}

func (tc *IdentityTxContext) Keccak256(input string) string {
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
