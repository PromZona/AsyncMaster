package router

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/flows/common"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func DispatchText(context runtime.Context, b *bot.BotData) error {
	chatID := context.ChatID()

	session := b.GetUserSession(chatID)

	if session == nil {
		user, err := db.GetUserByID(b.DB, chatID)
		if err != nil {
			return err
		}
		return common.GetMainMenuByRole(context, b.DB, user)
	}

	err := session.DispatchText(context)

	if session.IsDone() {
		delete(b.UserActiveSessions, chatID)
	}
	return err
}
