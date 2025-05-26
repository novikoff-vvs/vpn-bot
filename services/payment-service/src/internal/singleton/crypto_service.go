package singleton

import (
	"payment-service/config"
	"payment-service/internal/crypto"
	"sync"
)

var cryptoService *crypto.Service
var onceCryptoService sync.Once

func CryptoServiceBoot(cfg config.Crypto) {
	onceCryptoService.Do(func() {
		cryptoService = crypto.NewCryptoService([]byte(cfg.Key))
	})
}

func CryptoService() *crypto.Service {
	return cryptoService
}
