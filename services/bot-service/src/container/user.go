package container

import (
	"pkg/models"
	"time"
)

type CachedUser struct {
	User   models.VpnUser
	Expiry time.Time
}

type UserContainer struct {
	users map[int64]CachedUser
}

func (uc UserContainer) Register(user models.VpnUser) {
	uc.users[user.ChatId] = CachedUser{
		User:   user,
		Expiry: time.Now().Add(time.Hour),
	}
}

func (uc UserContainer) Put(user models.VpnUser) {
	uc.users[user.ChatId] = CachedUser{
		User:   user,
		Expiry: time.Now().Add(time.Hour),
	}
}

func (uc UserContainer) Get(chatId int64) (CachedUser, bool) {
	cachedUser, ok := uc.users[chatId]
	if !ok {
		return CachedUser{}, false
	}

	if time.Now().After(cachedUser.Expiry) {
		delete(uc.users, chatId)
		return CachedUser{}, false
	}

	return cachedUser, false
}

func NewUserContainer() *UserContainer {
	return &UserContainer{
		users: make(map[int64]CachedUser),
	}
}
