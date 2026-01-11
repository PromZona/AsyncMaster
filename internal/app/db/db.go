package db

import (
	"database/sql"
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	_ "github.com/lib/pq"
)

type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}

func CreateUser(e DBExecutor, user *bot.UserData) error {
	log.Print("Register User: ", user.ChatID, " ", user.TelegramName)

	result, err := e.Exec("insert into users (chat_id, telegram_name, player_name, role) values ($1, $2, $3, $4)",
		user.ChatID,
		user.TelegramName,
		user.PlayerName,
		user.Role)
	if err != nil {
		log.Fatal("Failed to add user to database ", err)
		return err
	}
	log.Print("DB: added new user ", result)

	return nil
}

func UpdateUser(e DBExecutor, user *bot.UserData) {
	_, err := e.Exec("UPDATE users SET telegram_name = $1, player_name = $2, role = $3 WHERE chat_id = $4",
		user.TelegramName,
		user.PlayerName,
		user.Role,
		user.ChatID)
	if err != nil {
		log.Print("ERROR: while updating user ", user.ChatID, ". ", err)
	}
}

func GetUserByID(e DBExecutor, chatID int64) *bot.UserData {
	var newUser bot.UserData
	queryResult := e.QueryRow("SELECT telegram_name, COALESCE(player_name, '') AS player_name, chat_id, role FROM users WHERE chat_id = $1", chatID)
	err := queryResult.Scan(&newUser.TelegramName, &newUser.PlayerName, &newUser.ChatID, &newUser.Role)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &newUser
}

func GetUserByName(e DBExecutor, playerName string) (*bot.UserData, error) {
	var user bot.UserData
	queryResult := e.QueryRow("SELECT telegram_name, COALESCE(player_name, '') AS player_name, chat_id from users where player_name = $1", playerName)
	err := queryResult.Scan(&user.TelegramName, &user.PlayerName, &user.ChatID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &user, nil
}

func GetUserPlayerNames(e DBExecutor) ([]string, error) {
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

func GetUserPlayerNamesAndChatID(e DBExecutor) (names []string, chatIDs []int64, err error) {
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

func EnsureUserExist(e DBExecutor, chatID int64) bool {
	var isExist bool
	queryResult := e.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE chat_id = $1)", chatID)
	queryResult.Scan(&isExist)

	return isExist
}

func CreateMessage(e DBExecutor, message *bot.Message) (*bot.Message, error) {
	err := e.QueryRow("INSERT INTO messages (chat_id, message_id, message_title, message_text) values ($1, $2, $3, $4) RETURNING id",
		message.ChatID, message.MessageID, message.Title, message.Text).
		Scan(&message.ID)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func GetMessageByID(e DBExecutor, messageID string) (*bot.Message, error) {
	var message bot.Message
	queryResult := e.QueryRow("SELECT chat_id, message_title, message_id FROM messages WHERE message_id = $1", messageID)
	err := queryResult.Scan(&message.ChatID, &message.Title, &message.MessageID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &message, nil
}

func UpdateMessage() {

}

func DeleteMessage() {

}

func CreateMesssageTransaction(e DBExecutor, transaction *bot.MessageTransaction) (*bot.MessageTransaction, error) {
	err := e.QueryRow("INSERT INTO message_transaction (from_chat, to_chat, message_id) VALUES ($1, $2, $3) RETURNING id, created_at",
		transaction.From, transaction.To, transaction.Message.ID).
		Scan(&transaction.ID, &transaction.CreatedAt)

	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func GetMessageTransaction() {

}

func UpdateMesssageTransaction() {

}

func DeleteMessageTransaction() {

}

func CreateMasterRequest(e DBExecutor, request *bot.MasterRequest) (*bot.MasterRequest, error) {
	err := e.QueryRow("INSERT INTO master_requests (text_request, to_player) VALUES ($1, $2) RETURNING id, created_at",
		request.TextRequest, request.To).
		Scan(&request.ID, &request.CreatedAt)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func GetMasterRequestByID(e DBExecutor, id int) (*bot.MasterRequest, error) {
	var request bot.MasterRequest
	queryResult := e.QueryRow("SELECT id, to_player, text_request, created_at, is_answered FROM master_requests WHERE id = $1", id)
	err := queryResult.Scan(&request.ID, &request.To, &request.TextRequest, &request.CreatedAt, &request.IsAnswered)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &request, nil
}

func UpdateMasterRequest(e DBExecutor, masterRequest *bot.MasterRequest) error {
	_, err := e.Exec("UPDATE master_requests SET to_player = $1, text_request = $2, text_response = $3, is_answered = $4 WHERE id = $5",
		masterRequest.To,
		masterRequest.TextRequest,
		masterRequest.TextResponse,
		masterRequest.IsAnswered,
		masterRequest.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMasterRequest(e DBExecutor) error {
	return nil
}

func CreateRollRequest(e DBExecutor, roll *bot.RollRequest, masterRequestID int) (*bot.RollRequest, error) {
	err := e.QueryRow("INSERT INTO roll_requests (title, dice_count, dice_sides, transaction_id) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		roll.Title, roll.DiceCount, roll.DiceSides, masterRequestID).
		Scan(&roll.ID, &roll.CreatedAt)

	if err != nil {
		return nil, err
	}
	return roll, nil
}

func GetRollRequestByID(e DBExecutor, rollID int) (*bot.RollRequest, error) {
	var request bot.RollRequest
	queryResult := e.QueryRow("SELECT id, created_at, title, dice_count, dice_sides FROM roll_requests WHERE id = $1", rollID)
	err := queryResult.Scan(&request.ID, &request.CreatedAt, &request.Title, &request.DiceCount, &request.DiceSides)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &request, nil
}

func UpdateRollRequest(e DBExecutor, roll *bot.RollRequest) error {
	_, err := e.Exec("UPDATE roll_requests SET title = $1, dice_count = $2, dice_sides = $3, roll_result = $4 WHERE id = $5",
		roll.Title,
		roll.DiceCount,
		roll.DiceSides,
		roll.RollResult,
		roll.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRollRequest(e DBExecutor) error {
	return nil
}
