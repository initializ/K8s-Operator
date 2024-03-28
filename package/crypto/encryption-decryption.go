package crypto

import (
	"encoding/base64"
	"encoding/hex"
)

// DecryptPrivateKey decrypts the private key using the derived key, IV, and authTag.
// DecryptPrivateKey decrypts the private key using the derived key, IV, and authTag.
func DecryptPrivateKey(orgID string, salt string, encryptedData string, iv string, authTag string) (string, string, error) {
	// Decode hex strings
	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		return "", "", err
	}

	encryptedDataBytes, err := hex.DecodeString(encryptedData)
	if err != nil {
		return "", "", err
	}

	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return "", "", err
	}

	authTagBytes, err := hex.DecodeString(authTag)
	if err != nil {
		return "", "", err
	}

	// Derive key from organization
	orgKey, err := DeriveKeyFromOrganisation(orgID, saltBytes, 32, 169696)
	if err != nil {
		return "", "", err
	}

	// Decrypt the private key using the derived key, iv, and authTag
	decryptedData, err := DecryptAESGCM(encryptedDataBytes, orgKey, ivBytes, authTagBytes)
	if err != nil {
		return "", "", err
	}

	privateKeyBase64 := base64.StdEncoding.EncodeToString(decryptedData)
	orgKeyHex := hex.EncodeToString(orgKey)

	return privateKeyBase64, orgKeyHex, nil
}

// DecryptKeyValuePair decrypts a key-value pair using AES 256 GCM.
func DecryptKeyValuePair(privateKey, keyIV, keyTag, keyCiphertext, valueIV, valueTag, valueCiphertext string) (string, string, error) {
	// Decode base64-encoded private key
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", "", err
	}

	// Convert hex strings to byte slices
	keyIVBytes, err := hex.DecodeString(keyIV)
	if err != nil {
		return "", "", err
	}
	valueIVBytes, err := hex.DecodeString(valueIV)
	if err != nil {
		return "", "", err
	}
	keyTagBytes, err := hex.DecodeString(keyTag)
	if err != nil {
		return "", "", err
	}
	valueTagBytes, err := hex.DecodeString(valueTag)
	if err != nil {
		return "", "", err
	}
	keyCiphertextBytes, err := hex.DecodeString(keyCiphertext)
	if err != nil {
		return "", "", err
	}
	valueCiphertextBytes, err := hex.DecodeString(valueCiphertext)
	if err != nil {
		return "", "", err
	}

	// Decrypt key
	decryptedKey, err := DecryptAESGCM(keyCiphertextBytes, decodedPrivateKey, keyIVBytes, keyTagBytes)
	if err != nil {
		return "", "", err
	}

	// Decrypt value
	decryptedValue, err := DecryptAESGCM(valueCiphertextBytes, decodedPrivateKey, valueIVBytes, valueTagBytes)
	if err != nil {
		return "", "", err
	}

	// Return decrypted key and value as strings
	return string(decryptedKey), string(decryptedValue), nil
}
