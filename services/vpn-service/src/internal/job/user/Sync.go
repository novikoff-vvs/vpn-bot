package user

import (
	"errors"
	"fmt"
	"github.com/novikoff-vvs/logger"
	"pkg/exceptions"
	pkg_user "pkg/infrastructure/client/user"
	"vpn-service/internal/service/vpn"
)

type SyncJob struct {
	client     *pkg_user.Client
	vpnService vpn.ServiceInterface
	lg         logger.Interface
}

func (j SyncJob) Run() {
	users, err := j.vpnService.GetAllUsers()
	if err != nil {
		j.lg.Error(fmt.Sprintf("get all users error: %s", err.Error()))
	}
	var uuids []string
	for _, user := range users {
		uuids = append(uuids, user.UUID)
		var req = pkg_user.GetUserByChatIdRequest{
			ChatId: user.ChatId,
		}
		_, err := j.client.GetByChatID(req)
		if errors.Is(err, exceptions.ErrModelNotFound) {
			j.lg.Info(fmt.Sprintf("Юзер удален, id: %s", user.UUID))
			var req = pkg_user.CreateUserRequest{
				ChatId: user.ChatId,
				Email:  user.Email,
				UUID:   user.UUID,
			}
			_, err := j.client.Create(req)
			if err != nil {
				j.lg.Error(fmt.Sprintf("Юзер удален, id: %s", user.UUID))
			}
			continue
		}
		if err != nil {
			j.lg.Error(fmt.Sprintf("get all users error: %s", err.Error()))
			continue
		}
	}
	syncedUsers, err := j.client.SyncUsers(pkg_user.SyncUsersRequest{UUIDs: uuids})
	if err != nil {
		j.lg.Error(fmt.Sprintf("SyncUsers error: %s", err.Error()))
		return
	}
	j.lg.Info(fmt.Sprintf("Synced users: %v", syncedUsers))
}

func NewSyncJob(client *pkg_user.Client, vpnService vpn.ServiceInterface, lg logger.Interface) *SyncJob {
	return &SyncJob{
		client:     client,
		vpnService: vpnService,
		lg:         lg,
	}
}
