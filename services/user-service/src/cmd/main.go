package main

import (
	"fmt"
	"github.com/novikoff-vvs/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"pkg/infrastructure/DB/gorm"
	"pkg/infrastructure/http"
	"pkg/singleton"
	"user-service/config"
	"user-service/docs"
	"user-service/internal/controller/subscription"
	"user-service/internal/controller/user"
	"user-service/internal/migration"
	"user-service/internal/plan"
	"user-service/internal/repository/plan/sqlite"
	sqliteSubscription "user-service/internal/repository/subscription/sqlite"
	sqliteUser "user-service/internal/repository/user/sqlite"
	subscription2 "user-service/internal/subscription"
	user2 "user-service/internal/user"
)

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name	Apache 2.0
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
	singleton.NatsPublisherBoot(cfg.NatsPublisher)

	db, err := migration.InitDBConnection(cfg.Database)
	if err != nil {
		LoggingService.Error(err.Error())
		return
	}

	newDatabaseService := gorm.NewDBService(db)

	userRepo := sqliteUser.NewUserRepository(newDatabaseService)
	subscrRepo := sqliteSubscription.NewSubscriptionRepository(newDatabaseService)

	userService := user2.NewUserService(userRepo, subscrRepo)
	planService := plan.NewPlanService(sqlite.NewPlanRepository(newDatabaseService))
	subscriptionRepo := sqliteSubscription.NewSubscriptionRepository(newDatabaseService)
	subscriptionService := subscription2.NewSubscriptionService(subscriptionRepo, planService)

	s := http.NewServer(LoggingService)

	if cfg.Base.Swagger {
		docs.SwaggerInfo.Title = "User Service"
		docs.SwaggerInfo.Description = "Сервис пользователей"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Base.AppPort)
		docs.SwaggerInfo.BasePath = "/api"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}

		s.GetWebGroup().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	user.RegisterRoutes(s, userService, LoggingService)
	subscription.RegisterRoutes(s, subscriptionService)
	err = s.Run(cfg.Base.AppPort)
	if err != nil {
		LoggingService.Error(err.Error())
		return
	}
}
