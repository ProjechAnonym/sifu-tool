package cert

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/go-acme/lego/v4/certcrypto"
)

func Encrypt(plaintext []byte, key string) (string, error) {
	hashKey := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(hashKey[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
func Decrypt(ciphertext string, key string) (*ecdsa.PrivateKey, error) {
	hashKey := sha256.Sum256([]byte(key))
	
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(hashKey[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParseECPrivateKey(plaintext)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func certCrypto(method string) string {
	switch method {
		case "ec256":
			return string(certcrypto.EC256)
		case "ec384":
			return string(certcrypto.EC384)
		case "rsa2048":
			return string(certcrypto.RSA2048)
	}
	return string(certcrypto.RSA2048)
}