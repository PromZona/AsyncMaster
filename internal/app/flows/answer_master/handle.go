package answermaster

import (
	"fmt"
	"math/rand"

	"github.com/PromZona/AsyncMaster/internal/app/db"
	tele "gopkg.in/telebot.v4"
)

func handleReplyToMaster(context tele.Context, s *Session, cbData string) error {
	if s.UserState != Idle {
		return context.Send("Expected other action")
	}

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
	s.UserState = AwaitText
	return context.Send("Write your reply:")
}

func handleText(context tele.Context, s *Session) error {
	if s.UserState != AwaitText {
		return context.Send("Expected other action")
	}

	if s.MasterRequest == nil {
		context.Send("Met unexpected error, contact administrator")
		return fmt.Errorf("master request is nil while trying to write text into it")
	}

	text := context.Text()
	s.MasterRequest.TextResponse = text

	err := db.UpdateMasterRequest(s.DB, s.MasterRequest)
	if err != nil {
		return err
	}

	s.Done = true
	return context.Send("Answer accepted. Please, roll requested dices if there are any left")
}

func handleRoll(context tele.Context, s *Session, cbData string) error {
	if s.UserState != Idle {
		return context.Send("Expected other action")
	}

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

	s.Done = true
	textToPlayer := fmt.Sprintf("Roll result:\n%s\n%dd%d: %d", roll.Title, roll.DiceCount, roll.DiceSides, roll.RollResult)
	return context.Send(textToPlayer)
}
