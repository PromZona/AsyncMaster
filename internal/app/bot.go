package app

import (
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
)

type UserState int

const (
	// Normal State of Being
	UserStateDefault UserState = 0

	// Registration Phase
	UserStateAwaitCodename = 1

	// Master Commands
	UserStateSendingAll            = 2
	UserStateAwaitSavingMessage    = 3
	UserStateAwaitTitleForMesssage = 4

	// Player Commands

)

type UserData struct {
	ChatID       tele.ChatID
	TelegramName string
	PlayerName   string
	State        UserState
}

func (user *UserData) Recipient() string {
	return strconv.FormatInt(int64(user.ChatID), 10)
}

type BotData struct {
	DB           *sql.DB
	MessageCache map[int64]*MessageStruct
}

func BotInit(db *sql.DB) *BotData {
	bot := BotData{
		DB:           db,
		MessageCache: make(map[int64]*MessageStruct),
	}
	return &bot
}

type MessageStruct struct {
	ID           int
	MessageTitle string
	MessageID    string
	ChatID       int64
}

func (msg MessageStruct) MessageSig() (string, int64) {
	return msg.MessageID, msg.ChatID
}

func (msg MessageStruct) MessageHash() string {
	return msg.MessageID + string(msg.ChatID)
}
