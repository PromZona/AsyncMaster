package app

import (
	"database/sql"
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
	// /send TO WHAT
	// TO can be player_name or EVERYONE
	// WHAT is message_id
	args := context.Args()
	if len(args) < 2 {
		return context.Send("Expected 2 arguments: TO(player_name or EVERYONE) and WHAT(message_id)")
	}

	var (
		playerNames []string
		err         error
	)

	if args[0] == "EVERYONE" {
		playerNames, err = getUserPlayerNames(bot.DB)
		log.Print("DEBUG ", playerNames)
		if err != nil {
			return err
		}
	} else {
		playerNames = append(playerNames, args[0])
	}

	messageID := args[1]
	for _, playerName := range playerNames {
		log.Print("DEBUG ", playerName)
		user, err := getUserByName(bot.DB, playerName)
		if err != nil {
			if err == sql.ErrNoRows {
				context.Send("Given player name not found")
				continue
			} else {
				log.Print("Error while searching player by name ", playerName, " ", err)
				context.Send("Error occured while searching for given player")
				continue
			}
		}

		message, err := getMessage(bot.DB, messageID)
		if err != nil {
			if err == sql.ErrNoRows {
				context.Send("Given messsage does not registered in the DB")
				continue
			} else {
				log.Print("Error while searching message by message_id ", messageID, " ", err)
				context.Send("Error occured while searching for message to send")
				continue
			}
		}

		_, err = context.Bot().Copy(user, message)
		if err != nil {
			log.Print("Error occured while coping message ", err)
			context.Send("Failed to send requested message")
			continue
		}

		log.Print("Sent message [", messageID, "] to [", playerName, "]")
	}
	return nil
}

func HandleText(context tele.Context, bot *BotData) error {
	id := context.Chat().ID
	user := getUser(bot.DB, id)
	if user == nil {
		password := os.Getenv("BOT_USER_PASSWORD")
		if context.Text() == password {
			user := registerUser(bot.DB, context)
			context.Send("Password is correct, welcome!")
			user.State = UserStateAwaitCodename
			updateUser(bot.DB, user)
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
