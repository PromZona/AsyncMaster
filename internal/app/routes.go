package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func HandleStartMessage(context tele.Context, bot *BotData) error {
	return context.Send("Hello, enter password to loging into the System")
}

func SendMainMenu(context tele.Context, bot *BotData) error {
	return context.Send("Main menu", bot.PlayerMenu)
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

func ErrorRecoveryMiddleware(bot *BotData) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(context tele.Context) error {
			err := next(context)
			if err != nil {
				log.Print("ERROR: ", err)
				bot.UserSessionState[context.Chat().ID] = UserStateDefault
				context.Send("Error occured while proccesing your request. Please, contact administrator. Returning you to main menu")
			}
			return err
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
			PlayerName:   playerName,
			Role:         RolePlayer,
		}

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
		message := Message{MessageID: messageID, ChatID: chatID}
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
		message.Title = title
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

func HandleElevateToMaster(context tele.Context, bot *BotData) error {
	args := context.Args()
	chatID := context.Chat().ID
	state := bot.UserSessionState[chatID]

	if state != UserStateDefault {
		return context.Send("Please finish previous action to activate this command")
	}

	if len(args) < 1 {
		return context.Send("Not enough arguments received. Send a password as argument for a command")
	}

	password := os.Getenv("BOT_MASTER_PASSWORD")
	if password != args[0] {
		return context.Send("Password is incorrect")
	}

	user := getUser(bot.DB, chatID)
	user.Role = RoleMaster
	updateUser(bot.DB, user)

	return context.Send("Role updated to Master")
}

func HandleText(context tele.Context, bot *BotData) error {
	chatID := context.Chat().ID
	state := bot.UserSessionState[chatID]

	switch state {
	case UserStateDefault:
		return SendMainMenu(context, bot)
	case UserStateAwaitSavingMessage, UserStateAwaitTitleForMesssage:
		return HandleSave(context, bot)
	case UserStateAwaitMessage:
		return handlePlayerMessage(context, bot)
	case UserStateAwaitTitle:
		bot.MessageCache[chatID].Title = context.Message().Text
		return handlePlayerFinalMessageSending(context, bot)

	}
	return nil
}

func handlePlayerMessage(context tele.Context, bot *BotData) error {
	chatID := context.Chat().ID

	message := &Message{
		Title:     "",
		MessageID: strconv.FormatInt(int64(context.Message().ID), 10),
		ChatID:    chatID,
		Text:      context.Text(),
	}

	bot.MessageCache[chatID] = message
	bot.UserSessionState[chatID] = UserStateAwaitTitleDecision

	return context.Send("Do you want to add title for a message?", getTitleQuestionKeyboard(bot))
}

func handlePlayerFinalMessageSending(context tele.Context, bot *BotData) error {
	chatID := context.Chat().ID

	transaction, ok := bot.MessageTransactionCache[chatID]
	if !ok {
		return fmt.Errorf("expected transaction message to exist while handling UserStateAwaitMessage state")
	}

	message, ok := bot.MessageCache[chatID]
	if !ok {
		return fmt.Errorf("expected message to exist while handling UserStateAwaitMessage state")
	}

	tx, err := bot.DB.Begin()
	if err != nil {
		return err
	}

	message, err = createMessage(tx, message)
	if err != nil {
		tx.Rollback()
		return err
	}

	transaction.Message = message
	transaction, err = createTransaction(tx, transaction)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	messageFromPlayerName := getUser(bot.DB, int64(transaction.From)).PlayerName
	messageToPlayerName := getUser(bot.DB, int64(transaction.To)).PlayerName

	formatedMessage := fmt.Sprintf("Title: %s\n\nFrom: %s\nTo: %s\n\n %s",
		message.Title,
		messageFromPlayerName,
		messageToPlayerName,
		message.Text)

	context.Bot().Send(transaction.To, formatedMessage)

	log.Printf("Player send message succesfully, transaction id: %d", transaction.ID)

	bot.UserSessionState[chatID] = UserStateDefault
	bot.ClearUserCache(chatID)

	context.Send("Message sent")
	return SendMainMenu(context, bot)
}

func getPlayerNamesKeyboard(bot *BotData) *tele.ReplyMarkup {
	result := &tele.ReplyMarkup{}

	playerNames, chatIDs, err := getUserPlayerNamesAndChatID(bot.DB)
	if err != nil {
		log.Print("Error while creating keyboard, ", err)
		return nil
	}
	if len(playerNames) != len(chatIDs) {
		log.Print("Error while creating keyboard, playerNames are not the same size as chatIDs: ", len(playerNames), "; ", len(chatIDs))
		return nil
	}

	var btnPlayerNames []tele.Btn

	for i, name := range playerNames {
		dataString := fmt.Sprintf("%s:%d", name, chatIDs[i])
		btnPlayerNames = append(btnPlayerNames, result.Data(name, "player_names", dataString))
	}

	result.Inline(
		result.Row(btnPlayerNames...),
		result.Row(bot.BtnCancel),
	)

	return result
}

func getTitleQuestionKeyboard(bot *BotData) *tele.ReplyMarkup {
	result := &tele.ReplyMarkup{}

	btnNo := result.Data("No", "no_title")
	btnYes := result.Data("Yes", "yes_title")

	result.Inline(
		result.Row(btnNo, btnYes),
		result.Row(bot.BtnCancel),
	)

	return result
}

func HandleCallbacks(context tele.Context, bot *BotData) error {
	context.Respond()

	chatID := context.Chat().ID
	state := bot.UserSessionState[chatID]

	rawCallbackData := context.Callback().Data
	cbUnique, cbData := parseCallbackDataString(rawCallbackData)

	// log.Printf("callback unique   = %q", cbUnique)
	// log.Printf("callback data     = %q", cbData)
	// log.Printf("raw callback data = %q", rawCallbackData)

	switch cbUnique {
	case "send":
		if state != UserStateDefault {
			return context.Send("This button is not available right now, please finish your previous action")
		}
		bot.UserSessionState[chatID] = UserStateAwaitResipient
		bot.MessageTransactionCache[chatID] = &MessageTransaction{
			From: tele.ChatID(chatID),
		}
		return context.Send("Names:", getPlayerNamesKeyboard(bot))
	case "player_names":
		if state != UserStateAwaitResipient {
			return context.Send("This button is not available right now, please finish your previous action")
		}

		splited := strings.SplitAfterN(cbData, ":", 2)
		if len(splited) != 2 {
			return fmt.Errorf("splitting callbackdata, met unexpected amount of data: %d", len(splited))
		}

		transaction, ok := bot.MessageTransactionCache[chatID]
		if !ok {
			return fmt.Errorf("retrieving transaction message. No message entry found")
		}

		toChatID, err := strconv.ParseInt(splited[1], 10, 64)
		if err != nil {
			return err
		}

		transaction.To = tele.ChatID(toChatID)
		bot.UserSessionState[chatID] = UserStateAwaitMessage
		return context.Send("Write your message:")
	case "no_title":
		if state != UserStateAwaitTitleDecision {
			return context.Send("This button is not available right now, please finish your previous action")
		}
		return handlePlayerFinalMessageSending(context, bot)
	case "yes_title":
		if state != UserStateAwaitTitleDecision {
			return context.Send("This button is not available right now, please finish your previous action")
		}
		bot.UserSessionState[chatID] = UserStateAwaitTitle
		return context.Send("Write title for your message:")
	default:
		log.Print("Error, met unknown state while receiving callback, ", state)
		return errors.New("unsupported state in callback")
	}
}
