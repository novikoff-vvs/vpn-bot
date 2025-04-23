package router

import (
	"bot-service/internal/bot/handlers"
	"context"
	"github.com/mr-linch/go-tg/tgb"
)

type Interface interface {
	RegisterCommandHandlers(handler handlers.CommandHandlerInterface)
}

type Router struct {
	router *tgb.Router
}

func NewRouter(router *tgb.Router) *Router {
	return &Router{
		router: router,
	}
}

func (r Router) RegisterCommandHandlers(handler handlers.CommandHandlerInterface) {
	for command, fn := range handler.GetHandlerFuncs() {
		r.router.Message(fn, tgb.Command(command))
	}
}

func (r Router) RegisterReactionHandlers(handler handlers.ReactionHandlerInterface) {
	r.router.MessageReaction(handler.GetReactionHandleFunc(), nil)
}

type Filter struct {
	data string
}

func (f Filter) Allow(ctx context.Context, update *tgb.Update) (bool, error) {
	if update.Update.CallbackQuery.Data == f.data {
		return true, nil
	}
	return false, nil
}

func (r Router) RegisterCallbackQueryHandlers(handler handlers.CallbackQueryHandlerInterface) {
	for data, fn := range handler.GetCallbackQueryHandlersFunc() {

		filters := Filter{
			data: data,
		}
		r.router.
			CallbackQuery(
				fn,
				filters,
			)
	}

}

func (r Router) RegisterMessageHandlers(handler handlers.MessageHandlerInterface) {
	for _, get := range handler.GetHandlerFuncs() {
		fn, filters := get()
		r.router.Message(fn, filters...)
	}
}
