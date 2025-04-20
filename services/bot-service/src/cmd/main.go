package main

import (
	"bot-service/config"
	"bot-service/internal/bot"
	"bot-service/internal/controllers"
	usrClient "bot-service/internal/infrastructure/client/user"
	"bot-service/internal/infrastructure/http"
	"bot-service/internal/repository/http/user"
	usrService "bot-service/internal/user"
	"bot-service/internal/vpn"
	"github.com/novikoff-vvs/logger"
)

func main() {
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}
	client := usrClient.NewUserClient(cfg.UserService)
	userRepo := user.NewHTTPUserRepository(client)

	lg, err := logger.NewZapLogger(cfg.Logger.Path, cfg.Logger.Name, cfg.Logger.IsOutput)
	if err != nil {
		panic(err)
	}
	vpnService, err := vpn.NewVPNService(cfg.VpnService, lg)
	if err != nil {
		panic(err)
	}

	userService := usrService.NewUserService(vpnService, userRepo)

	service := bot.NewService(cfg.BotSettings.Token, userService)
	go func() {
		err := service.Run()
		if err != nil {
			panic(err)
		}
	}()

	api := http.NewApiServer()

	api.SetupControllers([]http.Controller{
		controllers.NewWebhookController(vpnService),
	})

	err = api.Start()
	if err != nil {
		panic(err)
	}
}
