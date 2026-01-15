package masterrequest

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/PromZona/AsyncMaster/internal/app/flows/_example/contract"
	tele "gopkg.in/telebot.v4"
)

type Session struct {
	DB        *sql.DB
	UserState State
	Done      bool
}

func (s *Session) Name() string {
	return "_example"
}

func (s *Session) IsSupportedCallback(cb string) bool {
	callbacks := []string{
		contract.CBExample,
	}
	return slices.Contains(callbacks, cb)
}

func (s *Session) IsDone() bool {
	return s.Done
}

func (s *Session) DispatchCallback(context tele.Context, cbUnique string, cbData string) error {
	switch cbUnique {
	case contract.CBExample:
		return handleStartFlow(context, s)
	default:
		return fmt.Errorf("met unexpected callback unique: %s", cbUnique)
	}
}

func (s *Session) DispatchText(context tele.Context) error {
	switch s.UserState {
	case AwaitTextMessage:
		return handleTextMessage(context, s)
	default:
		return fmt.Errorf("met unsupported state while handling master request: %d", s.UserState)
	}
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
	AwaitTextMessage State = 1
)
