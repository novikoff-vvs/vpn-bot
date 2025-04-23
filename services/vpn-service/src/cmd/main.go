package main

import (
	"fmt"
	"github.com/novikoff-vvs/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"pkg/infrastructure/http"
	"pkg/singleton"
	"vpn-service/config"
	"vpn-service/docs"
	"vpn-service/internal/controller/vpn"
	"vpn-service/internal/cron"
	vpn2 "vpn-service/internal/service/vpn"
)

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// LoggingService @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
var LoggingService logger.Interface

func main() {
	cfg, err := config.LoadConfigs()
	if err != nil {
		fmt.Println("Error loading configs: ", err)
		return
	}

	LoggingService, err = logger.NewZapLogger(cfg.Logger.Path, cfg.Logger.Name, cfg.Logger.IsOutput)
	if err != nil {
		fmt.Println("Error initializing logger: ", err)
		return
	}
	LoggingService.Info("Initializing app")

	vpnService, err := vpn2.NewVPNService(cfg.Xui, LoggingService)
	if err != nil {
		return
	}
	singleton.UserClientBoot(cfg.UserService)

	worker, err := cron.NewWorker(vpnService, singleton.UserClient(), LoggingService)
	if err != nil {
		panic(err)
	}
	worker.Start()
	s := http.NewServer(LoggingService)

	if cfg.Base.Swagger {
		docs.SwaggerInfo.Title = "Vpn Service"
		docs.SwaggerInfo.Description = "Сервис взаимодействия с 3x-ui"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Base.AppPort)
		docs.SwaggerInfo.BasePath = "/api"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}

		s.GetWebGroup().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	vpn.RegisterRoutes(s, vpnService)
	err = s.Run(cfg.Base.AppPort)
	if err != nil {
		LoggingService.Error(err.Error())
		return
	}
}
