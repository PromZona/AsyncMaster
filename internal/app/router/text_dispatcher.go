package router

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

func DispatchText(context tele.Context, b *bot.BotData) error {
	chatID := context.Chat().ID

	session := b.GetUserSession(chatID)

	if session == nil {
		user, err := db.GetUserByID(b.DB, chatID)
		if err != nil {
			return err
		}
		return ui.MainMenuKeyboard(context, user.Role)
	}

	err := session.DispatchText(context)

	if session.IsDone() {
		delete(b.UserActiveSessions, chatID)
	}
	return err
}
