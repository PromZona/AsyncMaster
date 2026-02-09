package registration

import (
	"os"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	tele "gopkg.in/telebot.v4"
)

func StartMessage(context tele.Context) error {
	return context.Send("Hello, enter password to log in into the System")
}

func handlePassword(context tele.Context, s *Session) error {
	password := os.Getenv("BOT_USER_PASSWORD")
	if context.Text() == password {
		context.Send("Password is correct, welcome!")
		s.UserState = AwaitCodename
		return context.Send("Please enter your Player Name")
	} else {
		return StartMessage(context)
	}
}

func handlePlayerName(context tele.Context, s *Session) error {
	playerName := context.Text()
	user := &bot.UserData{
		ChatID:       context.Chat().ID,
		TelegramName: context.Sender().FirstName,
		PlayerName:   playerName,
		Role:         bot.RolePlayer,
		Faction:      &bot.Faction{},
	}

	s.User = user
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
	return nil
}
