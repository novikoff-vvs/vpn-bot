package main

import (
	"github.com/novikoff-vvs/logger"
	"log"
	"payment-service/config"
	"payment-service/internal/controller/yoomoney"
	"payment-service/internal/migration"
	"payment-service/internal/payment"
	paymentRepo "payment-service/internal/repository/payment"
	singleton2 "payment-service/internal/singleton"
	"pkg/infrastructure/DB/gorm"
	"pkg/infrastructure/http"
	"pkg/singleton"
)

func main() {
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}
	lg, err := logger.NewZapLogger(cfg.Logger.Path, cfg.Logger.Name, true) //
	if err != nil {
		log.Println(err.Error())
		return
	}
	db, err := migration.InitDBConnection(cfg.Database)
	if err != nil {
		lg.Error(err.Error())
		panic(err)
		return
	}
	dbService := gorm.NewDBService(db)

	bootSingletons(cfg)

	paymentRepository := paymentRepo.NewPaymentRepository(dbService)
	paymentServices := payment.NewPaymentService(paymentRepository)
	client := singleton.UserClient()

	server := http.NewServer(lg)
	server.RegisterStatic()
	yoomoney.RegisterRoutes(server, client, paymentServices)

	err = server.Run(cfg.Base.AppPort)

	if err != nil {
		lg.Error(err.Error())
		return
	}
}

func bootSingletons(cfg *config.Config) {
	singleton.UserClientBoot(cfg.UserService)
	singleton2.CryptoServiceBoot(cfg.Crypto)
	singleton.SubscriptionClientBoot(cfg.UserService)
}
