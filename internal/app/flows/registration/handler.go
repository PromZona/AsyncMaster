package registration

import (
	"errors"
	"log"
	"os"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	tele "gopkg.in/telebot.v4"
)

func HandleRegisterUser(context tele.Context, b *bot.BotData) error {
	chatID := context.Chat().ID
	state, ok := b.UserSessionState[chatID]
	if !ok {
		state = bot.UserStateAwaitPassword
		b.UserSessionState[chatID] = state
	}

	switch state {
	case bot.UserStateAwaitPassword:
		password := os.Getenv("BOT_USER_PASSWORD")
		if context.Text() == password {
			context.Send("Password is correct, welcome!")
			b.UserSessionState[chatID] = bot.UserStateAwaitCodename
			return context.Send("Please enter your Player Name")
		} else {
			return StartMessage(context)
		}
	case bot.UserStateAwaitCodename:
		playerName := context.Text()
		user := &bot.UserData{
			ChatID:       context.Chat().ID,
			TelegramName: context.Sender().FirstName,
			PlayerName:   playerName,
			Role:         bot.RolePlayer,
		}

		err := db.RegisterUser(b.DB, user)
		if err != nil {
			log.Print("Error while registering user, ", err)
			return context.Send("Failed to register you, contact administrator")
		}
		b.UserSessionState[chatID] = bot.UserStateDefault
		return context.Send("Your player name: " + playerName)
	default:
		log.Print("Error in handle register, met unexpected User State: ", state)
		return errors.New("unsupported user state")
	}
}

func StartMessage(context tele.Context) error {
	return context.Send("Hello, enter password to loging into the System")
}
