package common

import (
	"database/sql"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/ui"

	tele "gopkg.in/telebot.v4"
)

func GetMainMenuByRole(context tele.Context, DB *sql.DB, user *bot.UserData) error {
	if user.Role == bot.RolePlayer {
		menuData := db.GetPlayerMenuData(DB, user.ChatID)
		return ui.MainMenuPlayerKeyboard(context, user, menuData)
	}

	if user.Role == bot.RoleMaster {
		// TODO: Need to implement this
		// menuData := db.GetMasetMenuData(DB, user.ChatID)
		return ui.MainMenuMasterKeyboard(context, user, nil)
	}
	return context.Send("Met unexpected role")
}
