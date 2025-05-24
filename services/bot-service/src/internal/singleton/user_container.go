package singleton

import (
	"bot-service/container"
	"sync"
)

var userContainer *container.UserContainer
var userContainerOnce sync.Once

func userContainerBoot() {
	userContainerOnce.Do(func() {
		userContainer = container.NewUserContainer()
	})
}

func UserContainer() *container.UserContainer {
	return userContainer
}
