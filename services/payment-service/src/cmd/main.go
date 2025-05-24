package main

import (
	"github.com/novikoff-vvs/logger"
	"log"
	"payment-service/config"
	"payment-service/internal/controller/yoomoney"
	"pkg/infrastructure/client/user"
	"pkg/infrastructure/http"
)

func main() {
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}
	lg, err := logger.NewZapLogger(cfg.Logger.Path, cfg.Logger.Name, cfg.Logger.IsOutput)
	if err != nil {
		log.Println(err.Error())
		return
	}

	client := user.NewUserClient(cfg.UserService)
	server := http.NewServer(lg)
	server.RegisterStatic()
	yoomoney.RegisterRoutes(server, client)
	err = server.Run(cfg.Base.AppPort)
	if err != nil {
		lg.Error(err.Error())
		return
	}
}
