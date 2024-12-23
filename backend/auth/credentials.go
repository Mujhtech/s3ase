package auth

type CredentialType string

const (
	CredentialTypeBasic  CredentialType = "basic"
	CredentialTypeBearer CredentialType = "bearer"
)

type Credentials struct {
	Type     CredentialType `json:"type"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Token    string         `json:"token"`
}
