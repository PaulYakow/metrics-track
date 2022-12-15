package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type Decryptor struct {
	privateKey *rsa.PrivateKey
}

func NewDecryptor(privateKeyPath string) (*Decryptor, error) {
	bytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := convertBytesToPrivateKey(bytes)
	if err != nil {
		return nil, err
	}

	return &Decryptor{privateKey: privateKey}, err
}

func (d *Decryptor) Decrypt(data []byte) ([]byte, error) {
	msgLen := len(data)
	step := d.privateKey.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, d.privateKey, data[start:finish], nil)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

func convertBytesToPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	var err error

	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes

	privateKey, err := x509.ParsePKCS1PrivateKey(blockBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
