package common

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

func HandleCancelButton(ctx tele.Context, b *bot.BotData) error {
	ctx.Respond()
	chatID := ctx.Chat().ID
	b.ClearUserCache(chatID)
	b.UserSessionState[chatID] = bot.UserStateDefault
	return ui.MainMenu(ctx, b)
}
