package registry

import (
	"encoding/json"
	"fmt"
	ledgerapi "tfg/ledgerapi"
)

type CredentialStatus uint

const (
	ValidC CredentialStatus = iota + 1
	AskIssuer
	Revoked
	DeletedBySubject
)

func (status CredentialStatus) String() string {
	names := []string{"ValidC", "AskIssuer", "Revoked", "DeletedBySubject"}

	if status < ValidC || status > DeletedBySubject {
		return "Error"
	}

	return names[status-1]
}

// Status goes from 0 to 3
// Statuses are ValidC, AskIssuer, Revoked, DeletedBySubject
// SubjectCredential: Initially ValidC: Only DeletedBySubject
// IssuerCredentials: Initially ValidC: Only AskIssuer or Revoked, no backwards transitions.
type Credential struct {
	PSMHash   string           `json:"PSMHash"`
	Status    CredentialStatus `json:"status"`
	Type      string           `json:"type"`
	URI       string           `json:"URI"`
	IssuerID  string           `json:"issuerId"`
	SubjectID string           `json:"subjectID"`
}

// CreateCredential creates a key for credentials
func CreateCredentialKey(credentialType string, hash string) string {
	return ledgerapi.MakeKey(credentialType, hash)
}

// GetSplitKey returns values which should be used to form key
func (c *Credential) GetSplitKey() []string {
	return []string{c.Type, c.PSMHash}
}

// Serialize formats the credential as JSON bytes
func (c *Credential) Serialize() ([]byte, error) {
	return json.Marshal(c)
}

// Deserialize formats the credential from JSON bytes
func Deserialize(bytes []byte, c *Credential) error {
	err := json.Unmarshal(bytes, c)

	if err != nil {
		return fmt.Errorf("Error deserializing credential. %s", err.Error())
	}

	return nil
}
