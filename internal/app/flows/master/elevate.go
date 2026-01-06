package master

import (
	"os"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

func HandleElevateToMaster(context tele.Context, b *bot.BotData) error {
	args := context.Args()
	chatID := context.Chat().ID
	state := b.UserSessionState[chatID]

	if state != bot.UserStateDefault {
		return context.Send("Please finish previous action to activate this command")
	}

	if len(args) < 1 {
		return context.Send("Not enough arguments received. Send a password as argument for a command")
	}

	password := os.Getenv("BOT_MASTER_PASSWORD")
	if password != args[0] {
		return context.Send("Password is incorrect")
	}

	user := db.GetUserByID(b.DB, chatID)
	user.Role = bot.RoleMaster
	db.UpdateUser(b.DB, user)

	return context.Send("Role updated to Master", ui.MainMenuKeyboard(context, b))
}
