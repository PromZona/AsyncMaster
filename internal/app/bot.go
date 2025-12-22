package app

import (
	"database/sql"

	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
)

type UserState int

const (
	// Registration Phase
	UserStateCodename UserState = iota

	// Normal State of Being
	UserStateDefault

	// Master Commands
	UserStateSendingAll

	// Player Commands

)

type UserData struct {
	ChatId       tele.ChatID
	TelegramName string
	PlayerName   string
	State        UserState
}

type BotData struct {
	DB    *sql.DB
	Users map[int64]*UserData
}
