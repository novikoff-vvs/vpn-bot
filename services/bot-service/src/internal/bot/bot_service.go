package bot

import (
	"bot-service/internal/bot/handlers"
	router2 "bot-service/internal/bot/router"
	"bot-service/internal/repository/http/vpn"
	usrService "bot-service/internal/user"
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
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
	callbackQueryH := handlers.NewCallbackQueryHandler(s.userService, s.vpnRepo)
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
