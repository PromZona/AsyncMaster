package app

import (
	"log"
	"os"
	"strconv"

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

func HandleSave(context tele.Context, bot *BotData) error {
	chatID := context.Chat().ID
	user := getUser(bot.DB, chatID)

	switch user.State {
	case UserStateDefault:
		user.State = UserStateAwaitSavingMessage
		updateUser(bot.DB, user)
		return context.Send("Send message that you want to save")
	case UserStateAwaitSavingMessage:
		messageID := strconv.FormatInt(int64(context.Message().ID), 10)
		message := MessageStruct{MessageID: messageID, ChatID: chatID}
		bot.MessageCache[chatID] = &message
		user.State = UserStateAwaitTitleForMesssage
		updateUser(bot.DB, user)
		return context.Send("Message received. Give your message a title")
	case UserStateAwaitTitleForMesssage:
		message, ok := bot.MessageCache[chatID]
		if !ok {
			log.Print("Error while creating message: Message is not found")
			return context.Send("Error occured")
		}
		title := context.Message().Text
		message.MessageTitle = title
		createMessage(bot.DB, message)

		user.State = UserStateDefault
		updateUser(bot.DB, user)
		return context.Send("Message has been saved. ID")
	default:
		log.Print("Error: Met unsupported state while handling Save action: ", user.State)
	}

	return nil
}

func HandleSend(context tele.Context, bot *BotData) error {
	// args := context.Args()

	return nil
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
	case UserStateAwaitCodename:
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
	case UserStateAwaitSavingMessage, UserStateAwaitTitleForMesssage:
		return HandleSave(context, bot)
	}
	return nil
}
