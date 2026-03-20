package db

import (
	"database/sql"
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	"github.com/lib/pq"
)

type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}

func CreateUser(e DBExecutor, user *bot.UserData) ([]uint8, error) {
	var id []uint8

	err := e.QueryRow("insert into users (chat_id, telegram_name, player_name, role) values ($1, $2, $3, $4) RETURNING id",
		user.ChatID,
		user.TelegramName,
		user.PlayerName,
		user.Role).Scan(&id)
	if err != nil {
		return nil, err
	}
	return id, nil
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

func GetUsersAll(e DBExecutor) ([]bot.UserData, error) {
	var result []bot.UserData

	rows, err := e.Query(`
		SELECT
			u.telegram_name,
			COALESCE(player_name, '') AS player_name,
			u.chat_id,
			u.role,

			f.id,
			f.name,
			f.description,
			f.resources
		FROM users u
		JOIN factions f ON f.user_id = u.id
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user bot.UserData
		user.Faction = &bot.Faction{}

		innerErr := rows.Scan(
			&user.TelegramName,
			&user.PlayerName,
			&user.ChatID,
			&user.Role,
			&user.Faction.ID,
			&user.Faction.Name,
			&user.Faction.Description,
			&user.Faction.Resources)
		if innerErr != nil {
			log.Print(err)
			return nil, err
		}
		result = append(result, user)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
		return nil, err
	}

	return result, nil
}

func GetUserByID(e DBExecutor, chatID int64) (*bot.UserData, error) {
	var newUser bot.UserData
	newUser.Faction = &bot.Faction{}

	var factionID sql.NullInt32
	var factionName sql.NullString
	var factionDescription sql.NullString
	var factionResources sql.NullString

	queryResult := e.QueryRow(`
		SELECT
			u.telegram_name,
			COALESCE(player_name, '') AS player_name,
			u.chat_id,
			u.role,

			f.id,
			f.name,
			f.description,
			f.resources
		FROM users u
		LEFT JOIN factions f ON f.user_id = u.id
		WHERE chat_id = $1
		`, chatID)
	err := queryResult.Scan(
		&newUser.TelegramName,
		&newUser.PlayerName,
		&newUser.ChatID,
		&newUser.Role,
		&factionID,
		&factionName,
		&factionDescription,
		&factionResources)

	if factionID.Valid {
		newUser.Faction.ID = int(factionID.Int32)
		newUser.Faction.Name = factionName.String
		newUser.Faction.Description = factionDescription.String
		newUser.Faction.Resources = factionResources.String
	} else {
		newUser.Faction = nil
	}

	if err != nil {
		return nil, err
	}
	return &newUser, nil
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

func GetNamesAndChatIDsOfAll(e DBExecutor) (names []string, chatIDs []int64, err error) {
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

// Returns only players chat ids, without master
func GetNamesAndChatIDsOfPlayers(e DBExecutor) (names []string, chatIDs []int64, err error) {
	var rows *sql.Rows
	rows, err = e.Query("SELECT player_name, chat_id FROM users WHERE role = $1", bot.RolePlayer)
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
		transaction.From, pq.Array(transaction.To), transaction.Message.ID).
		Scan(&transaction.ID, &transaction.CreatedAt)

	if err != nil {
		log.Printf("SCAN ERROR: %#v\n", err)
		return nil, err
	}
	return transaction, nil
}

func GetLastMessageTransactions(e DBExecutor, toPlayerChatID int64) ([]*bot.MessageTransaction, error) {
	rows, err := e.Query(`
		SELECT
			mt.id,
			mt.created_at,
			mt.from_chat,
			mt.to_chat,

			m.id,
			m.message_title,
			m.message_id,
			m.chat_id,
			m.message_text
		FROM message_transaction mt
		JOIN messages m ON m.id = mt.message_id
		WHERE mt.to_chat = $1
		ORDER BY mt.created_at DESC
		LIMIT 10
	`, toPlayerChatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*bot.MessageTransaction

	for rows.Next() {
		mt := &bot.MessageTransaction{}
		msg := &bot.Message{}

		err := rows.Scan(
			&mt.ID,
			&mt.CreatedAt,
			&mt.From,
			&mt.To,

			&msg.ID,
			&msg.Title,
			&msg.MessageID,
			&msg.ChatID,
			&msg.Text,
		)
		if err != nil {
			return nil, err
		}

		mt.Message = msg
		result = append(result, mt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func GetMessageTransaction(e DBExecutor, transactionID int64) (*bot.MessageTransaction, error) {
	rows, err := e.Query(`
		SELECT
			mt.id,
			mt.created_at,
			mt.from_chat,
			mt.to_chat,

			m.id,
			m.message_title,
			m.message_id,
			m.chat_id,
			m.message_text
		FROM message_transaction mt
		JOIN messages m ON m.id = mt.message_id
		WHERE mt.id = $1
		ORDER BY mt.created_at DESC
		LIMIT 10
	`, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result *bot.MessageTransaction

	for rows.Next() {
		mt := &bot.MessageTransaction{}
		msg := &bot.Message{}

		err := rows.Scan(
			&mt.ID,
			&mt.CreatedAt,
			&mt.From,
			&mt.To,

			&msg.ID,
			&msg.Title,
			&msg.MessageID,
			&msg.ChatID,
			&msg.Text,
		)
		if err != nil {
			return nil, err
		}

		mt.Message = msg
		result = mt
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
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
	rows, err := e.Query(`
		SELECT
		 mr.id,
		 mr.to_player,
		 mr.text_request,
		 mr.created_at,
		 mr.updated_at,
		 mr.state,

		 r.id AS roll_id,
		 r.created_at AS roll_created_at,
		 r.title,
		 r.dice_count,
		 r.dice_sides,
		 r.roll_result
		FROM master_requests mr
		LEFT JOIN roll_requests r ON mr.id = r.transaction_id 
		WHERE id = $1`,
		id)
	if err != nil {
		return nil, err
	}

	var request bot.MasterRequest
	var rolls []*bot.RollRequest

	for rows.Next() {
		var id sql.NullInt32
		var createdAt sql.NullTime
		var title sql.NullString
		var diceCount sql.NullInt32
		var diceSides sql.NullInt32
		var rollResult sql.NullInt32

		err := rows.Scan(
			&request.ID,
			&request.To,
			&request.TextRequest,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.State,
			&id,
			&createdAt,
			&title,
			&diceCount,
			&diceSides,
			&rollResult)
		if err != nil {
			return nil, err
		}

		if id.Valid {
			rolls = append(rolls, &bot.RollRequest{
				ID:         int(id.Int32),
				CreatedAt:  createdAt.Time,
				Title:      title.String,
				DiceCount:  int(diceCount.Int32),
				DiceSides:  int(diceSides.Int32),
				RollResult: int(rollResult.Int32),
			})
		}
	}
	request.RollRequests = rolls
	return &request, nil
}

func GetFirstUnansweredMasterRequest(e DBExecutor, chatID int64) (*bot.MasterRequest, error) {
	rows, err := e.Query(`
		SELECT
		 mr.id,
		 mr.to_player,
		 mr.text_request,
		 mr.created_at,
		 mr.updated_at,
		 mr.state,

		 r.id AS roll_id,
		 r.created_at AS roll_created_at,
		 r.title,
		 r.dice_count,
		 r.dice_sides,
		 r.roll_result
		FROM master_requests mr
		LEFT JOIN roll_requests r ON mr.id = r.transaction_id 
		WHERE to_player = $1 AND state = $2
		ORDER BY created_at ASC`,
		chatID, bot.MRUnasnwered)
	if err != nil {
		return nil, err
	}

	var request bot.MasterRequest
	var rolls []*bot.RollRequest

	for rows.Next() {
		var id sql.NullInt32
		var createdAt sql.NullTime
		var title sql.NullString
		var diceCount sql.NullInt32
		var diceSides sql.NullInt32
		var rollResult sql.NullInt32

		err := rows.Scan(
			&request.ID,
			&request.To,
			&request.TextRequest,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.State,
			&id,
			&createdAt,
			&title,
			&diceCount,
			&diceSides,
			&rollResult)
		if err != nil {
			return nil, err
		}

		if id.Valid {
			rolls = append(rolls, &bot.RollRequest{
				ID:         int(id.Int32),
				CreatedAt:  createdAt.Time,
				Title:      title.String,
				DiceCount:  int(diceCount.Int32),
				DiceSides:  int(diceSides.Int32),
				RollResult: int(rollResult.Int32),
			})
		}
	}
	request.RollRequests = rolls
	return &request, nil
}

func GetFirstAnsweredMasterRequest(e DBExecutor) (*bot.MasterRequest, error) {

	rows, err := e.Query(`
		SELECT
		 mr.id,
		 mr.to_player,
		 mr.text_request,
		 mr.created_at,
		 mr.updated_at,
		 mr.state,

		 r.id AS roll_id,
		 r.created_at AS roll_created_at,
		 r.title,
		 r.dice_count,
		 r.dice_sides,
		 r.roll_result
		FROM
		(
    		SELECT *
    		FROM master_requests
    		WHERE state = $1
    		ORDER BY updated_at ASC
    		LIMIT 1
		) mr
		LEFT JOIN roll_requests r ON mr.id = r.transaction_id`,
		bot.MRAnswered)
	if err != nil {
		return nil, err
	}

	var request bot.MasterRequest
	var rolls []*bot.RollRequest

	for rows.Next() {
		var id sql.NullInt32
		var createdAt sql.NullTime
		var title sql.NullString
		var diceCount sql.NullInt32
		var diceSides sql.NullInt32
		var rollResult sql.NullInt32

		err := rows.Scan(
			&request.ID,
			&request.To,
			&request.TextRequest,
			&request.CreatedAt,
			&request.UpdatedAt,
			&request.State,
			&id,
			&createdAt,
			&title,
			&diceCount,
			&diceSides,
			&rollResult)
		if err != nil {
			return nil, err
		}

		if id.Valid {
			rolls = append(rolls, &bot.RollRequest{
				ID:         int(id.Int32),
				CreatedAt:  createdAt.Time,
				Title:      title.String,
				DiceCount:  int(diceCount.Int32),
				DiceSides:  int(diceSides.Int32),
				RollResult: int(rollResult.Int32),
			})
		}
	}
	request.RollRequests = rolls
	return &request, nil
}

func UpdateMasterRequest(e DBExecutor, masterRequest *bot.MasterRequest) error {

	/*
		WHAT DOES NOT WORK:
		- Player answering with text for master request.
		WHAT TO DO:
		- Delete all testing db data. Let's start clean.
		- Track down if there are some missing logic in handling this case.
	*/

	_, err := e.Exec(`
		UPDATE
		 master_requests
		SET
			to_player = $1,
			text_request = $2,
			text_response = $3,
			state = $4
		WHERE
			id = $5`,
		masterRequest.To,
		masterRequest.TextRequest,
		masterRequest.TextResponse,
		masterRequest.State,
		masterRequest.ID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMasterRequestState(e DBExecutor, masterRequestID int64, state bot.MasterRequestState) error {
	_, err := e.Exec("UPDATE master_requests SET state = $1 WHERE id = $2",
		state,
		masterRequestID)
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

func CreateFaction(e DBExecutor, faction *bot.Faction, userUUID []uint8) (*bot.Faction, error) {
	err := e.QueryRow("INSERT INTO factions (name, description, resources, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
		faction.Name, faction.Description, faction.Resources, userUUID).
		Scan(&faction.ID)

	if err != nil {
		return nil, err
	}
	return faction, nil
}
func GetFaction(e DBExecutor) {

}

func UpdateFaction(e DBExecutor) {

}
func DeleteFaction(e DBExecutor) {

}

func GetPlayerMenuData(e DBExecutor, chatID int64) *ui.PlayerMenu {
	result := &ui.PlayerMenu{}
	err := e.QueryRow(`
		SELECT
			u.player_name,
			f.name,
			f.description,
			f.resources,
			(
        		SELECT COUNT(*)
        		FROM master_requests mr
        		WHERE mr.to_player = $1 AND mr.state = $2
    		) AS unanswered_master_requests_count
		FROM users u
		JOIN factions f
			ON f.user_id = u.id
		`, chatID, bot.MRUnasnwered).Scan(
		&result.PlayerName,
		&result.FactionName,
		&result.FactionDescription,
		&result.FactionResources,
		&result.UnansweredMasterRequests,
	)
	if err != nil {
		log.Printf("Failed to get player menu data from DB. Returning empty menu to player: %d", chatID)
		log.Printf("%s", err)
		return result
	}

	return result
}

func GetMasterMenuData(e DBExecutor, chatID int64) *ui.MasterMenu {
	return &ui.MasterMenu{}
}
