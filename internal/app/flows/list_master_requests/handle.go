package listmasterrequests

import (
	"fmt"
	"strings"
	"time"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

// Player wants to answer master
func handleStartFlow(context tele.Context, s *Session) error {
	if s.UserState != FlowStart {
		return context.Send("This action is not available right now, finish previous action first")
	}

	masterRequest, err := db.GetFirstUnansweredMasterRequest(s.DB, context.Chat().ID)
	if err != nil {
		return err
	}

	formattedMessage := fmt.Sprintf("MASTER REQUEST\n\n%s", masterRequest.TextRequest)
	_, err = context.Bot().Send(masterRequest.To, formattedMessage, ui.AnswerMasterKeyboard(masterRequest))

	s.Done = true
	return err
}

// Master wants to check first Answered request
func handleGetFirstAnswered(context tele.Context, s *Session) error {
	if s.UserState != FlowStart {
		return context.Send("This action is not available right now, finish previous action first")
	}

	masterRequest, err := db.GetFirstAnsweredMasterRequest(s.DB)
	if err != nil {
		return err
	}

	player, err := db.GetUserByID(s.DB, int64(masterRequest.To))
	if err != nil {
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

	s.Done = true
	return err
}

func handleMarkAsRead(context tele.Context, s *Session, cbData string) error {
	if s.UserState != FlowStart {
		return context.Send("This action is not available right now, finish previous action first")
	}

	var masterRequestID int64
	_, err := fmt.Sscanf(cbData, "%d", &masterRequestID)
	if err != nil {
		return err
	}

	err = db.UpdateMasterRequestState(s.DB, int64(masterRequestID), bot.MRCheckedByMaster)
	if err != nil {
		return err
	}

	s.Done = true
	return context.Send("Master Request marked as read!")
}
