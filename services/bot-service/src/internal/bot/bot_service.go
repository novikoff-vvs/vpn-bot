package bot

import (
	"bot-service/internal/bot/handlers"
	router2 "bot-service/internal/bot/router"
	"bot-service/internal/repository/http/vpn"
	usrService "bot-service/internal/user"
	"context"
	"encoding/json"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"github.com/nats-io/nats.go"
	"log"
	"pkg/events"
	singleton2 "pkg/singleton"
)

type Service struct {
	bot         *tg.Client
	userService usrService.ServiceInterface
	vpnRepo     vpn.RepositoryInterface
}

func NewService(token string, userService usrService.ServiceInterface, vpnRepo vpn.RepositoryInterface) *Service {
	return &Service{
		bot:         tg.New(token),
		userService: userService,
		vpnRepo:     vpnRepo,
	}
}

func (s *Service) Run() error {
	client := s.bot
	router := tgb.NewRouter()

	r := router2.NewRouter(router)
	commandH := handlers.NewCommandHandler(s.userService)
	reactionH := handlers.NewReactionHandler(s.bot)
	callbackQueryH := handlers.NewCallbackQueryHandler(s.userService, s.vpnRepo, s.bot)
	messageH := handlers.NewMessageHandler(s.userService, s.vpnRepo)

	r.RegisterCommandHandlers(commandH)
	r.RegisterReactionHandlers(reactionH)
	r.RegisterCallbackQueryHandlers(callbackQueryH)
	r.RegisterMessageHandlers(messageH)

	_, err := singleton2.NatsPublisher().Subscribe("events.subscription.refreshed", "bot_service_subscription_refreshed_consumer", func(msg *nats.Msg) {
		var event events.SubscriptionRefreshed
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
		}
		log.Println("Сообщение получено!")

		m := tg.HTML.Text(
			tg.HTML.Bold("💸 Оплата успешно получена!\n"),
			tg.HTML.Bold("🎉 Ваша подписка обновлена — спасибо, что остаетесь с нами!\n"),
			"Теперь все функции доступны без ограничений. 💫\n\n",
			"Если появятся вопросы — мы всегда рядом. 🤝",
		)
		err = client.SendMessage(tg.UserID(event.ChatId), m).ParseMode(tg.HTML).DoVoid(context.Background())

		if err != nil {
			log.Println(err)
		}
		err = msg.Ack()
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Subscribed")
	return tgb.NewPoller(
		router,
		client,
		tgb.WithPollerAllowedUpdates(
			tg.UpdateTypeMessage,
			tg.UpdateTypeMessageReaction,
			tg.UpdateTypeCallbackQuery,
		),
	).Run(context.Background())
}
