package handlers

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

type ReactionHandlerInterface interface {
	GetReactionHandleFunc() tgb.MessageReactionHandler
}

type ReactionHandler struct {
	bot *tg.Client
}

func NewReactionHandler(bot *tg.Client) *ReactionHandler {
	return &ReactionHandler{
		bot: bot,
	}
}

func (r ReactionHandler) GetReactionHandleFunc() tgb.MessageReactionHandler {
	return r.MessageReaction
}

func (r ReactionHandler) MessageReaction(ctx context.Context, reaction *tgb.MessageReactionUpdate) error {
	if len(reaction.OldReaction) > len(reaction.NewReaction) {
		return nil
	}
	return r.bot.SendMessage(reaction.Chat.ID, "You are breathtaking!").DoVoid(ctx)
}
