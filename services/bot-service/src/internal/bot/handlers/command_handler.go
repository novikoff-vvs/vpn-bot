package handlers

import (
	notify_user "bot-service/internal/repository/pgsql/notify-user"
	"bot-service/internal/singleton"
	usrService "bot-service/internal/user"
	"context"
	"errors"
	"github.com/mr-linch/go-tg/tgb"
	"log"
	"pkg/exceptions"
)

type CommandHandlerInterface interface {
	GetHandlerFuncs() map[string]tgb.MessageHandler
}

type CommandHandler struct {
	userService    usrService.ServiceInterface
	notifyUserRepo *notify_user.NotifyUserRepository
}

func (h CommandHandler) GetHandlerFuncs() map[string]tgb.MessageHandler {
	return map[string]tgb.MessageHandler{
		"start":       h.StartCommandHandle,
		"instruction": h.InstructionCommandHandle,
	}
}

func NewCommandHandler(userService usrService.ServiceInterface, notifyUserRepo *notify_user.NotifyUserRepository) *CommandHandler {
	return &CommandHandler{
		userService:    userService,
		notifyUserRepo: notifyUserRepo,
	}
}

func (h CommandHandler) StartCommandHandle(ctx context.Context, msg *tgb.MessageUpdate) error {
	var err error
	usr, err := h.userService.UserGetByChatId(int64(msg.Chat.ID))
	if errors.Is(err, exceptions.ErrModelNotFound) {
		return singleton.MessageBuilder().GetFirstMessage(msg).AddRequestContactKeyboard().Build().DoVoid(ctx)
	}
	if err != nil {
		return err
	}

	err = h.notifyUserRepo.Create(usr.ChatId)
	if err != nil {
		log.Println(err.Error())
	}

	return singleton.MessageBuilder().GetReturnMessage(msg).AddRequestMainMenuKeyboard(usr.UUID).Build().DoVoid(ctx)
}

func (h CommandHandler) InstructionCommandHandle(ctx context.Context, msg *tgb.MessageUpdate) error {

	return singleton.MessageBuilder().GetInstructionMessage(msg).Build().DoVoid(ctx)
}
