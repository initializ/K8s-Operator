package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"github.com/xdg-go/pbkdf2"
)

// DecryptAESGCM decrypts data using AES-256-GCM symmetric encryption.
func DecryptAESGCM(encryptedData []byte, key []byte, iv []byte, authTag []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCMWithNonceSize(block, len(iv))
	if err != nil {
		return nil, err
	}

	var nonce = iv
	var ciphertext = append(encryptedData, authTag...)

	decrypted, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// DeriveKeyFromOrganisation derives a cryptographic key from an organization using PBKDF2.
func DeriveKeyFromOrganisation(organizationID string, salt []byte, keyLength int, iterations int) ([]byte, error) {
	derivedKey := pbkdf2.Key([]byte(organizationID), salt, iterations, keyLength, sha512.New)
	return derivedKey, nil
}
