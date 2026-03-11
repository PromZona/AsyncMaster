package list_faction

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/PromZona/AsyncMaster/internal/app/flows/list_factions/contract"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

type Session struct {
	DB   *sql.DB
	Done bool
}

func (s *Session) Name() string {
	return "list_factions"
}

func (s *Session) IsSupportedCallback(cb string) bool {
	callbacks := []string{
		contract.CBListFactions,
	}
	return slices.Contains(callbacks, cb)
}

func (s *Session) IsDone() bool {
	return s.Done
}

func (s *Session) DispatchCallback(context runtime.Context, cbUnique string, cbData string) error {
	switch cbUnique {
	case contract.CBListFactions:
		return handleListFactions(context, s)
	default:
		return fmt.Errorf("met unexpected callback unique: %s", cbUnique)
	}
}

func (s *Session) DispatchText(context runtime.Context) error {
	return fmt.Errorf("list_factions route does not support Text dispatch")
}

func NewSession(db *sql.DB) *Session {
	return &Session{
		DB:   db,
		Done: false,
	}
}
