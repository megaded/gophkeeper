package manager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"gophkeeper/internal/config"
	"io"
)

type Crypter interface {
	Encrypt(content string) ([]byte, error)
	Decrypt(content []byte) (string, error)
}

type CryptoManager struct {
	key string
}

func NewCryptoManager(cfg config.Config) CryptoManager {
	return CryptoManager{key: cfg.SecretKey}
}

func (c *CryptoManager) Encrypt(content string) ([]byte, error) {
	byteMsg := []byte(content)
	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)
	return cipherText, nil
}

func (c *CryptoManager) Decrypt(content []byte) (string, error) {
	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		return "", err
	}

	if len(content) < aes.BlockSize {
		return "", err
	}

	iv := content[:aes.BlockSize]
	cipherText := content[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
