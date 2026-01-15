package answermaster

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/flows/answer_master/contract"
	tele "gopkg.in/telebot.v4"
)

type Session struct {
	DB            *sql.DB
	UserState     State
	Done          bool
	MasterRequest *bot.MasterRequest
}

func (s *Session) Name() string {
	return "answer_master"
}
func (s *Session) IsSupportedCallback(cbUnique string) bool {
	supported := []string{
		contract.CBReplyToMaster, contract.CBRollRequest,
	}
	return slices.Contains(supported, cbUnique)
}
func (s *Session) IsDone() bool {
	return s.Done
}

func (s *Session) DispatchCallback(context tele.Context, cbUnique string, cbData string) error {
	switch cbUnique {
	case contract.CBReplyToMaster:
		return handleReplyToMaster(context, s, cbData)
	case contract.CBRollRequest:
		return handleRoll(context, s, cbData)
	default:
		return fmt.Errorf("met unexpected callback unique: %s", cbUnique)
	}
}

func (s *Session) DispatchText(context tele.Context) error {
	switch s.UserState {
	case AwaitText:
		return handleText(context, s)
	default:
		return fmt.Errorf("met unsupported state while handling master request: %d", s.UserState)
	}
}

func NewSession(db *sql.DB) *Session {
	return &Session{
		DB:            db,
		UserState:     Idle,
		Done:          false,
		MasterRequest: nil,
	}
}

type State int

const (
	Idle      State = 0
	AwaitText State = 0
)
