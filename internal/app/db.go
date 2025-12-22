package app

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"
)

func registerUser(db *sql.DB, context tele.Context) *UserData {
	log.Print("Register User: ", context.Chat().ID, " ", context.Chat().FirstName)

	result, err := db.Exec("insert into users (chat_id, telegram_name) values ($1, $2)", context.Chat().ID, context.Chat().FirstName)
	if err != nil {
		log.Fatal("Failed to add user to db. ", err)
		return nil
	}
	log.Print("DB: added new user ", result)

	return getUser(db, context.Chat().ID)
}

func updateUser(db *sql.DB, user *UserData) {
	log.Print("updateUser: ID = ", user.ChatId)
	_, err := db.Exec("UPDATE users SET telegram_name = $1, player_name = $2, state = $3 WHERE chat_id = $4", user.TelegramName, user.PlayerName, user.State, user.ChatId)
	if err != nil {
		log.Print("ERROR: while updating user ", user.ChatId, ". ", err)
	}
}

func getUser(db *sql.DB, chatId int64) *UserData {
	log.Print("getUser: requested ", chatId)

	var newUser UserData
	queryResult := db.QueryRow("select telegram_name, COALESCE(player_name, '') AS player_name, state, chat_id from users where chat_id = $1", chatId)
	err := queryResult.Scan(&newUser.TelegramName, &newUser.PlayerName, &newUser.State, &newUser.ChatId)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &newUser
}

func ensureUser(db *sql.DB, chatId int64) bool {
	var isExist bool
	queryResult := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE chat_id = $1)", chatId)
	queryResult.Scan(&isExist)

	return isExist
}

func getOrCreateUser(db *sql.DB, context tele.Context) *UserData {
	chatId := context.Chat().ID
	user := getUser(db, chatId)
	if user == nil {
		user = registerUser(db, context)
	}
	return user
}
