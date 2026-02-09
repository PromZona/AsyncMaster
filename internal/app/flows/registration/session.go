package registration

import (
	"database/sql"
	"fmt"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	tele "gopkg.in/telebot.v4"
)

type Session struct {
	DB        *sql.DB
	UserState State
	Done      bool
	User      *bot.UserData
}

func (s *Session) Name() string {
	return "registration"
}

func (s *Session) IsDone() bool {
	return s.Done
}

func (s *Session) IsSupportedCallback(cbUnique string) bool {
	return false
}

func (s *Session) DispatchCallback(context tele.Context, cbUnique string, cbData string) error {
	return fmt.Errorf("registration callback does not support callbacks. cbUnique: %s, cbData: %s", cbUnique, cbData)
}

func (s *Session) DispatchText(context tele.Context) error {
	switch s.UserState {
	case AwaitPassword:
		return handlePassword(context, s)
	case AwaitCodename:
		return handlePlayerName(context, s)
	case AwaitFactionName:
		return handleFactionName(context, s)
	case AwaitFactionDescription:
		return handleFactionDescription(context, s)
	default:
		return fmt.Errorf("registration met unsupported user state while dispatching text: %d", s.UserState)
	}
}

type State int

const (
	AwaitPassword           State = 0
	AwaitCodename           State = 1
	AwaitFactionName        State = 2
	AwaitFactionDescription State = 3
)
