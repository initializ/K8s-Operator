package api

import "time"

type KeyData struct {
	OrgID                   string `json:"org_id"`
	PublicKey               string `json:"publicKey"`
	EncryptedPrivateKey     string `json:"encryptedPrivateKey"`
	EncryptedPrivateIv      string `json:"encryptedPrivateIv"`
	EncryptedPrivateSalt    string `json:"encryptedPrivateSalt"`
	EncryptedPrivateAuthTag string `json:"encryptedPrivateAuthTag"`
}

type GetEncryptedWorkspaceKeyResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []KeyData `json:"data"`
}

type GetEncryptedWorkspaceKeyRequest struct {
	OrgID string `json:"org_id"`
}

type GetEncryptedSecretsRequest struct {
	SecretIDs   []string `json:"secret_ids"`
	Environment string   `json:"environment"`
	OrgId       string   `json:"org_id"`
}

type SecretData struct {
	ID                    string    `json:"id"`
	WorkspaceID           string    `json:"workspace_id"`
	Environment           string    `json:"environment"`
	SecretKeyCiphertext   string    `json:"secretKeyCiphertext"`
	SecretKeyIV           string    `json:"secretKeyIV"`
	SecretKeyTag          string    `json:"secretKeyTag"`
	SecretValueCiphertext string    `json:"secretValueCiphertext"`
	SecretValueIV         string    `json:"secretValueIV"`
	SecretValueTag        string    `json:"secretValueTag"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

type GetEncryptedSecretResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    []SecretData `json:"data"`
}

