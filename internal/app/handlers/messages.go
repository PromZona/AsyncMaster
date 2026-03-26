package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
)

func InitialSendMessage(context runtime.Context, s *bot.Session) error {
	chatID := context.ChatID()

	s.DraftTransaction = &bot.MessageTransaction{
		From:    chatID,
		To:      nil,
		Message: nil,
	}
	s.DraftTransaction.From = chatID

	if s.IsSendEveryone {
		_, ids, err := db.GetNamesAndChatIDsOfAll(s.DB)
		if err != nil {
			return err
		}

		chatids := make([]int64, len(ids))
		copy(chatids, ids)
		s.DraftTransaction.To = chatids
		return context.Send("Write your message:")
	}

	playerNames, chatIDs, err := db.GetNamesAndChatIDsOfAll(s.DB)
	if err != nil {
		context.Send("Error happened while processing your request, contact administrator")
		return err
	}
	return context.Send("Names:", ui.PlayerNamesKeyboard(playerNames, chatIDs))
}

func InitialSendEveryone(context runtime.Context, s *bot.Session) error {
	s.IsSendEveryone = true
	return InitialSendMessage(context, s)
}

func PlayerNameProcess(context runtime.Context, s *bot.Session) error {
	_, cbData := ParseCallbackDataString(context.Callback())
	splited := strings.SplitAfterN(cbData, ":", 2)
	if len(splited) != 2 {
		return fmt.Errorf("splitting callbackdata, met unexpected amount of data: %d", len(splited))
	}

	toChatID, err := strconv.ParseInt(splited[1], 10, 64)
	if err != nil {
		return err
	}

	s.DraftTransaction.To = append(s.DraftTransaction.To, toChatID)
	return context.Send("Write your message:")
}

func MessageTextProcess(context runtime.Context, s *bot.Session) error {
	chatID := context.ChatID()

	message := &bot.Message{
		Title:     "",
		MessageID: strconv.FormatInt(context.MessageID(), 10),
		ChatID:    chatID,
		Text:      context.MessageText(),
	}

	s.DraftMessage = message
	return context.Send("Do you want to add title for a message?", ui.YesNoKeyboard(bot.HID_message_title_yes, bot.HID_message_title_no))
}

func MessageTitleProcess(context runtime.Context, s *bot.Session) error {
	s.DraftMessage.Title = context.MessageText()
	return messageCreationFinilize(context, s)
}

func MessageYesTitle(context runtime.Context, s *bot.Session) error {
	return context.Send("Write title for your message:")
}

func MessageNoTitle(context runtime.Context, s *bot.Session) error {
	return messageCreationFinilize(context, s)
}

func messageCreationFinilize(context runtime.Context, s *bot.Session) error {
	chatID := context.ChatID()

	transaction := s.DraftTransaction
	message := s.DraftMessage
	if transaction == nil || message == nil {
		return fmt.Errorf("expected transaction message and message to exist while sending message")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	message, err = db.CreateMessage(tx, message)
	if err != nil {
		tx.Rollback()
		return err
	}

	transaction.Message = message
	transaction, err = db.CreateMesssageTransaction(tx, transaction)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	userFrom, err := db.GetUserByID(s.DB, int64(transaction.From))
	if err != nil {
		return err
	}

	for _, toChatID := range transaction.To {
		messageFromPlayerName := userFrom.PlayerName

		formatedMessage := fmt.Sprintf("Title: %s\n\nFrom: %s\n\n %s",
			message.Title,
			messageFromPlayerName,
			message.Text)

		context.SendTo(toChatID, formatedMessage)

		log.Printf("Send message succesfully. From: %d to %d, transaction id: %d", userFrom.ChatID, toChatID, transaction.ID)
	}

	context.Send("Message sent")
	user, err := db.GetUserByID(s.DB, chatID)
	if err != nil {
		return err
	}
	return GetMainMenuByRole(context, s.DB, user)
}

func ListMessages(context runtime.Context, s *bot.Session) error {
	messages, err := db.GetLastMessageTransactions(s.DB, context.ChatID())
	if err != nil {
		return err
	}

	return context.Send("Your last 10 messages, pick one", ui.UserMessagesKeyboard(messages))
}

func GetMessage(context runtime.Context, s *bot.Session) error {
	_, cbData := ParseCallbackDataString(context.Callback())
	transactionID, err := strconv.ParseInt(cbData, 10, 64)
	if err != nil {
		return err
	}

	transaction, err := db.GetMessageTransaction(s.DB, transactionID)
	if err != nil {
		return err
	}

	user_from, err := db.GetUserByID(s.DB, int64(transaction.From))
	if err != nil {
		return err
	}
	messageFromPlayerName := user_from.PlayerName

	formatedMessage := fmt.Sprintf("Title: %s\n\nFrom: %s\n\n %s",
		transaction.Message.Title,
		messageFromPlayerName,
		transaction.Message.Text)

	return context.Send(formatedMessage)
}
