package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
)

func ListFactions(context runtime.Context, s *bot.Session) error {

	data, err := db.GetUsersAll(s.DB)
	if err != nil {
		return err
	}

	var text strings.Builder
	text.WriteString("Фракции\n\n")

	for _, user := range data {
		faction := fmt.Sprintf(`
			Фракия: %s
			Лидер Фракции: %s
			Описание:
			%s
			---
			`,
			user.Faction.Name,
			user.PlayerName,
			user.Faction.Description,
		)

		text.WriteString(faction)
	}

	return context.Send(text.String())
}

func FactionsUpdate(context runtime.Context, s *bot.Session) error {
	users, err := db.GetUsersAll(s.DB)
	if err != nil {
		return err
	}

	return context.Send("Pick a faction to change", ui.FactionsUpdateFactionsList(users))
}

func FactionUpdatePickWhat(context runtime.Context, s *bot.Session) error {
	_, cbData := ParseCallbackDataString(context.Callback())
	chatID, err := strconv.ParseInt(cbData, 10, 64)
	if err != nil {
		return err
	}

	user, err := db.GetUserByID(s.DB, chatID)
	if err != nil {
		return err
	}

	s.FactionUpdateUser = user
	text := fmt.Sprintf(`
		What to update
		
		Player Name : %s
		Faction Name: %s
		Faction Description:
		%s
		Faction Resources:
		%s`,
		user.PlayerName,
		user.Faction.Name,
		user.Faction.Description,
		user.Faction.Resources)

	return context.Send(text, ui.FactionsUpdateWhatToUpdate())
}

func FactionUpdateResources(context runtime.Context, s *bot.Session) error {
	text := fmt.Sprintf("Current resources:\n\n%s\n\nEnter new resources", s.FactionUpdateUser.Faction.Resources)
	return context.Send(text, ui.CancelMenu())
}

func FactionUpdateResourcesText(context runtime.Context, s *bot.Session) error {
	text := context.MessageText()
	s.FactionUpdateUser.Faction.Resources = text

	err := db.UpdateFaction(s.DB, int64(s.FactionUpdateUser.Faction.ID), s.FactionUpdateUser.Faction)
	if err != nil {
		return err
	}

	return context.Send("Updated successfuly")
}
