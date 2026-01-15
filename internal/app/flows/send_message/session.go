package sendmessage

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/flows/send_message/contract"
	tele "gopkg.in/telebot.v4"
)

type Session struct {
	DB               *sql.DB
	UserState        State
	DraftMessage     *bot.Message
	DraftTransaction *bot.MessageTransaction
	Done             bool
}

func (s *Session) Name() string {
	return "send_message"
}

func (s *Session) IsDone() bool {
	return s.Done
}

func (s *Session) IsSupportedCallback(cbUnique string) bool {
	slice := []string{
		contract.CBSend, contract.CBPlayerNames, contract.CBTitleNo, contract.CBTitleYes,
	}
	return slices.Contains(slice, cbUnique)
}

func (s *Session) DispatchCallback(context tele.Context, cbUnique string, cbData string) error {
	switch cbUnique {
	case contract.CBSend:
		return handleInitialSend(context, s)
	case contract.CBPlayerNames:
		return handlePlayerName(context, s, cbData)
	case contract.CBTitleNo:
		return handleNoTitle(context, s)
	case contract.CBTitleYes:
		return handleYesTitle(context, s)
	default:
		return fmt.Errorf("sendmessage met unsupported <callback unique> while dispatching callback: %s", cbUnique)
	}
}

func (s *Session) DispatchText(context tele.Context) error {
	switch s.UserState {
	case AwaitMessage:
		return handleMessageText(context, s)
	case AwaitTitle:
		return handleMessageTitle(context, s)
	default:
		return fmt.Errorf("sendmessage met unsupported user state while dispatching text: %d", s.UserState)
	}
}

func NewSession(db *sql.DB) *Session {
	return &Session{
		DB:               db,
		UserState:        FlowStart,
		DraftMessage:     &bot.Message{},
		DraftTransaction: &bot.MessageTransaction{},
		Done:             false,
	}
}

type State int

const (
	FlowStart          State = 0
	AwaitResipient           = 1
	AwaitMessage             = 2
	AwaitTitleDecision       = 3
	AwaitTitle               = 4
)
