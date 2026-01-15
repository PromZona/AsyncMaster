package listmessages

import (
	"fmt"
	"strconv"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

func handleStartFlow(context tele.Context, s *Session) error {
	if s.UserState != FlowStart {
		return context.Send("This action is not available right now, finish previous action first")
	}

	messages, err := db.GetLastMessageTransactions(s.DB, context.Chat().ID)
	if err != nil {
		return err
	}

	s.UserState = AwaitMessagePick
	return context.Send("Your last 10 messages, pick one", ui.UserMessagesKeyboard(messages))
}

func handleMessagePick(context tele.Context, s *Session, cbData string) error {
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

	messageFromPlayerName := db.GetUserByID(s.DB, int64(transaction.From)).PlayerName

	player := db.GetUserByID(s.DB, int64(transaction.To))
	messageToPlayerName := player.PlayerName

	formatedMessage := fmt.Sprintf("Title: %s\n\nFrom: %s\nTo: %s\n\n %s",
		transaction.Message.Title,
		messageFromPlayerName,
		messageToPlayerName,
		transaction.Message.Text)

	context.Send(formatedMessage)
	return finilize(context, s, player)
}

func finilize(context tele.Context, s *Session, player *bot.UserData) error {
	s.Done = true
	return ui.MainMenuKeyboard(context, player.Role)
}
