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
	UserStateAwaitPassword = 100
	UserStateAwaitCodename = 101

	// Master Commands
	UserStateAwaitSavingMessage    = 200
	UserStateAwaitTitleForMesssage = 201

	// Player Commands

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
	DB                    *sql.DB
	UserSessionState      map[int64]UserState
	MessageCache          map[int64]*MessageStruct
	UserRegistrationCache map[int64]*UserData

	//Player Menu Data
	PlayerMenu    *tele.ReplyMarkup
	BtnPlayerSend tele.Btn
}

func BotInit(db *sql.DB) *BotData {
	bot := &BotData{
		DB:                    db,
		UserSessionState:      make(map[int64]UserState),
		MessageCache:          make(map[int64]*MessageStruct),
		UserRegistrationCache: make(map[int64]*UserData),
		PlayerMenu:            &tele.ReplyMarkup{},
	}

	bot.BtnPlayerSend = bot.PlayerMenu.Data("Send To Player", "send")

	bot.PlayerMenu.Inline(
		bot.PlayerMenu.Row(bot.BtnPlayerSend),
	)

	return bot
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
