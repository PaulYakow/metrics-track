package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type Cryptographer struct {
	publicKey *rsa.PublicKey
}

func NewCryptographer(publicKeyPath string) (*Cryptographer, error) {
	bytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := convertBytesToPublicKey(bytes)
	if err != nil {
		return nil, err
	}

	return &Cryptographer{publicKey: publicKey}, err
}

func (c *Cryptographer) Encrypt(data []byte) ([]byte, error) {
	msgLen := len(data)
	step := c.publicKey.Size() - 2*sha512.Size - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, c.publicKey, data[start:finish], nil)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func convertBytesToPublicKey(keyBytes []byte) (*rsa.PublicKey, error) {
	var err error

	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes

	publicKey, err := x509.ParsePKCS1PublicKey(blockBytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
