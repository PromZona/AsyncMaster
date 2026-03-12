package masterrequest

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
)

func handleStartFlow(context runtime.Context, s *Session) error {
	if s.UserState != FlowStart {
		return context.Send("This action is not available right now, finish previous action first")
	}

	playerNames, chatIDs, err := db.GetNamesAndChatIDsOfPlayers(s.DB)
	if err != nil {
		return err
	}

	if s.IsSendEveryone {
		s.Resipients = append(s.Resipients, chatIDs...)
		s.UserState = AwaitText
		return context.Send("Type text which will be sent to players")
	}

	s.UserState = AwaitResipient
	return context.Send("Pick a player to send to", ui.PlayerNamesKeyboard(playerNames, chatIDs))
}

func handleResipient(context runtime.Context, s *Session, cbData string) error {
	if s.UserState != AwaitResipient {
		return context.Send("This action is not available right now, finish previous action first")
	}

	splited := strings.SplitAfterN(cbData, ":", 2)
	if len(splited) != 2 {
		return fmt.Errorf("splitting callbackdata, met unexpected amount of data: %d", len(splited))
	}

	toChatID, err := strconv.ParseInt(splited[1], 10, 64)
	if err != nil {
		return err
	}

	s.Resipients = append(s.Resipients, toChatID)
	s.UserState = AwaitText
	return context.Send("Type text which will be sent to player")
}

func handleText(context runtime.Context, s *Session) error {
	if s.UserState != AwaitText {
		return context.Send("This action is not available right now, finish previous action first")
	}

	s.RequestData.TextRequest = context.MessageText()

	s.UserState = AwaitRollDecision
	return context.Send("Do you want to add dice request?", ui.YesNoKeyboard())
}

func handleYes(context runtime.Context, s *Session) error {
	if s.UserState != AwaitRollDecision {
		return context.Send("This action is not available right now, finish previous action first")
	}

	s.UserState = AwaitRoll
	return sendRollQuestion(context)
}

func handleNo(context runtime.Context, s *Session) error {
	if s.UserState != AwaitRollDecision {
		return context.Send("This action is not available right now, finish previous action first")
	}

	return finilize(context, s)
}

func handleRoll(context runtime.Context, s *Session) error {
	if s.UserState != AwaitRoll {
		return context.Send("This action is not available right now, finish previous action first")
	}

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
	s.UserState = AwaitRollDecision
	return context.Send("Do you want to add dice request?", ui.YesNoKeyboard())
}

func finilize(context runtime.Context, s *Session) error {
	masterRequest := s.RequestData
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
	}

	formattedMessage := fmt.Sprintf("MASTER REQUEST\n\n%s", masterRequest.TextRequest)
	context.SendTo(masterRequest.To, formattedMessage, ui.AnswerMasterKeyboard(masterRequest))

	s.Done = true
	return context.Send("Message send to resipient")
}

func sendRollQuestion(context runtime.Context) error {
	return context.Send("Write a roll in a format: <count>d<dice> <name of the role>\nExample: 1d6 Roll on money")
}
