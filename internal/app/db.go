package app

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func registerUser(db *sql.DB, user *UserData) error {
	log.Print("Register User: ", user.ChatID, " ", user.TelegramName)

	result, err := db.Exec("insert into users (chat_id, telegram_name, player_name) values ($1, $2, $3)", user.ChatID, user.TelegramName, user.PlayerName)
	if err != nil {
		log.Fatal("Failed to add user to db. ", err)
		return err
	}
	log.Print("DB: added new user ", result)

	return nil
}

func updateUser(db *sql.DB, user *UserData) {
	_, err := db.Exec("UPDATE users SET telegram_name = $1, player_name = $2 WHERE chat_id = $3", user.TelegramName, user.PlayerName, user.ChatID)
	if err != nil {
		log.Print("ERROR: while updating user ", user.ChatID, ". ", err)
	}
}

func getUser(db *sql.DB, chatId int64) *UserData {
	var newUser UserData
	queryResult := db.QueryRow("SELECT telegram_name, COALESCE(player_name, '') AS player_name, chat_id FROM users WHERE chat_id = $1", chatId)
	err := queryResult.Scan(&newUser.TelegramName, &newUser.PlayerName, &newUser.ChatID)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &newUser
}

func getUserByName(db *sql.DB, playerName string) (*UserData, error) {
	var user UserData
	queryResult := db.QueryRow("SELECT telegram_name, COALESCE(player_name, '') AS player_name, chat_id from users where player_name = $1", playerName)
	err := queryResult.Scan(&user.TelegramName, &user.PlayerName, &user.ChatID)
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

func getUserPlayerNamesAndChatID(db *sql.DB) (names []string, chatIDs []int64, err error) {
	var rows *sql.Rows
	rows, err = db.Query("SELECT player_name, chat_id FROM users")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var chatID int64
		err = rows.Scan(&name, &chatID)
		if err != nil {
			return nil, nil, err
		}
		names = append(names, name)
		chatIDs = append(chatIDs, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	return names, chatIDs, nil
}

func ensureUser(db *sql.DB, chatId int64) bool {
	var isExist bool
	queryResult := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE chat_id = $1)", chatId)
	queryResult.Scan(&isExist)

	return isExist
}

func createMessage(db *sql.DB, message *Message) {
	result, err := db.Exec("insert into messages (chat_id, message_id, message_title) values ($1, $2, $3)", message.ChatID, message.MessageID, message.MessageTitle)
	if err != nil {
		log.Print("Error whule creating Message: ", err)
		return
	}
	log.Print("Created new Message: ", result)
}

func getMessage(db *sql.DB, messageID string) (*Message, error) {
	var message Message
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
