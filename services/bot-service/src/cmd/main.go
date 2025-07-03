package main

import (
	"bot-service/config"
	"bot-service/internal/bot"
	"bot-service/internal/handler/nats"
	"bot-service/internal/migration"
	"bot-service/internal/repository/http/user"
	"bot-service/internal/repository/http/vpn"
	notify_user "bot-service/internal/repository/pgsql/notify-user"
	"bot-service/internal/singleton"
	usrService "bot-service/internal/user"
	"github.com/novikoff-vvs/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	usrClient "pkg/infrastructure/client/user"
	vpn2 "pkg/infrastructure/client/vpn"
	singleton2 "pkg/singleton"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}
	singleton.Boot(cfg)

	lg, err := logger.NewZapLogger(cfg.Logger.Path, cfg.Logger.Name, cfg.Logger.IsOutput)
	if err != nil {
		panic(err)
	}
	lg.Info("Bot-service started", zap.Field{
		Key:    "service",
		Type:   zapcore.StringType,
		String: "Bot-service",
	})

	db, err := migration.InitDBConnection(cfg.Database)
	if err != nil {
		lg.Error(err.Error())
		return
	}

	notifyUserRepo := notify_user.NewNotifyUserRepository(db)

	singleton2.NatsPublisherBoot(cfg.Nats)

	userClient := usrClient.NewUserClient(cfg.UserService)
	userRepo := user.NewHTTPUserRepository(userClient)
	vpnClient := vpn2.NewVpnClient(cfg.VpnService, lg)
	vpnRepo := vpn.NewHTTPVPNUserRepository(vpnClient)

	userService := usrService.NewUserService(vpnRepo, userRepo)

	service := bot.NewService(cfg.BotSettings.Token, userService, vpnRepo, notifyUserRepo)
	var errChan = make(chan error, 1)
	go func(errChan chan error) {
		err := service.Run()
		if err != nil {
			lg.Error(err.Error(), zap.Field{
				Key:    "service",
				Type:   zapcore.StringType,
				String: "Bot-service",
			})
			errChan <- err
		}
	}(errChan)

	subscriptionService := nats.NewSubscriptionHandler(service)
	_, err = subscriptionService.SubscribeToEvents()
	if err != nil {
		lg.Error(err.Error(), zap.Field{
			Key:    "service",
			Type:   zapcore.StringType,
			String: "Bot-service",
		})
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	lg.Info("Bot-service started", zap.Field{
		Key:    "service",
		Type:   zapcore.StringType,
		String: "Bot-service",
	})

	select {
	case <-ch:
		break
	case <-errChan:
		{
			return
		}
	}

	lg.Info("Bot-service down", zap.Field{
		Key:    "service",
		Type:   zapcore.StringType,
		String: "Bot-service",
	})
}
