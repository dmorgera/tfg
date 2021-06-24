package identity

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

const timeToLive int64 = 10000

func (c *Contract) PrepareAlastriaID(ctx IdentityTxContextInterface, clientId string) error {
	// TODO: When entity contracts are created, we have to check if caller is an identityIssuer.
	// callerIdHash := ctx.Keccak256(ctx.GetClientIdentity().GetID())
	// caller, err := ctx.GetEntityList().GetEntity(issuerIdHash)
	// if err != nil {
	// 	return err
	// }
	// if !caller.IsIssuer {
	// 	return nil
	// }

	identity := Identity{IDHash: ctx.Keccak256(clientId), Pending: true, TTL: time.Now().Unix() + timeToLive}
	return ctx.GetIdentityList().AddIdentity(&identity)
}

func (c *Contract) CreateAlastriaIdentity(ctx IdentityTxContextInterface) error {
	clientId, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return err
	}
	idHash := ctx.Keccak256(clientId)
	caller, err := ctx.GetIdentityList().GetIdentity(idHash)
	if err != nil {
		return err
	}

	if caller == nil || !caller.Pending || caller.TTL < time.Now().Unix() {
		return fmt.Errorf("Error creating ID")
	}
	caller.Pending = false

	return ctx.GetIdentityList().UpdateIdentity(caller)
}
