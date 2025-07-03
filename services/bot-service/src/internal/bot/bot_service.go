package bot

import (
	"bot-service/internal/bot/handlers"
	router2 "bot-service/internal/bot/router"
	"bot-service/internal/repository/http/vpn"
	notify_user "bot-service/internal/repository/pgsql/notify-user"
	usrService "bot-service/internal/user"
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"log"
	"pkg/events"
)

type Service struct {
	bot            *tg.Client
	userService    usrService.ServiceInterface
	vpnRepo        vpn.RepositoryInterface
	notifyUserRepo *notify_user.NotifyUserRepository
}

func NewService(token string, userService usrService.ServiceInterface, vpnRepo vpn.RepositoryInterface, notifyUserRepo *notify_user.NotifyUserRepository) *Service {
	return &Service{
		bot:            tg.New(token),
		userService:    userService,
		vpnRepo:        vpnRepo,
		notifyUserRepo: notifyUserRepo,
	}
}

func (s *Service) Run() error {
	client := s.bot
	router := tgb.NewRouter()

	r := router2.NewRouter(router)
	commandH := handlers.NewCommandHandler(s.userService, s.notifyUserRepo)
	reactionH := handlers.NewReactionHandler(s.bot)
	callbackQueryH := handlers.NewCallbackQueryHandler(s.userService, s.vpnRepo, s.bot)
	messageH := handlers.NewMessageHandler(s.userService, s.vpnRepo)

	r.RegisterCommandHandlers(commandH)
	r.RegisterReactionHandlers(reactionH)
	r.RegisterCallbackQueryHandlers(callbackQueryH)
	r.RegisterMessageHandlers(messageH)

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

func (s *Service) NotifySubscriptionRefreshed(event events.SubscriptionRefreshed) (err error) {
	m := tg.HTML.Text(
		tg.HTML.Bold("💸 Оплата успешно получена!\n"),
		tg.HTML.Bold("🎉 Ваша подписка обновлена — спасибо, что остаетесь с нами!\n"),
		"Все функции доступны без ограничений. 💫\n\n",
		"Если появятся вопросы — мы всегда рядом. 🤝",
	)
	err = s.bot.SendMessage(tg.UserID(event.ChatId), m).ParseMode(tg.HTML).DoVoid(context.Background())

	if err != nil {
		//todo добавить логирование
		log.Println(err)
		return err
	}
	return nil
}

func (s *Service) NotifyDeactivatedUser(event events.UserDeactivated) error {
	m := tg.HTML.Text(
		tg.HTML.Bold("😵 Вы деактивированы!"),
		tg.HTML.Blockquote("Кажется вы не оплатили подписку или нарушили наши правила! 😔"),
		"Если появятся вопросы — то обращайтесь в теххподдержку",
	)
	err := s.bot.SendMessage(tg.UserID(event.ChatId), m).ParseMode(tg.HTML).DoVoid(context.Background())
	if err != nil {
		//todo добавить логирование
		log.Println(err)
		return err
	}
	return nil
}

func (s *Service) NotifyNewMessage(event events.NewMessage) error {
	m := tg.HTML.Text(event.Message)
	users, err := s.notifyUserRepo.All()
	if err != nil {
		return err
	}

	for _, user := range users {
		err = s.bot.SendMessage(tg.UserID(user.ChatId), m).ParseMode(tg.HTML).DoVoid(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(fmt.Sprintf("Отправлено %d", user.ChatId))
	}

	return nil
}
