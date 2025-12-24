package app

import (
	"log"
	"os"

	tele "gopkg.in/telebot.v4"
)

func HandleStartMessage(context tele.Context, bot *BotData) error {
	return context.Send("Hello, enter password to loging into the System")
}

func MasterSendMessageToAll(context tele.Context, bot *BotData) error {
	chatId := context.Chat().ID
	user := getUser(bot.DB, chatId)
	if user == nil {
		return HandleStartMessage(context, bot)
	}

	user.State = UserStateSendingAll
	updateUser(bot.DB, user)
	return context.Send("What to send?")
}

func HandleText(context tele.Context, bot *BotData) error {
	id := context.Chat().ID
	user := getUser(bot.DB, id)
	if user == nil {
		password := os.Getenv("BOT_USER_PASSWORD")
		if context.Text() == password {
			_ = registerUser(bot.DB, context)
			context.Send("Password is correct, welcome!")
			return context.Send("Please enter your Player Name")
		} else {
			return HandleStartMessage(context, bot)
		}
	}

	switch user.State {
	case UserStateDefault:
		return context.Send("What do you want from me?")
	case UserStateCodename:
		playerName := context.Text()
		user.PlayerName = playerName
		user.State = UserStateDefault
		updateUser(bot.DB, user)
		return context.Send("Your player name:" + playerName)
	case UserStateSendingAll:
		message := context.Message().Text
		log.Print(string(message))
		user.State = UserStateDefault
		updateUser(bot.DB, user)

		return context.Send("Thanks")
	}
	return nil
}
