package job

import "pkg/service/vpn"

type SyncUsersJob struct {
	vpnService vpn.ServiceInterface
}

func NewSyncUsersJob(vpnService vpn.ServiceInterface) *SyncUsersJob {
	return &SyncUsersJob{vpnService: vpnService}
}

func (j SyncUsersJob) Run() {
}
