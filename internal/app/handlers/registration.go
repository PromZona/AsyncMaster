package handlers

import (
	"fmt"
	"os"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func RegistrationStartMessage(context runtime.Context) error {
	return context.Send("Hello, enter password to log in into the System")
}

func RegistrationPassword(context runtime.Context, s *bot.Session) error {
	passwordPlayer := os.Getenv("BOT_USER_PASSWORD")
	passwordMaster := os.Getenv("BOT_MASTER_PASSWORD")

	password := context.MessageText()
	if password != passwordMaster && password != passwordPlayer {
		return RegistrationStartMessage(context)
	}

	s.RegistrationUser = &bot.UserData{
		ChatID:       context.ChatID(),
		TelegramName: context.FirstName(),
		PlayerName:   "",
		Faction:      &bot.Faction{},
	}
	s.RegistrationState = bot.RegistrationAwaitCodename

	switch password {
	case passwordPlayer:
		context.Send("Player password is correct, welcome!")
		s.RegistrationUser.Role = bot.RolePlayer
		return context.Send("Please enter your Player Name")
	case passwordMaster:
		context.Send("Master password is correct, welcome!")
		s.RegistrationUser.Role = bot.RoleMaster
		return context.Send("Please, enter your Master Name")
	default:
		return fmt.Errorf("expected to handle password while registering user, but received this: %s", password)
	}
}

func RegistrationPlayerName(context runtime.Context, s *bot.Session) error {
	playerName := context.MessageText()

	s.RegistrationUser.PlayerName = playerName
	if s.RegistrationUser.Role == bot.RoleMaster {
		return registrationFinilize(context, s)
	}

	s.RegistrationState = bot.RegistrationAwaitFactionName
	return context.Send("In this game you control a faction of your own. And you charachter is a leader\nNow you need to create your faction\nWrite name for a faction:")
}

func RegistrationFactionName(context runtime.Context, s *bot.Session) error {
	factionName := context.MessageText()

	s.RegistrationUser.Faction.Name = factionName
	s.RegistrationState = bot.RegistrationAwaitFactionDescription
	return context.Send("Now describe your faction. 1 paragraph of text:")
}

func RegistrationFactionDescription(context runtime.Context, s *bot.Session) error {
	factionDesc := context.MessageText()

	s.RegistrationUser.Faction.Description = factionDesc
	return registrationFinilize(context, s)
}

func registrationFinilize(context runtime.Context, s *bot.Session) error {
	id, err := db.CreateUser(s.DB, s.RegistrationUser)
	if err != nil {
		return err
	}

	_, err = db.CreateFaction(s.DB, s.RegistrationUser.Faction, id)
	if err != nil {
		return err
	}

	s.RegistrationState = bot.RegistrationFinished
	return context.Send("You are ready...\nMaster will contact you soon")
}
