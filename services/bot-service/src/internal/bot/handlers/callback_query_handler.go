package handlers

import (
	"bot-service/internal/repository/http/vpn"
	"bot-service/internal/singleton"
	"bot-service/internal/user"
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"log"
	"pkg/models"
)

type CallbackQueryHandlerInterface interface {
	GetCallbackQueryHandlersFunc() map[string]tgb.CallbackQueryHandler
	GetFilter() tgb.Filter
}

type CallbackQueryHandler struct {
	userService user.ServiceInterface
	vpnRepo     vpn.RepositoryInterface
	bot         *tg.Client
}

func NewCallbackQueryHandler(userService user.ServiceInterface, vpnRepo vpn.RepositoryInterface, bot *tg.Client) *CallbackQueryHandler {
	return &CallbackQueryHandler{
		userService: userService,
		vpnRepo:     vpnRepo,
		bot:         bot,
	}
}

func (h CallbackQueryHandler) GetFilter() tgb.Filter {
	return nil
}

func (h CallbackQueryHandler) GetCallbackQueryHandlersFunc() map[string]tgb.CallbackQueryHandler {
	return map[string]tgb.CallbackQueryHandler{
		"get_link": h.GetVessaLink,
		"payment":  h.OpenPayment,
	}
}

func (h CallbackQueryHandler) GetConfigHandle(context.Context, *tgb.CallbackQueryUpdate) error {
	panic("implement me")

}

func (h CallbackQueryHandler) GetVessaLink(ctx context.Context, update *tgb.CallbackQueryUpdate) error {
	cachedUser, _ := singleton.UserContainer().Get(int64(update.CallbackQuery.From.ID))
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
	log.Println(err) //TODO добавить логер
	return update.
		Client.
		SendMessage(update.CallbackQuery.From.ID, "Чем еще могу помочь?").
		ParseMode(tg.HTML).
		ReplyMarkup(singleton.MessageBuilder().GetMainMenuKeyboad(cachedUser.User.UUID)).
		DoVoid(ctx)
}

func (h CallbackQueryHandler) OpenPayment(ctx context.Context, update *tgb.CallbackQueryUpdate) error {
	chatId := int64(update.CallbackQuery.From.ID)
	var usr models.VpnUser
	var err error
	cachedUser, ok := singleton.UserContainer().Get(chatId)
	if !ok {
		usr, err = h.userService.UserGetByChatId(chatId)
		if err != nil {
			return err
		}
	} else {
		usr = cachedUser.User
	}
	err = h.bot.SendMessage(update.CallbackQuery.From.ID, tg.HTML.Text(
		tg.HTML.Text(tg.HTML.Blockquote("Для продолжения оплаты, просьба, перейти в наше приложение!")),
	)).
		ReplyMarkup(singleton.MessageBuilder().GetPaymentMenuKeyboard(usr.UUID)).
		ParseMode(tg.HTML).DoVoid(ctx)
	//TODO добавить логер
	return err
}
