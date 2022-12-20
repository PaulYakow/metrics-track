package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// Cryptographer хранит публичный ключ для шифрования данных и функцию чтения данных (для получения ключа).
type Cryptographer struct {
	publicKey  *rsa.PublicKey
	dataReader func() ([]byte, error)
}

// NewCryptographer - создаёт объект Cryptographer с заданным путём к публичному ключу.
func NewCryptographer(publicKeyPath string) (*Cryptographer, error) {
	c := &Cryptographer{dataReader: func() ([]byte, error) {
		bytes, err := os.ReadFile(publicKeyPath)
		return bytes, err
	}}

	keyBytes, err := c.dataReader()
	publicKey, err := convertBytesToPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}

	c.publicKey = publicKey
	return c, nil
}

// Encrypt - шифрует данные на основе публичного ключа.
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
