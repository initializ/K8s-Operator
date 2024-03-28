package model

type Key struct {
	OrgID                   string `json:"org_id" bson:"org_id"`
	PublicKey               string `json:"publicKey" bson:"publicKey"`
	EncryptedPrivateKey     string `json:"encryptedPrivateKey" bson:"encryptedPrivateKey"`
	EncryptedIv             string `json:"encryptedPrivateIv" bson:"encryptedPrivateIv"`
	EncryptedSalt           string `json:"encryptedPrivateSalt" bson:"encryptedPrivateSalt"`
	EncryptedPrivateAuthTag string `json:"encryptedPrivateAuthTag" bson:"encryptedPrivateAuthTag"`
}