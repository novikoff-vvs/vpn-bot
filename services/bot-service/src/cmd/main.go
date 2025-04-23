package main

import (
	"bot-service/config"
	"bot-service/internal/bot"
	"bot-service/internal/repository/http/user"
	"bot-service/internal/repository/http/vpn"
	"bot-service/internal/singleton"
	usrService "bot-service/internal/user"
	"github.com/novikoff-vvs/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	usrClient "pkg/infrastructure/client/user"
	vpn2 "pkg/infrastructure/client/vpn"
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

	userClient := usrClient.NewUserClient(cfg.UserService)
	userRepo := user.NewHTTPUserRepository(userClient)
	vpnClient := vpn2.NewVpnClient(cfg.VpnService, lg)
	vpnRepo := vpn.NewHTTPVPNUserRepository(vpnClient)

	userService := usrService.NewUserService(vpnRepo, userRepo)

	service := bot.NewService(cfg.BotSettings.Token, userService, vpnRepo)
	go func() {
		err := service.Run()
		if err != nil {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	lg.Info("Bot-service started", zap.Field{
		Key:    "service",
		Type:   zapcore.StringType,
		String: "Bot-service",
	})
	<-ch
	lg.Info("Bot-service down", zap.Field{
		Key:    "service",
		Type:   zapcore.StringType,
		String: "Bot-service",
	})
}
