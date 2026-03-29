package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
)

// CREATE MASTER REQUEST SCENARIO

func InitialMasterRequest(context runtime.Context, s *bot.Session) error {
	playerNames, chatIDs, err := db.GetNamesAndChatIDsOfPlayers(s.DB)
	if err != nil {
		return err
	}

	if s.IsSendEveryone {
		s.Resipients = append(s.Resipients, chatIDs...)
		return context.Send("Type text which will be sent to players")
	}

	return context.Send("Pick a player to send to", ui.PlayerNamesKeyboard(bot.HID_mr_resipient, playerNames, chatIDs))
}

func InitialMasterRequestEveryone(context runtime.Context, s *bot.Session) error {
	s.IsSendEveryone = true
	return InitialMasterRequest(context, s)
}

func MasterRequestResipient(context runtime.Context, s *bot.Session) error {
	_, cbData := ParseCallbackDataString(context.Callback())
	splited := strings.SplitAfterN(cbData, ":", 2)
	if len(splited) != 2 {
		return fmt.Errorf("splitting callbackdata, met unexpected amount of data: %d", len(splited))
	}

	toChatID, err := strconv.ParseInt(splited[1], 10, 64)
	if err != nil {
		return err
	}

	s.Resipients = append(s.Resipients, toChatID)
	return context.Send("Type text which will be sent to player")
}

func MasterRequestText(context runtime.Context, s *bot.Session) error {
	s.MasterRequest = &bot.MasterRequest{
		To:           0,
		TextRequest:  "",
		TextResponse: "",
		State:        bot.MRUnasnwered,
		RollRequests: nil,
	}
	s.MasterRequest.TextRequest = context.MessageText()
	return context.Send("Do you want to add dice request?", ui.YesNoKeyboard(bot.HID_mr_dice_yes, bot.HID_mr_dice_no))
}

func MasterRequestDiceYes(context runtime.Context, s *bot.Session) error {
	return context.Send("Write a roll in a format: <count>d<dice> <name of the role>\nExample: 1d6 Roll on money")
}

func MasterRequestDiceNo(context runtime.Context, s *bot.Session) error {
	return masterRequestCreationFinilize(context, s)
}

func MasterRequestAddRoll(context runtime.Context, s *bot.Session) error {
	args := strings.SplitAfterN(context.MessageText(), " ", 2)
	if len(args) != 2 {
		return context.Send(fmt.Sprintf("Expected 2 arguments, but received: %d", len(args)))
	}
	rollString := args[0]
	title := args[1]

	var diceCount, diceSides int
	_, err := fmt.Sscanf(rollString, "%dd%d", &diceCount, &diceSides)
	if err != nil {
		return err
	}

	s.RollRequests = append(s.RollRequests, &bot.RollRequest{
		Title:     title,
		DiceCount: diceCount,
		DiceSides: diceSides,
	})
	return context.Send("Do you want to add dice request?", ui.YesNoKeyboard(bot.HID_mr_dice_yes, bot.HID_mr_dice_no))
}

func masterRequestCreationFinilize(context runtime.Context, s *bot.Session) error {
	masterRequest := s.MasterRequest
	if masterRequest == nil {
		return fmt.Errorf("master request is nil while submiting data to database")
	}

	for _, res := range s.Resipients {

		masterRequest.To = res

		tx, err := s.DB.Begin()
		if err != nil {
			return err
		}

		masterRequest, err = db.CreateMasterRequest(tx, masterRequest)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, value := range s.RollRequests {
			_, err = db.CreateRollRequest(tx, value, masterRequest.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		masterRequest.RollRequests = s.RollRequests

		err = tx.Commit()
		if err != nil {
			return err
		}

		formattedMessage := fmt.Sprintf("Мастер требует внимания\n\n%s", masterRequest.TextRequest)
		context.SendTo(masterRequest.To, formattedMessage, ui.AnswerMasterKeyboard(masterRequest))

		log.Printf("Send Master Request to %d. master request id: %d", res, masterRequest.ID)
	}

	return context.Send("Message send to resipient")
}

// END CREATE MASTER REQUEST SCENARIO

//-------------------------------------

// MASTER REQUEST ANSWER SCENARIO

func MasterRequestAnswer(context runtime.Context, s *bot.Session) error {
	_, cbData := ParseCallbackDataString(context.Callback())

	var masterRequestID int
	_, err := fmt.Sscanf(cbData, "%d", &masterRequestID)
	if err != nil {
		return err
	}

	masterRequest, err := db.GetMasterRequestByID(s.DB, masterRequestID)
	if err != nil {
		return err
	}

	s.MasterRequest = masterRequest
	return context.Send("Напиши ответ Мастеру")
}

func MasterRequestAnswerText(context runtime.Context, s *bot.Session) error {
	if s.MasterRequest == nil {
		return fmt.Errorf("master request is nil while trying to write text into it")
	}

	text := context.MessageText()
	s.MasterRequest.TextResponse = text
	s.MasterRequest.State = bot.MRAnswered

	err := db.UpdateMasterRequest(s.DB, s.MasterRequest)
	if err != nil {
		return err
	}

	log.Printf("Player %d replied to master request %d", context.ChatID(), s.MasterRequest.ID)
	return context.Send("Ответ был принят\nНе забудь бросить кости, если Мастер спросил")
}

func MasterRequestAnswerRoll(context runtime.Context, s *bot.Session) error {
	_, cbData := ParseCallbackDataString(context.Callback())

	var rollID int
	_, err := fmt.Sscanf(cbData, "%d", &rollID)
	if err != nil {
		return err
	}

	roll, err := db.GetRollRequestByID(s.DB, rollID)
	if err != nil {
		return err
	}

	rollResult := 0
	for i := 0; i < roll.DiceCount; i++ {
		rollResult += rand.Intn(roll.DiceSides) + 1
	}

	roll.RollResult = rollResult
	db.UpdateRollRequest(s.DB, roll)

	log.Printf("Player %d made a roll %d", context.ChatID(), roll.ID)
	textToPlayer := fmt.Sprintf("Результат броска\n%s\n%dd%d: %d", roll.Title, roll.DiceCount, roll.DiceSides, roll.RollResult)
	return context.Send(textToPlayer)
}

// END MASTER REQUEST ANSWER SCENARIO

//--------------------------------------

// MASTER REQUEST GET SCENARIO

func MasterRequestFirstUnanswered(context runtime.Context, s *bot.Session) error {
	masterRequest, err := db.GetFirstUnansweredMasterRequest(s.DB, context.ChatID())
	if err != nil {
		return err
	}

	formattedMessage := fmt.Sprintf("Мастер требует внимания\n\n%s", masterRequest.TextRequest)
	err = context.SendTo(masterRequest.To, formattedMessage, ui.AnswerMasterKeyboard(masterRequest))

	return err
}

func MasterRequestGetFirstAnswered(context runtime.Context, s *bot.Session) error {
	log.Print("DEBUG: Start of the function")
	masterRequest, err := db.GetFirstAnsweredMasterRequest(s.DB)
	if err != nil {
		log.Print("DEBUG: get first answered master request != nil")
		if errors.Is(err, sql.ErrNoRows) {
			log.Print("DEBUG: err is ErrNoRows")
			return context.Send("No entries found")
		}
		return err
	}

	log.Print("DEBUG: Before Get User")
	player, err := db.GetUserByID(s.DB, int64(masterRequest.To))
	if err != nil {
		log.Print("DEBUG: User retrieval, err != nil")
		return err
	}

	var diceRollsText strings.Builder
	for _, roll := range masterRequest.RollRequests {
		text := fmt.Sprintf("Roll: %s\n%d (%dd%d)\n\n---", roll.Title, roll.RollResult, roll.DiceCount, roll.DiceSides)
		diceRollsText.WriteString(text)
	}

	formattedMessage := fmt.Sprintf("Answered By: %s\nAt: %s\n\nREQUEST:\n\n%s\n\nRESPONSE:\n\n%s\n\nDice Rolls:\n\n%s",
		player.PlayerName,
		masterRequest.UpdatedAt.Format(time.RFC850),
		masterRequest.TextRequest,
		masterRequest.TextResponse,
		diceRollsText.String())
	err = context.Send(formattedMessage, ui.MasterRequestMarkRead(int64(masterRequest.ID)))

	return err
}

func MasterRequestMarkAsRead(context runtime.Context, s *bot.Session) error {
	_, cbData := ParseCallbackDataString(context.Callback())
	var masterRequestID int64
	_, err := fmt.Sscanf(cbData, "%d", &masterRequestID)
	if err != nil {
		return err
	}

	err = db.UpdateMasterRequestState(s.DB, int64(masterRequestID), bot.MRCheckedByMaster)
	if err != nil {
		return err
	}

	return context.Send("Master Request marked as read!")
}

// END MASTER REQUEST GET SCENARIO
