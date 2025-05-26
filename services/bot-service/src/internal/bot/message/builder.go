package message

import (
	"bot-service/config"
	"fmt"
	"github.com/mr-linch/go-tg"
)

type message interface {
	Answer(text string) *tg.SendMessageCall
}

type Builder struct {
	sndMsg *tg.SendMessageCall
	cfg    config.PaymentService
}

func NewSendMessageCallBuilder(cfg config.PaymentService) *Builder {
	return &Builder{
		cfg: cfg,
	}
}
func (b Builder) GetFirstMessage(msg message) Builder {
	b.sndMsg = msg.Answer(
		tg.HTML.Text(
			tg.HTML.Bold("👋 Добро пожаловать!"),
			"",
			tg.HTML.Text("Для завершения регистрация, пожалуйста, пришлите ваш номер телефона."),
		),
	).ParseMode(tg.HTML)
	return b
}
func (b Builder) GetReturnMessage(msg message) Builder {
	b.sndMsg = msg.Answer(
		tg.HTML.Text(
			tg.HTML.Bold("👋 C возвращением!"),
		),
	).ParseMode(tg.HTML)
	return b
}
func (b Builder) GetSuccessRegister(msg message, vessaLink string) Builder {
	b.sndMsg = msg.Answer(tg.HTML.Text(
		tg.HTML.Bold("🎉 Регистрация успешно завершена!"),
		"",
		tg.HTML.Text(tg.HTML.Bold("🔐 Важно:"), tg.HTML.Blockquote("Эта ссылка является вашим личным доступом.\n\nНикому не передавайте её – это может привести к потере аккаунта.")),
		"",
		tg.HTML.Line(tg.HTML.Bold("🔗 Ваша подписка: "), tg.HTML.Link("SUBSCRIPTION-URL", vessaLink)),
	)).ParseMode(tg.HTML)
	return b
}
func (b Builder) GetVessaLinkMessage(msg message, vessaLink string) Builder {
	b.sndMsg = msg.Answer(tg.HTML.Line(tg.HTML.Bold("🔗 Ваша подписка: "), tg.HTML.Link("SUBSCRIPTION-URL", vessaLink)))
	return b
}
func (b Builder) GetCustomMessage(msg *tg.SendMessageCall) Builder {
	b.sndMsg = msg
	return b
}
func (b Builder) AddRequestContactKeyboard() Builder {
	inlineKeyboard := tg.NewReplyKeyboardMarkup(
		[]tg.KeyboardButton{
			tg.NewKeyboardButtonRequestContact("📱 Отправить номер"),
		},
	).WithResizeKeyboardMarkup()
	b.sndMsg.ReplyMarkup(inlineKeyboard)
	return b
}
func (b Builder) AddRequestMainMenuKeyboard(uuid string) Builder {

	b.sndMsg.ReplyMarkup(b.GetMainMenuKeyboad(uuid))
	return b
}
func (b Builder) GetPaymentMenuKeyboard(uuid string) tg.InlineKeyboardMarkup {
	return tg.NewInlineKeyboardMarkup(
		[]tg.InlineKeyboardButton{
			{
				Text: "💸 Оплата",
				WebApp: &tg.WebAppInfo{
					URL: fmt.Sprintf("%sweb/yoomoney/%s", b.cfg.Url, uuid),
				},
			},
		})
}
func (b Builder) RemoveKeyboard() Builder {
	b.sndMsg = b.sndMsg.ReplyMarkup(tg.NewReplyKeyboardRemove())
	return b
}
func (b Builder) Build() *tg.SendMessageCall {
	return b.sndMsg
}
func (b Builder) GetMainMenuKeyboad(uuid string) tg.InlineKeyboardMarkup {
	keyboard := b.GetPaymentMenuKeyboard(uuid).InlineKeyboard
	return tg.NewInlineKeyboardMarkup(
		[]tg.InlineKeyboardButton{
			{
				Text:         "🔗 Моя ссылка",
				CallbackData: "get_link",
			},
		}, keyboard[0])
}
