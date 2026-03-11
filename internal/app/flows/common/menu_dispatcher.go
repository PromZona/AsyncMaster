package common

import (
	"database/sql"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
)

func GetMainMenuByRole(context runtime.Context, DB *sql.DB, user *bot.UserData) error {
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
