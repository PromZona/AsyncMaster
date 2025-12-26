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
	_, err := db.Exec("UPDATE users SET telegram_name = $1, player_name = $2, state = $3 WHERE chat_id = $4", user.TelegramName, user.PlayerName, user.State, user.ChatID)
	if err != nil {
		log.Print("ERROR: while updating user ", user.ChatID, ". ", err)
	}
}

func getUser(db *sql.DB, chatId int64) *UserData {
	var newUser UserData
	queryResult := db.QueryRow("select telegram_name, COALESCE(player_name, '') AS player_name, state, chat_id from users where chat_id = $1", chatId)
	err := queryResult.Scan(&newUser.TelegramName, &newUser.PlayerName, &newUser.State, &newUser.ChatID)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &newUser
}

func getUserByName(db *sql.DB, playerName string) (*UserData, error) {
	var user UserData
	queryResult := db.QueryRow("select telegram_name, COALESCE(player_name, '') AS player_name, state, chat_id from users where player_name = $1", playerName)
	err := queryResult.Scan(&user.TelegramName, &user.PlayerName, &user.State, &user.ChatID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &user, nil
}

func getUserPlayerNames(db *sql.DB) ([]string, error) {
	var result []string
	rows, err := db.Query("SELECT player_name FROM users")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		innerErr := rows.Scan(&name)
		if innerErr != nil {
			log.Print(err)
			return nil, err
		}
		result = append(result, name)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
		return nil, err
	}

	return result, nil
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

func createMessage(db *sql.DB, message *MessageStruct) {
	result, err := db.Exec("insert into messages (chat_id, message_id, message_title) values ($1, $2, $3)", message.ChatID, message.MessageID, message.MessageTitle)
	if err != nil {
		log.Print("Error whule creating Message: ", err)
		return
	}
	log.Print("Created new Message: ", result)
}

func getMessage(db *sql.DB, messageID string) (*MessageStruct, error) {
	var message MessageStruct
	queryResult := db.QueryRow("SELECT chat_id, message_title, message_id FROM messages WHERE message_id = $1", messageID)
	err := queryResult.Scan(&message.ChatID, &message.MessageTitle, &message.MessageID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &message, nil
}

func updateMessage() {

}

func deleteMessage() {

}
