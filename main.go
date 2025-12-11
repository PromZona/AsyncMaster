package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
)

type UserState int

const (
	// Registration Phase
	UserStateFirstTimePassword UserState = iota
	UserStateCodename

	// Normal State of Being
	UserStateDefault

	// Master Commands
	UserStateSendingAll

	// Player Commands

)

type UserData struct {
	chatId       tele.ChatID
	telegramName string
	playerName   string
	state        UserState
}

type BotData struct {
	db    *sql.DB
	users map[int64]*UserData
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	psqlInfo := os.Getenv("DB_CONNECTION_STRING")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Print("Database successfully connected!")

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	botData := BotData{db: db, users: make(map[int64]*UserData)}

	b.Handle(tele.OnText, botData.handleText)
	b.Handle("/sendAll", botData.masterSendMessageToAll)

	b.Start()
}

func (bot *BotData) masterSendMessageToAll(context tele.Context) error {

	id := context.Chat().ID
	isUserExist := bot.ensureUser(id)
	var user *UserData
	if !isUserExist {
		log.Print("Not found user: ", id)
		user = bot.registerUser(context)
	} else {
		user = bot.getUser(id)
	}
	user.state = UserStateSendingAll
	bot.updateUser(bot.db, user)
	return context.Send("What to send?")
}

func (bot *BotData) handleText(context tele.Context) error {

	return nil

	id := context.Chat().ID
	user := bot.users[id]
	switch user.state {
	case UserStateDefault:
		return context.Send("What do you want from me?")
	case UserStateSendingAll:
		message := context.Message().Text
		log.Print(string(message))
		return context.Send("Thanks")
	}
	return nil
}

func (bot *BotData) registerUser(context tele.Context) *UserData {
	log.Print("Register User: ", context.Chat().ID, " ", context.Chat().FirstName)

	result, err := bot.db.Exec("insert into users (chat_id, telegram_name) values ($1, $2)", context.Chat().ID, context.Chat().FirstName)
	if err != nil {
		log.Fatal("Failed to add user to db. ", err)
		return nil
	}
	log.Print("DB: added new user ", result)

	return bot.getUser(context.Chat().ID)
}

func (bot *BotData) updateUser(db *sql.DB, user *UserData) {
	log.Print("updateUser: ID = ", user.chatId)
	_, err := db.Exec("UPDATE users SET telegram_name = $1, player_name = $2, state = $3 WHERE chat_id = $4", user.telegramName, user.playerName, user.state, user.chatId)
	if err != nil {
		log.Print("ERROR: while updating user ", user.chatId, ". ", err)
	}
}

func (bot *BotData) ensureUser(chatId int64) bool {

	if bot.users[chatId] != nil {
		return true
	}

	var isExist bool
	queryResult := bot.db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE chat_id = $1)", chatId)
	queryResult.Scan(&isExist)

	return isExist
}

func (bot *BotData) getUser(chatId int64) *UserData {
	log.Print("getUser: requested ", chatId)

	var newUser UserData
	queryResult := bot.db.QueryRow("select telegram_name, COALESCE(player_name, '') AS player_name, state, chat_id from users where chat_id = $1", chatId)
	err := queryResult.Scan(&newUser.telegramName, &newUser.playerName, &newUser.state, &newUser.chatId)
	if err != nil {
		log.Print("ERROR: getUser ", chatId, ". ", err)
		return nil
	}
	return &newUser
}
