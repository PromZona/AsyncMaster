package registration

import (
	"log"
	"os"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	tele "gopkg.in/telebot.v4"
)

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
	}

	err := db.CreateUser(s.DB, user)
	if err != nil {
		log.Print("Error while registering user, ", err)
		return context.Send("Failed to register you, contact administrator")
	}
	s.Done = true
	return context.Send("Your player name: " + playerName)
}

func StartMessage(context tele.Context) error {
	return context.Send("Hello, enter password to loging into the System")
}
