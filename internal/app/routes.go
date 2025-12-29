package app

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"

	tele "gopkg.in/telebot.v4"
)

func HandleStartMessage(context tele.Context, bot *BotData) error {
	return context.Send("Hello, enter password to loging into the System")
}

func RegistrationCheckMiddleware(bot *BotData) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(context tele.Context) error {
			chatID := context.Chat().ID

			if !ensureUser(bot.DB, context.Chat().ID) {
				if _, ok := bot.UserSessionState[chatID]; !ok {
					bot.UserSessionState[chatID] = UserStateAwaitPassword
				}
				return HandleRegisterUser(context, bot)
			}

			if _, ok := bot.UserSessionState[chatID]; !ok {
				bot.UserSessionState[chatID] = UserStateDefault
			}
			return next(context)
		}
	}
}

func HandleRegisterUser(context tele.Context, bot *BotData) error {
	chatID := context.Chat().ID
	state, ok := bot.UserSessionState[chatID]
	if !ok {
		state = UserStateAwaitPassword
		bot.UserSessionState[chatID] = state
	}

	switch state {
	case UserStateAwaitPassword:
		password := os.Getenv("BOT_USER_PASSWORD")
		if context.Text() == password {
			context.Send("Password is correct, welcome!")
			bot.UserSessionState[chatID] = UserStateAwaitCodename
			return context.Send("Please enter your Player Name")
		} else {
			return HandleStartMessage(context, bot)
		}
	case UserStateAwaitCodename:
		playerName := context.Text()
		user := &UserData{
			ChatID:       context.Chat().ID,
			TelegramName: context.Sender().FirstName,
			PlayerName:   playerName}

		err := registerUser(bot.DB, user)
		if err != nil {
			log.Print("Error while registering user, ", err)
			return context.Send("Failed to register you, contact administrator")
		}
		bot.UserSessionState[chatID] = UserStateDefault
		return context.Send("Your player name: " + playerName)
	default:
		log.Print("Error in handle register, met unexpected User State: ", state)
		return errors.New("unsupported user state")
	}
}

func HandleSave(context tele.Context, bot *BotData) error {
	chatID := context.Chat().ID
	state, ok := bot.UserSessionState[chatID]
	if !ok {
		state = UserStateDefault
		bot.UserSessionState[chatID] = state
	}

	switch state {
	case UserStateDefault:
		state = UserStateAwaitSavingMessage
		bot.UserSessionState[chatID] = state
		return context.Send("Send message that you want to save")
	case UserStateAwaitSavingMessage:
		messageID := strconv.FormatInt(int64(context.Message().ID), 10)
		message := MessageStruct{MessageID: messageID, ChatID: chatID}
		bot.MessageCache[chatID] = &message
		state = UserStateAwaitTitleForMesssage
		bot.UserSessionState[chatID] = state
		return context.Send("Message received. Give your message a title")
	case UserStateAwaitTitleForMesssage:
		message, ok := bot.MessageCache[chatID]
		if !ok {
			log.Print("Error while creating message: Message is not found")
			bot.UserSessionState[chatID] = UserStateDefault
			return context.Send("Error occured")
		}
		title := context.Message().Text
		message.MessageTitle = title
		createMessage(bot.DB, message)

		bot.UserSessionState[chatID] = UserStateDefault
		return context.Send("Message has been saved")
	default:
		log.Print("Error: Met unsupported state while handling Save action: ", state)
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
		if err != nil {
			return err
		}
	} else {
		playerNames = append(playerNames, args[0])
	}

	messageID := args[1]
	for _, playerName := range playerNames {
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
	chatID := context.Chat().ID
	state := bot.UserSessionState[chatID]

	switch state {
	case UserStateDefault:
		return context.Send("What do you want from me?", bot.PlayerMenu)
	case UserStateAwaitSavingMessage, UserStateAwaitTitleForMesssage:
		return HandleSave(context, bot)
	}
	return nil
}
