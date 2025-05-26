package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

type Service struct {
	CryptoKey []byte
}

func (s Service) Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(s.CryptoKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	plainText = plainText + ":" + time.Now().Format("2006-01-02_15:04:05")

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (s Service) Decrypt(encText string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(s.CryptoKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	return aesGCM.Open(nil, nonce, cipherText, nil)
}

func NewCryptoService(key []byte) *Service {
	return &Service{
		CryptoKey: key,
	}
}
