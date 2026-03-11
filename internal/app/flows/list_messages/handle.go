package listmessages

import (
	"fmt"
	"strconv"

	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
)

func handleStartFlow(context runtime.Context, s *Session) error {
	if s.UserState != FlowStart {
		return context.Send("This action is not available right now, finish previous action first")
	}

	messages, err := db.GetLastMessageTransactions(s.DB, context.ChatID())
	if err != nil {
		return err
	}

	s.UserState = AwaitMessagePick
	return context.Send("Your last 10 messages, pick one", ui.UserMessagesKeyboard(messages))
}

func handleMessagePick(context runtime.Context, s *Session, cbData string) error {
	if s.UserState != AwaitMessagePick {
		return context.Send("This action is not available right now, finish previous action first")
	}

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

	s.Done = true
	return context.Send(formatedMessage)
}
