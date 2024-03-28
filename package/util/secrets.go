package util

import (
	"Initializ-Operator/package/api"
	"Initializ-Operator/package/crypto"
	"github.com/go-resty/resty/v2"
)

type KeyValue struct {
	Key   string
	Value string
}

func GetPlainTextSecretsViaServiceToken(serviceToken string, org_id string, envSlug string, objectIds []string) ([]KeyValue, error) {
	var decryptedPairs []KeyValue

	httpClient := resty.New()
	request := api.GetEncryptedWorkspaceKeyRequest{
		OrgID: org_id,
	}

	response, err := api.CallGetEncryptedWorkspaceKey(httpClient, serviceToken, request)
	if err != nil {
		return nil, err
	}

	salt := response.Data[0].EncryptedPrivateSalt
	encryptedData := response.Data[0].EncryptedPrivateKey
	iv := response.Data[0].EncryptedPrivateIv
	authTag := response.Data[0].EncryptedPrivateAuthTag
	orgID := response.Data[0].OrgID

	decryptedData, _, err := crypto.DecryptPrivateKey(orgID, salt, encryptedData, iv, authTag)
	if err != nil {
		return nil, err
	}
	privateKey := decryptedData

	secretsRequest := api.GetEncryptedSecretsRequest{
		SecretIDs:   objectIds,
		Environment: envSlug,
		OrgId:       org_id,
	}

	secretResponse, err := api.CallGetEncryptedSecrets(secretsRequest, serviceToken)
	if err != nil {
		return nil, err
	}

	// Decrypt each secret key-value pair using the privateKey
	for _, secret := range secretResponse.Data {
		decryptedKey, decryptedValue, err := crypto.DecryptKeyValuePair(
			privateKey,
			secret.SecretKeyIV,
			secret.SecretKeyTag,
			secret.SecretKeyCiphertext,
			secret.SecretValueIV,
			secret.SecretValueTag,
			secret.SecretValueCiphertext,
		)
		if err != nil {
			return nil, err
		}

		decryptedPairs = append(decryptedPairs, KeyValue{
			Key:   decryptedKey,
			Value: decryptedValue,
		})
	}

	return decryptedPairs, nil
}
