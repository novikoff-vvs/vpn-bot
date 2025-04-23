package handlers

import (
	"bot-service/internal/singleton"
	usrService "bot-service/internal/user"
	"bot-service/internal/utils"
	"context"
	"errors"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"pkg/exceptions"
	"pkg/models"
	"strings"
)

type MessageHandlerInterface interface {
	GetHandlerFuncs() []func() (tgb.MessageHandler, []tgb.Filter)
}

type MessageHandler struct {
	userService usrService.ServiceInterface
}

func NewMessageHandler(userService usrService.ServiceInterface) *MessageHandler {
	return &MessageHandler{userService: userService}
}

type ContactFilter struct {
}

func (f ContactFilter) Allow(ctx context.Context, update *tgb.Update) (bool, error) {
	return update.Message.Contact != nil, nil
}

func (h MessageHandler) GetHandlerFuncs() []func() (tgb.MessageHandler, []tgb.Filter) {
	return []func() (tgb.MessageHandler, []tgb.Filter){
		func() (tgb.MessageHandler, []tgb.Filter) {
			return h.ContactHandle, []tgb.Filter{ContactFilter{}}
		},
	}
}

func (h MessageHandler) ContactHandle(ctx context.Context, update *tgb.MessageUpdate) error {
	var user models.VpnUser
	user, err := h.userService.UserGetByChatId(int64(update.Chat.ID))
	if errors.Is(err, exceptions.ErrModelNotFound) {
		user, err = h.userService.UserRegisterByChatId(int64(update.Chat.ID), "Авторегистрация из бота", strings.TrimPrefix(update.Contact.PhoneNumber, "+"))
		if err != nil {
			return err
		}

		return singleton.MessageBuilder().GetSuccessRegister(update, utils.BuildVlessLink(user.UUID)).AddRequestMainMenuKeyboard("POEBEN").RemoveKeyboard().Build().DoVoid(ctx)
	}
	if err != nil {
		//todo нужно сказать пользователю, что у нас ошибка
		return err
	}

	return singleton.MessageBuilder().GetCustomMessage(update.Answer(tg.HTML.Text(
		tg.HTML.Bold("😻 Вы уже зарегистрированы!"),
		"",
		tg.HTML.Text("Более не нужно делиться номером телефона."),
	)).ParseMode(tg.HTML)).
		RemoveKeyboard().Build().DoVoid(ctx)
}
