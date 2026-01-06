package ui

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	tele "gopkg.in/telebot.v4"
)

func HandleCancelButton(ctx tele.Context, b *bot.BotData) error {
	ctx.Respond()
	chatID := ctx.Chat().ID
	b.ClearUserCache(chatID)
	b.UserSessionState[chatID] = bot.UserStateDefault
	return MainMenu(ctx, b)
}
