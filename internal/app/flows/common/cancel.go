package common

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

const CBCancel = "cancel"

func HandleCancelButton(ctx tele.Context, b *bot.BotData) error {
	ctx.Respond()
	chatID := ctx.Chat().ID
	b.ClearUserCache(chatID)

	user, err := db.GetUserByID(b.DB, chatID)
	if err != nil {
		return nil
	}

	return ui.MainMenuKeyboard(ctx, user.Role)
}
