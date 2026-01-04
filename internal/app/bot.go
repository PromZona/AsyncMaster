package app

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
)

type UserState int

const (
	// Normal State of Being
	UserStateDefault UserState = 0

	// Registration Phase
	UserStateAwaitPassword = 100
	UserStateAwaitCodename = 101

	// Master Commands
	UserStateAwaitSavingMessage    = 200
	UserStateAwaitTitleForMesssage = 201

	// Player Commands
	UserStateAwaitResipient     = 300
	UserStateAwaitMessage       = 301
	UserStateAwaitTitleDecision = 302
	UserStateAwaitTitle         = 303
)

type UserData struct {
	ChatID       int64
	TelegramName string
	PlayerName   string
}

func (user *UserData) Recipient() string {
	return strconv.FormatInt(int64(user.ChatID), 10)
}

type BotData struct {
	DB                      *sql.DB
	UserSessionState        map[int64]UserState
	MessageCache            map[int64]*Message
	MessageTransactionCache map[int64]*MessageTransaction
	UserRegistrationCache   map[int64]*UserData

	//Player Menu Data
	PlayerMenu    *tele.ReplyMarkup
	BtnPlayerSend tele.Btn
}

func BotInit(db *sql.DB) *BotData {
	bot := &BotData{
		DB:                      db,
		UserSessionState:        make(map[int64]UserState),
		MessageCache:            make(map[int64]*Message),
		MessageTransactionCache: make(map[int64]*MessageTransaction),
		UserRegistrationCache:   make(map[int64]*UserData),
		PlayerMenu:              &tele.ReplyMarkup{},
	}

	bot.BtnPlayerSend = bot.PlayerMenu.Data("Send To Player", "send")

	bot.PlayerMenu.Inline(
		bot.PlayerMenu.Row(bot.BtnPlayerSend),
	)

	return bot
}

func (b *BotData) ClearUserCache(chatID int64) {
	delete(b.MessageCache, chatID)
	delete(b.MessageTransactionCache, chatID)
	delete(b.UserRegistrationCache, chatID)
}

type Message struct {
	ID           int
	MessageTitle string
	MessageID    string
	ChatID       int64 // from which chat to copy
}

func (msg Message) MessageSig() (string, int64) {
	return msg.MessageID, msg.ChatID
}

func (msg Message) MessageHash() string {
	return msg.MessageID + strconv.FormatInt(msg.ChatID, 10)
}

type MessageTransaction struct {
	ID        int
	CreatedAt time.Time
	From      tele.ChatID
	To        tele.ChatID

	Message *Message
}

func parseCallbackDataString(callbackData string) (unique, data string) {
	trimmed := strings.Trim(callbackData, "\f")
	parts := strings.SplitN(trimmed, "|", 2)
	count := len(parts)
	if count == 2 {
		return parts[0], parts[1]
	}
	if count == 1 {
		return parts[0], ""
	}
	return "", ""
}
