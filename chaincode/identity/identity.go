package identity

import (
	"encoding/json"
	"fmt"
	ledgerapi "tfg/ledgerapi"
)

type Identity struct {
	Pending bool   `json:"pending"`
	IDHash  string `json:"idHash"`
	TTL     int64  `json:"ttl"`
}

// CreateIdentity creates a key for credentials
func CreateIdentityKey(idHash string) string {
	return ledgerapi.MakeKey(idHash)
}

// GetSplitKey returns values which should be used to form key
func (i *Identity) GetSplitKey() []string {
	return []string{i.IDHash}
}

// Serialize formats the credential as JSON bytes
func (i *Identity) Serialize() ([]byte, error) {
	return json.Marshal(i)
}

// Deserialize formats the credential from JSON bytes
func Deserialize(bytes []byte, i *Identity) error {
	err := json.Unmarshal(bytes, i)

	if err != nil {
		return fmt.Errorf("Error deserializing credential. %s", err.Error())
	}

	return nil
}
