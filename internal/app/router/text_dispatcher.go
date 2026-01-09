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
		return ui.MainMenuKeyboard(context, db.GetUserByID(b.DB, chatID).Role)
	}

	err := session.DispatchText(context)

	if session.IsDone() {
		delete(b.UserActiveSessions, chatID)
	}
	return err
}
