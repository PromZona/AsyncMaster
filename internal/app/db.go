package app

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}

func registerUser(e DBExecutor, user *UserData) error {
	log.Print("Register User: ", user.ChatID, " ", user.TelegramName)

	result, err := e.Exec("insert into users (chat_id, telegram_name, player_name) values ($1, $2, $3)", user.ChatID, user.TelegramName, user.PlayerName)
	if err != nil {
		log.Fatal("Failed to add user to database ", err)
		return err
	}
	log.Print("DB: added new user ", result)

	return nil
}

func updateUser(e DBExecutor, user *UserData) {
	_, err := e.Exec("UPDATE users SET telegram_name = $1, player_name = $2 WHERE chat_id = $3", user.TelegramName, user.PlayerName, user.ChatID)
	if err != nil {
		log.Print("ERROR: while updating user ", user.ChatID, ". ", err)
	}
}

func getUser(e DBExecutor, chatId int64) *UserData {
	var newUser UserData
	queryResult := e.QueryRow("SELECT telegram_name, COALESCE(player_name, '') AS player_name, chat_id FROM users WHERE chat_id = $1", chatId)
	err := queryResult.Scan(&newUser.TelegramName, &newUser.PlayerName, &newUser.ChatID)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &newUser
}

func getUserByName(e DBExecutor, playerName string) (*UserData, error) {
	var user UserData
	queryResult := e.QueryRow("SELECT telegram_name, COALESCE(player_name, '') AS player_name, chat_id from users where player_name = $1", playerName)
	err := queryResult.Scan(&user.TelegramName, &user.PlayerName, &user.ChatID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &user, nil
}

func getUserPlayerNames(e DBExecutor) ([]string, error) {
	var result []string
	rows, err := e.Query("SELECT player_name FROM users")
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

func getUserPlayerNamesAndChatID(e DBExecutor) (names []string, chatIDs []int64, err error) {
	var rows *sql.Rows
	rows, err = e.Query("SELECT player_name, chat_id FROM users")
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

func ensureUser(e DBExecutor, chatId int64) bool {
	var isExist bool
	queryResult := e.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE chat_id = $1)", chatId)
	queryResult.Scan(&isExist)

	return isExist
}

func createMessage(e DBExecutor, message *Message) (*Message, error) {
	err := e.QueryRow("INSERT INTO messages (chat_id, message_id, message_title) values ($1, $2, $3) RETURNING id",
		message.ChatID, message.MessageID, message.MessageTitle).
		Scan(&message.ID)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func getMessage(e DBExecutor, messageID string) (*Message, error) {
	var message Message
	queryResult := e.QueryRow("SELECT chat_id, message_title, message_id FROM messages WHERE message_id = $1", messageID)
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

func createTransaction(e DBExecutor, transaction *MessageTransaction) (*MessageTransaction, error) {
	err := e.QueryRow("INSERT INTO message_transaction (from_chat, to_chat, message_id) VALUES ($1, $2, $3) RETURNING id, created_at",
		transaction.From, transaction.To, transaction.Message.ID).
		Scan(&transaction.ID, &transaction.CreatedAt)

	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func getTransaction() {

}

func updateTransaction() {

}

func deleteTransaction() {

}
