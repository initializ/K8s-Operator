// models/secret_model.go
package model

import (
	"time"
)

// Secret represents the structure of the secret data
type Secret struct {
	ID                    string    `json:"id" bson:"_id,omitempty"`
	WorkspaceID           string    `json:"workspace_id" bson:"workspace_id"`
	Environment           string    `json:"environment" bson:"environment"`
	SecretKeyCiphertext   string    `json:"secretKeyCiphertext" bson:"secretKeyCiphertext"`
	SecretKeyIV           string    `json:"secretKeyIV" bson:"secretKeyIV"`
	SecretKeyTag          string    `json:"secretKeyTag" bson:"secretKeyTag"`
	SecretValueCiphertext string    `json:"secretValueCiphertext" bson:"secretValueCiphertext"`
	SecretValueIV         string    `json:"secretValueIV" bson:"secretValueIV"`
	SecretValueTag        string    `json:"secretValueTag" bson:"secretValueTag"`
	CreatedAt             time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt" bson:"updatedAt"`
}
