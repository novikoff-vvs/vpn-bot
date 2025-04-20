package handlers

import (
	"bot-service/internal/vpn"
	"context"
	"github.com/mr-linch/go-tg/tgb"
)

type CallbackQueryHandlerInterface interface {
	GetCallbackQueryHandlersFunc() map[string]tgb.CallbackQueryHandler
	GetFilter() tgb.Filter
}

type CallbackQueryHandler struct {
	vpnService vpn.ServiceInterface
}

func NewCallbackQueryHandler(vpnService vpn.ServiceInterface) *CallbackQueryHandler {
	return &CallbackQueryHandler{
		vpnService: vpnService,
	}
}

func (h CallbackQueryHandler) GetFilter() tgb.Filter {
	return nil
}

func (h CallbackQueryHandler) GetCallbackQueryHandlersFunc() map[string]tgb.CallbackQueryHandler {
	return map[string]tgb.CallbackQueryHandler{
		"register": h.GetVessaLink,
	}
}

func (h CallbackQueryHandler) GetConfigHandle(context.Context, *tgb.CallbackQueryUpdate) error {
	panic("implement me")
	return nil
}

func (h CallbackQueryHandler) GetVessaLink(ctx context.Context, update *tgb.CallbackQueryUpdate) error {
	panic(update.Update.Message)
	_, err := h.vpnService.UserGetByChatId(int64(update.Message.Chat().ID))
	if err != nil {
		return err
	}

	return update.AnswerText("qwrqwr", true).DoVoid(ctx)
}
