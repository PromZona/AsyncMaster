package sendmessage

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

func DispatchText(context tele.Context, b *bot.BotData) error {
	chatID := context.Chat().ID
	state := b.UserSessionState[chatID]

	switch state {
	case bot.UserStateAwaitMessage:
		return handleMessageText(context, b)
	case bot.UserStateAwaitTitle:
		return handleMessageTitle(context, b)
	default:
		return fmt.Errorf("sendmessage met unsupported user state while dispatching text: %d", state)
	}
}

func DispatchCallback(context tele.Context, b *bot.BotData, cbUnique string, cbData string) error {

	switch cbUnique {
	case "send":
		return handleInitialSend(context, b)
	case "player_names":
		return handlePlayerName(context, b, cbData)
	case "no_title":
		return handleNoTitle(context, b)
	case "yes_title":
		return handleYesTitle(context, b)
	default:
		return fmt.Errorf("sendmessage met unsupported <callback unique> while dispatching callback: %s", cbUnique)
	}
}

func handleMessageText(context tele.Context, b *bot.BotData) error {
	chatID := context.Chat().ID

	message := &bot.Message{
		Title:     "",
		MessageID: strconv.FormatInt(int64(context.Message().ID), 10),
		ChatID:    chatID,
		Text:      context.Text(),
	}

	b.MessageCache[chatID] = message
	b.UserSessionState[chatID] = bot.UserStateAwaitTitleDecision

	return context.Send("Do you want to add title for a message?", ui.TitleQuestion(b))
}

func handleMessageTitle(context tele.Context, b *bot.BotData) error {
	chatID := context.Chat().ID
	b.MessageCache[chatID].Title = context.Message().Text
	return finilize(context, b)
}

func handleInitialSend(context tele.Context, b *bot.BotData) error {
	chatID := context.Chat().ID
	state := b.UserSessionState[chatID]

	if state != bot.UserStateDefault {
		return context.Send("This button is not available right now, please finish your previous action")
	}
	b.UserSessionState[chatID] = bot.UserStateAwaitResipient
	b.MessageTransactionCache[chatID] = &bot.MessageTransaction{
		From: tele.ChatID(chatID),
	}
	return context.Send("Names:", ui.PlayerNames(b))
}

func handlePlayerName(context tele.Context, b *bot.BotData, cbData string) error {
	chatID := context.Chat().ID
	state := b.UserSessionState[chatID]

	if state != bot.UserStateAwaitResipient {
		return context.Send("This button is not available right now, please finish your previous action")
	}

	splited := strings.SplitAfterN(cbData, ":", 2)
	if len(splited) != 2 {
		return fmt.Errorf("splitting callbackdata, met unexpected amount of data: %d", len(splited))
	}

	transaction, ok := b.MessageTransactionCache[chatID]
	if !ok {
		return fmt.Errorf("retrieving transaction message. No message entry found")
	}

	toChatID, err := strconv.ParseInt(splited[1], 10, 64)
	if err != nil {
		return err
	}

	transaction.To = tele.ChatID(toChatID)
	b.UserSessionState[chatID] = bot.UserStateAwaitMessage
	return context.Send("Write your message:")
}

func handleYesTitle(context tele.Context, b *bot.BotData) error {
	chatID := context.Chat().ID
	state := b.UserSessionState[chatID]

	if state != bot.UserStateAwaitTitleDecision {
		return context.Send("This button is not available right now, please finish your previous action")
	}
	b.UserSessionState[chatID] = bot.UserStateAwaitTitle
	return context.Send("Write title for your message:")
}

func handleNoTitle(context tele.Context, b *bot.BotData) error {

	chatID := context.Chat().ID
	state := b.UserSessionState[chatID]
	if state != bot.UserStateAwaitTitleDecision {
		return context.Send("This button is not available right now, please finish your previous action")
	}
	return finilize(context, b)
}

func finilize(context tele.Context, b *bot.BotData) error {
	chatID := context.Chat().ID

	transaction, ok := b.MessageTransactionCache[chatID]
	if !ok {
		return fmt.Errorf("expected transaction message to exist while handling UserStateAwaitMessage state")
	}

	message, ok := b.MessageCache[chatID]
	if !ok {
		return fmt.Errorf("expected message to exist while handling UserStateAwaitMessage state")
	}

	tx, err := b.DB.Begin()
	if err != nil {
		return err
	}

	message, err = db.CreateMessage(tx, message)
	if err != nil {
		tx.Rollback()
		return err
	}

	transaction.Message = message
	transaction, err = db.CreateTransaction(tx, transaction)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	messageFromPlayerName := db.GetUser(b.DB, int64(transaction.From)).PlayerName
	messageToPlayerName := db.GetUser(b.DB, int64(transaction.To)).PlayerName

	formatedMessage := fmt.Sprintf("Title: %s\n\nFrom: %s\nTo: %s\n\n %s",
		message.Title,
		messageFromPlayerName,
		messageToPlayerName,
		message.Text)

	context.Bot().Send(transaction.To, formatedMessage)

	log.Printf("Player send message succesfully, transaction id: %d", transaction.ID)

	b.UserSessionState[chatID] = bot.UserStateDefault
	b.ClearUserCache(chatID)

	context.Send("Message sent")
	return ui.MainMenu(context, b)
}
