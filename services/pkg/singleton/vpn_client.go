package singleton

import (
	"github.com/novikoff-vvs/logger"
	"pkg/config"
	"pkg/infrastructure/client/vpn"
	"sync"
)

var vpnClient *vpn.Client
var onceVpn sync.Once

func VpnClientBoot(cfg config.VpnService, lg logger.Interface) {
	onceVpn.Do(func() {
		vpnClient = vpn.NewVpnClient(cfg, lg)
	})
}

func VpnClient() *vpn.Client {
	return vpnClient
}
