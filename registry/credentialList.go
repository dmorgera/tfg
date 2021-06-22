package registry

import (
	ledgerapi "tfg/ledgerapi"
)

// CredentialListInterface defines functionality needed
// to interact with the world state on behalf
// of a credential
type CredentialListInterface interface {
	AddCredential(*Credential) error
	GetCredential(string, string) (*Credential, error)
	UpdateCredential(*Credential) error
}

type list struct {
	stateList ledgerapi.StateListInterface
}

func (cl *list) AddCredential(cred *Credential) error {
	return cl.stateList.AddState(cred)
}

func (cl *list) GetCredential(credentialType string, hash string) (*Credential, error) {
	c := new(Credential)

	err := cl.stateList.GetState(CreateCredentialKey(credentialType, hash), c)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (cl *list) UpdateCredential(credential *Credential) error {
	return cl.stateList.UpdateState(credential)
}

func (cl *list) GetCredentialList() {
}

// NewList creates a new list from context
func newCredentialList(ctx RegistryTxContextInterface) *list {
	stateList := new(ledgerapi.StateList)
	stateList.Ctx = ctx
	stateList.Name = "org.net.credlist"
	stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*Credential))
	}

	list := new(list)
	list.stateList = stateList

	return list
}
