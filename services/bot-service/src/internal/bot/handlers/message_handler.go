package handlers

import (
	"bot-service/internal/bot/message"
	"bot-service/internal/models"
	usrService "bot-service/internal/user"
	"bot-service/internal/utils"
	"context"
	"errors"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"pkg/exceptions"
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
	var user models.User
	user, err := h.userService.UserGetByChatId(int64(update.Message.Chat.ID))
	if errors.Is(err, exceptions.ErrModelNotFound) {
		err = h.userService.UserRegisterByChatId(int64(update.Message.Chat.ID), "Авторегистрация из бота", strings.TrimPrefix(update.Message.Contact.PhoneNumber, "+"))
		if err != nil {
			return err
		}

		return message.NewSendMessageCallBuilder().GetSuccessRegister(update, utils.BuildVlessLink(user.UUID)).RemoveKeyboard().Build().DoVoid(ctx)
	}
	if err != nil {
		//todo нужно сказать пользователю, что у нас ошибка
		return err
	}

	return message.NewSendMessageCallBuilder().GetCustomMessage(update.Answer(tg.HTML.Text(
		tg.HTML.Bold("😻 Вы уже зарегистрированы!"),
		"",
		tg.HTML.Text("Более не нужно делиться номером телефона."),
	)).ParseMode(tg.HTML)).
		RemoveKeyboard().Build().DoVoid(ctx)
}
