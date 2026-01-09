package sendmessage

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
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
		Send, PlayerNames, NoTitle, YesTitle,
	}
	return slices.Contains(slice, cbUnique)
}

func (s *Session) DispatchCallback(context tele.Context, cbUnique string, cbData string) error {
	switch cbUnique {
	case "send":
		return handleInitialSend(context, s)
	case "player_names":
		return handlePlayerName(context, s, cbData)
	case "no_title":
		return handleNoTitle(context, s)
	case "yes_title":
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

const (
	Send        string = "send"
	PlayerNames        = "player_names"
	NoTitle            = "no_title"
	YesTitle           = "yes_title"
)
