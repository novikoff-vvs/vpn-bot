package handlers

import (
	"bot-service/internal/singleton"
	usrService "bot-service/internal/user"
	"context"
	"errors"
	"github.com/mr-linch/go-tg/tgb"
	"pkg/exceptions"
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
	var err error
	_, err = h.userService.UserGetByChatId(int64(msg.Chat.ID))
	if errors.Is(err, exceptions.ErrModelNotFound) {
		return singleton.MessageBuilder().GetFirstMessage(msg).AddRequestContactKeyboard().Build().DoVoid(ctx)
	}
	if err != nil {
		return err
	}

	return singleton.MessageBuilder().GetReturnMessage(msg).AddRequestMainMenuKeyboard().Build().DoVoid(ctx)
}
