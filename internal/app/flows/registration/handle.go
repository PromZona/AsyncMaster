package registration

import (
	"fmt"
	"os"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	tele "gopkg.in/telebot.v4"
)

func StartMessage(context tele.Context) error {
	return context.Send("Hello, enter password to log in into the System")
}

func handlePassword(context tele.Context, s *Session) error {
	passwordPlayer := os.Getenv("BOT_USER_PASSWORD")
	passwordMaster := os.Getenv("BOT_MASTER_PASSWORD")

	password := context.Text()
	if password != passwordMaster && password != passwordPlayer {
		return StartMessage(context)
	}

	s.User = &bot.UserData{
		ChatID:       context.Chat().ID,
		TelegramName: context.Sender().FirstName,
		PlayerName:   "",
		Faction:      &bot.Faction{},
	}
	s.UserState = AwaitCodename

	switch password {
	case passwordPlayer:
		context.Send("Player password is correct, welcome!")
		s.User.Role = bot.RolePlayer
		return context.Send("Please enter your Player Name")
	case passwordMaster:
		context.Send("Master password is correct, welcome!")
		s.User.Role = bot.RoleMaster
		return context.Send("Please, enter your Master Name")
	default:
		return fmt.Errorf("expected to handle password while registering user, but received this: %s", password)
	}
}

func handlePlayerName(context tele.Context, s *Session) error {
	playerName := context.Text()

	s.User.PlayerName = playerName
	if s.User.Role == bot.RoleMaster {
		return finilize(context, s)
	}

	s.UserState = AwaitFactionName
	return context.Send("In this game you control a faction of your own. And you charachter is a leader\nNow you need to create your faction\nWrite name for a faction:")
}

func handleFactionName(context tele.Context, s *Session) error {
	factionName := context.Text()

	s.User.Faction.Name = factionName
	s.UserState = AwaitFactionDescription
	return context.Send("Now describe your faction. 1 paragraph of text:")
}

func handleFactionDescription(context tele.Context, s *Session) error {
	factionDesc := context.Text()

	s.User.Faction.Description = factionDesc
	return finilize(context, s)
}

func finilize(context tele.Context, s *Session) error {
	err := db.CreateUser(s.DB, s.User)
	if err != nil {
		return err
	}

	_, err = db.CreateFaction(s.DB, s.User.Faction)
	if err != nil {
		return err
	}

	s.Done = true
	return context.Send("You are ready...\nMaster will contact you soon")
}
