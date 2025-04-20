package bot

import (
	"bot-service/internal/bot/handlers"
	router2 "bot-service/internal/bot/router"
	usrService "bot-service/internal/user"
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

type Service struct {
	bot         *tg.Client
	userService usrService.ServiceInterface
}

func NewService(token string, userService usrService.ServiceInterface) *Service {
	return &Service{
		bot:         tg.New(token),
		userService: userService,
	}
}

func (s *Service) Run() error {
	client := s.bot
	router := tgb.NewRouter()

	r := router2.NewRouter(router)
	commandH := handlers.NewCommandHandler(s.userService)
	reactionH := handlers.NewReactionHandler(s.bot)
	callbackQueryH := handlers.NewCallbackQueryHandler(s.userService)
	messageH := handlers.NewMessageHandler(s.userService)

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
