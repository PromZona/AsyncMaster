package listmessages

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/PromZona/AsyncMaster/internal/app/flows/list_messages/contract"
	tele "gopkg.in/telebot.v4"
)

type Session struct {
	DB        *sql.DB
	UserState State
	Done      bool
}

func (s *Session) Name() string {
	return "list_messages"
}

func (s *Session) IsSupportedCallback(cb string) bool {
	callbacks := []string{
		contract.CBGetMessageList, contract.CBGetMessage,
	}
	return slices.Contains(callbacks, cb)
}

func (s *Session) IsDone() bool {
	return s.Done
}

func (s *Session) DispatchCallback(context tele.Context, cbUnique string, cbData string) error {
	switch cbUnique {
	case contract.CBGetMessageList:
		return handleStartFlow(context, s)
	case contract.CBGetMessage:
		return handleMessagePick(context, s, cbData)
	default:
		return fmt.Errorf("met unexpected callback unique: %s", cbUnique)
	}
}

func (s *Session) DispatchText(context tele.Context) error {
	return fmt.Errorf("met unsupported state while handling master request: %d", s.UserState)
}

func NewSession(db *sql.DB) *Session {
	return &Session{
		DB:        db,
		UserState: FlowStart,
		Done:      false,
	}
}

type State int

const (
	FlowStart        State = 0
	AwaitMessagePick State = 1
)
