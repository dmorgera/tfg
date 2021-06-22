package registry

import (
	"encoding/hex"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

// Contract chaincode that defines
// the business logic for managing credentials
type Contract struct {
	contractapi.Contract
}

func (c *Contract) AddSubjectCredential(ctx RegistryTxContextInterface, subjectCredentialHash string, uri string) error {
	id, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	credential := Credential{PSMHash: subjectCredentialHash, Status: ValidC, Type: "subject", URI: uri, SubjectID: keccak256(id)}

	return ctx.GetCredentialList().AddCredential(&credential)
}

func (c *Contract) AddIssuerCredential(ctx RegistryTxContextInterface, issuerCredentialHash string) error {
	id, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	credential := Credential{PSMHash: issuerCredentialHash, Status: ValidC, Type: "issuer", SubjectID: keccak256(id)}

	return ctx.GetCredentialList().AddCredential(&credential)
}

func (c *Contract) DeleteSubjectCredential(ctx RegistryTxContextInterface, subjectCredentialHash string) error {
	credential, err := ctx.GetCredentialList().GetCredential("subject", subjectCredentialHash)
	if err != nil {
		return err
	}
	if credential.Status != DeletedBySubject {
		credential.Status = DeletedBySubject
	}

	err = ctx.GetCredentialList().UpdateCredential(credential)
	if err != nil {
		return err
	}

	return nil
}

func (c *Contract) GetSubjectCredentialStatus(ctx RegistryTxContextInterface, idHash string, subjectCredentialHash string) (CredentialStatus, error) {
	credential, err := ctx.GetCredentialList().GetCredential("subject", subjectCredentialHash)
	if err != nil {
		return 4, err
	}
	if credential == nil || credential.SubjectID != idHash {
		return 4, fmt.Errorf("Credential with hash %s doesn't exist for that subject", subjectCredentialHash)
	}
	return credential.Status, nil
}

func (c *Contract) GetIssuerCredentialStatus(ctx RegistryTxContextInterface, idHash string, issuerCredentialHash string) (CredentialStatus, error) {
	credential, err := ctx.GetCredentialList().GetCredential("issuer", issuerCredentialHash)
	if err != nil {
		return 4, err
	}
	if credential == nil || credential.IssuerID != idHash {
		return 4, fmt.Errorf("Credential with hash %s doesn't exist for that issuer", issuerCredentialHash)
	}
	return credential.Status, nil
}

func (c *Contract) UpdateCredentialStatus(ctx RegistryTxContextInterface, issuerCredentialHash string, status CredentialStatus) error {
	credential, err := ctx.GetCredentialList().GetCredential("issuer", issuerCredentialHash)
	if err != nil {
		return err
	}
	if status > credential.Status {
		if status == AskIssuer || status == Revoked {
			credential.Status = status
		}
	}

	err = ctx.GetCredentialList().UpdateCredential(credential)
	if err != nil {
		return err
	}

	return nil
}

func (c *Contract) GetSubjectCredentialList(ctx RegistryTxContextInterface, idHash string) ([]string, error) {
	credIterator, err := ctx.GetStub().GetStateByRange("", "")

}

func keccak256(input string) string {
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
