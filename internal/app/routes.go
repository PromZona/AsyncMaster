package app

import (
	"log"

	tele "gopkg.in/telebot.v4"
)

func MasterSendMessageToAll(context tele.Context, bot *BotData) error {

	id := context.Chat().ID
	isUserExist := ensureUser(bot.DB, bot, id)
	var user *UserData
	if !isUserExist {
		log.Print("Not found user: ", id)
		user = registerUser(bot.DB, context)
	} else {
		user = getUser(bot.DB, id)
	}
	user.State = UserStateSendingAll
	updateUser(bot.DB, user)
	return context.Send("What to send?")
}

func HandleText(context tele.Context, bot *BotData) error {

	return nil

	id := context.Chat().ID
	user := bot.Users[id]
	switch user.State {
	case UserStateDefault:
		return context.Send("What do you want from me?")
	case UserStateSendingAll:
		message := context.Message().Text
		log.Print(string(message))
		return context.Send("Thanks")
	}
	return nil
}
