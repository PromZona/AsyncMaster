package handlers

import (
	"database/sql"
	"os"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
)

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

func HandleElevateToMaster(context runtime.Context, b *bot.BotData) error {
	args := context.Args()
	chatID := context.ChatID()

	if len(args) < 1 {
		return context.Send("Not enough arguments received. Send a password as argument for a command")
	}

	password := os.Getenv("BOT_MASTER_PASSWORD")
	if password != args[0] {
		return context.Send("Password is incorrect")
	}

	user, err := db.GetUserByID(b.DB, chatID)
	if err != nil {
		return err
	}

	user.Role = bot.RoleMaster
	db.UpdateUser(b.DB, user)

	user, err = db.GetUserByID(b.DB, chatID)
	if err != nil {
		return err
	}

	context.Send("Role updated to Master")
	return GetMainMenuByRole(context, b.DB, user)
}

func ParseCallbackDataString(callbackData string) (unique, data string) {
	trimmed := strings.Trim(callbackData, "\f")
	parts := strings.SplitN(trimmed, "|", 2)
	count := len(parts)
	if count == 2 {
		return parts[0], parts[1]
	}
	if count == 1 {
		return parts[0], ""
	}
	return "", ""
}
