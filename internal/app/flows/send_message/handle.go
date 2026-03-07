package sendmessage

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/flows/common"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	"gopkg.in/telebot.v4"
	tele "gopkg.in/telebot.v4"
)

func handleMessageText(context tele.Context, s *Session) error {
	chatID := context.Chat().ID

	message := &bot.Message{
		Title:     "",
		MessageID: strconv.FormatInt(int64(context.Message().ID), 10),
		ChatID:    chatID,
		Text:      context.Text(),
	}

	s.DraftMessage = message
	s.UserState = AwaitTitleDecision

	return context.Send("Do you want to add title for a message?", ui.YesNoKeyboard())
}

func handleMessageTitle(context tele.Context, s *Session) error {
	s.DraftMessage.Title = context.Message().Text
	return finilize(context, s)
}

func handleInitialSend(context tele.Context, s *Session) error {
	chatID := context.Chat().ID

	s.DraftTransaction.From = tele.ChatID(chatID)

	if s.IsSendEveryone {
		_, ids, err := db.GetNamesAndChatIDsOfAll(s.DB)
		if err != nil {
			return err
		}

		chatids := make([]telebot.ChatID, len(ids))
		for i, v := range ids {
			chatids[i] = telebot.ChatID(v)
		}
		s.DraftTransaction.To = chatids
		s.UserState = AwaitMessage
		return context.Send("Write your message:")
	}

	playerNames, chatIDs, err := db.GetNamesAndChatIDsOfAll(s.DB)
	if err != nil {
		context.Send("Error happened while processing your request, contact administrator")
		return err
	}
	s.UserState = AwaitResipient
	return context.Send("Names:", ui.PlayerNamesKeyboard(playerNames, chatIDs))
}

func handlePlayerName(context tele.Context, s *Session, cbData string) error {
	if s.UserState != AwaitResipient {
		return context.Send("This button is not available right now, please finish your previous action")
	}

	splited := strings.SplitAfterN(cbData, ":", 2)
	if len(splited) != 2 {
		return fmt.Errorf("splitting callbackdata, met unexpected amount of data: %d", len(splited))
	}

	toChatID, err := strconv.ParseInt(splited[1], 10, 64)
	if err != nil {
		return err
	}

	s.DraftTransaction.To = append(s.DraftTransaction.To, tele.ChatID(toChatID))
	s.UserState = AwaitMessage
	return context.Send("Write your message:")
}

func handleYesTitle(context tele.Context, s *Session) error {
	if s.UserState != AwaitTitleDecision {
		return context.Send("This button is not available right now, please finish your previous action")
	}
	s.UserState = AwaitTitle
	return context.Send("Write title for your message:")
}

func handleNoTitle(context tele.Context, s *Session) error {
	if s.UserState != AwaitTitleDecision {
		return context.Send("This button is not available right now, please finish your previous action")
	}
	return finilize(context, s)
}

func finilize(context tele.Context, s *Session) error {
	chatID := context.Chat().ID

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

		context.Bot().Send(toChatID, formatedMessage)

		log.Printf("Send message succesfully. From: %d to %d, transaction id: %d", userFrom.ChatID, toChatID, transaction.ID)
	}

	context.Send("Message sent")
	s.Done = true

	user, err := db.GetUserByID(s.DB, chatID)
	if err != nil {
		return err
	}
	return common.GetMainMenuByRole(context, s.DB, user)
}
