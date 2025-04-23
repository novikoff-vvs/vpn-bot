package singleton

import (
	"pkg/config"
	"pkg/infrastructure/client/user"
	"sync"
)

var userClient *user.Client
var onceUser sync.Once

func UserClientBoot(cfg config.UserService) {
	onceUser.Do(func() {
		userClient = user.NewUserClient(cfg)
	})
}

func UserClient() *user.Client {
	return userClient
}
