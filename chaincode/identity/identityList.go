package identity

import (
	ledgerapi "tfg/ledgerapi"
)

// IdentityListInterface defines functionality needed
// to interact with the world state on behalf
// of an identity
type IdentityListInterface interface {
	AddIdentity(*Identity) error
	GetIdentity(string) (*Identity, error)
	UpdateIdentity(*Identity) error
}

type list struct {
	stateList ledgerapi.StateListInterface
}

func (il *list) AddIdentity(i *Identity) error {
	return il.stateList.AddState(i)
}

func (il *list) GetIdentity(idHash string) (*Identity, error) {
	i := new(Identity)

	err := il.stateList.GetState(CreateIdentityKey(idHash), i)

	if err != nil {
		return nil, err
	}

	return i, nil
}

func (il *list) UpdateIdentity(credential *Identity) error {
	return il.stateList.UpdateState(credential)
}

// NewList creates a new list from context
func newIdentityList(ctx IdentityTxContextInterface) *list {
	stateList := new(ledgerapi.StateList)
	stateList.Ctx = ctx
	stateList.Name = "org.net.idlist"
	stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*Identity))
	}

	list := new(list)
	list.stateList = stateList

	return list
}
