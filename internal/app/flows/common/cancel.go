package common

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

const CBCancel = "cancel"

func HandleCancelButton(ctx runtime.Context, b *bot.BotData) error {
	ctx.Respond()
	chatID := ctx.ChatID()
	b.ClearUserCache(chatID)

	user, err := db.GetUserByID(b.DB, chatID)
	if err != nil {
		return nil
	}

	return GetMainMenuByRole(ctx, b.DB, user)
}
