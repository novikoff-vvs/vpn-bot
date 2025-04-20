package handlers

import (
	"bot-service/internal/bot/message"
	"bot-service/internal/models"
	usrService "bot-service/internal/user"
	"context"
	"github.com/mr-linch/go-tg/tgb"
)

type CommandHandlerInterface interface {
	GetHandlerFuncs() map[string]tgb.MessageHandler
}

type CommandHandler struct {
	userService usrService.ServiceInterface
}

func (h CommandHandler) GetHandlerFuncs() map[string]tgb.MessageHandler {
	return map[string]tgb.MessageHandler{
		"start": h.StartCommandHandle,
	}
}

func NewCommandHandler(userService usrService.ServiceInterface) *CommandHandler {
	return &CommandHandler{
		userService: userService,
	}
}

func (h CommandHandler) StartCommandHandle(ctx context.Context, msg *tgb.MessageUpdate) error {
	var user models.User
	var err error
	user, err = h.userService.UserGetByChatId(int64(msg.Chat.ID))
	if err != nil {
		return err
	}

	if user.Email == "" {
		return message.NewSendMessageCallBuilder().GetFirstMessage(msg).AddRequestContactKeyboard().Build().DoVoid(ctx)
	}

	return message.NewSendMessageCallBuilder().GetReturnMessage(msg).Build().DoVoid(ctx)
}
