package registry

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing credentials
type Contract struct {
	contractapi.Contract
}

func (c *Contract) AddSubjectCredential(ctx RegistryTxContextInterface, subjectCredentialHash string, uri string) error {
	_, err := ctx.GetCredentialList().GetCredential("subject", subjectCredentialHash)
	if err == nil {
		return fmt.Errorf("Error")
	}

	id, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	credential := Credential{PSMHash: subjectCredentialHash, Status: ValidC, Type: "subject", URI: uri, SubjectID: ctx.Keccak256(id)}
	fmt.Printf("id: %s\n", id)
	fmt.Printf("idhash: %s\n", credential.SubjectID)

	return ctx.GetCredentialList().AddCredential(&credential)
}

func (c *Contract) AddIssuerCredential(ctx RegistryTxContextInterface, issuerCredentialHash string) error {
	_, err := ctx.GetCredentialList().GetCredential("issuer", issuerCredentialHash)
	if err == nil {
		return fmt.Errorf("Error")
	}

	id, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	credential := Credential{PSMHash: issuerCredentialHash, Status: ValidC, Type: "issuer", IssuerID: ctx.Keccak256(id)}

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

// Extremely unefficient, with CouchDB this gets much better
func (c *Contract) GetSubjectCredentialList(ctx RegistryTxContextInterface, idHash string) ([]string, error) {
	credIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("org.net.credlist", []string{"subject"})
	if err != nil {
		return nil, err
	}
	defer credIterator.Close()
	var i int
	var credentials []string
	for i = 0; credIterator.HasNext(); i++ {
		responseRange, err := credIterator.Next()
		if err != nil {
			return nil, err
		}
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(responseRange.Key)
		if err != nil {
			return nil, err
		}
		hash := compositeKeyParts[1]
		c, err := ctx.GetCredentialList().GetCredential("subject", hash)
		if c.SubjectID == idHash {
			credentials = append(credentials, hash)
			fmt.Println(hash)
		}
	}
	return credentials, nil
}
