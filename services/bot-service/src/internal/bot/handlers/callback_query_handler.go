package handlers

import (
	"bot-service/internal/repository/http/vpn"
	"bot-service/internal/singleton"
	"bot-service/internal/user"
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

type CallbackQueryHandlerInterface interface {
	GetCallbackQueryHandlersFunc() map[string]tgb.CallbackQueryHandler
	GetFilter() tgb.Filter
}

type CallbackQueryHandler struct {
	userService user.ServiceInterface
	vpnRepo     vpn.RepositoryInterface
}

func NewCallbackQueryHandler(userService user.ServiceInterface, vpnRepo vpn.RepositoryInterface) *CallbackQueryHandler {
	return &CallbackQueryHandler{
		userService: userService,
		vpnRepo:     vpnRepo,
	}
}

func (h CallbackQueryHandler) GetFilter() tgb.Filter {
	return nil
}

func (h CallbackQueryHandler) GetCallbackQueryHandlersFunc() map[string]tgb.CallbackQueryHandler {
	return map[string]tgb.CallbackQueryHandler{
		"get_link": h.GetVessaLink,
	}
}

func (h CallbackQueryHandler) GetConfigHandle(context.Context, *tgb.CallbackQueryUpdate) error {
	panic("implement me")

}

func (h CallbackQueryHandler) GetVessaLink(ctx context.Context, update *tgb.CallbackQueryUpdate) error {
	link, err := h.vpnRepo.GetSubscriptionLinkByChatId(int64(update.CallbackQuery.From.ID))
	if err != nil {
		return err
	}
	err = update.Update.Reply(ctx, tg.NewEditMessageTextCall(update.CallbackQuery.From.ID, update.Message.MessageID(), tg.HTML.Text(
		tg.HTML.Text(tg.HTML.Bold("🔐 Важно:"), tg.HTML.Blockquote("Эта ссылка является вашим личным доступом.\n\nНикому не передавайте её – это может привести к потере аккаунта.")),
		"",
		tg.HTML.Line(tg.HTML.Bold("🔗 Ваша подписка: "), tg.HTML.Link("SUBSCRIPTION-URL", link)),
	)).
		ParseMode(tg.HTML))

	return update.
		Client.
		SendMessage(update.CallbackQuery.From.ID, "Чем еще могу помочь?").
		ParseMode(tg.HTML).
		ReplyMarkup(singleton.MessageBuilder().GetMainMenuKeyboad()).
		DoVoid(ctx)
}
