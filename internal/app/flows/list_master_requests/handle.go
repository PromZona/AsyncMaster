package listmasterrequests

import (
	"fmt"

	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

func handleStartFlow(context tele.Context, s *Session) error {
	if s.UserState != FlowStart {
		return context.Send("This action is not available right now, finish previous action first")
	}

	masterRequest, err := db.GetLastMasterRequest(s.DB, context.Chat().ID)
	if err != nil {
		return err
	}

	formattedMessage := fmt.Sprintf("MASTER REQUEST\n\n%s", masterRequest.TextRequest)
	_, err = context.Bot().Send(masterRequest.To, formattedMessage, ui.AnswerMasterKeyboard(masterRequest))

	s.Done = true
	return err
}
