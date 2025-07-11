package singleton

import (
	"github.com/novikoff-vvs/logger"
	"pkg/config"
)

// Boot TODO нужно подумать как сделать единй конфиг для юзер сервиса, хм
func BootAll(userCfg config.UserService, vpnCfg config.VpnService, lg logger.Interface, publisherCfg config.NatsPublisher) {
	UserClientBoot(userCfg)
	VpnClientBoot(vpnCfg, lg)
	SubscriptionClientBoot(userCfg)
	NatsPublisherBoot(publisherCfg)
}
