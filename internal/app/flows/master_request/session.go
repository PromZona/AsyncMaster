package masterrequest

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	tele "gopkg.in/telebot.v4"
)

type Session struct {
	DB           *sql.DB
	UserState    State
	Done         bool
	RequestData  *bot.MasterRequest
	RollRequests []*bot.RollRequest
}

func (s *Session) Name() string {
	return "master_request"
}

func (s *Session) IsSupportedCallback(cb string) bool {
	callbacks := []string{
		"start_master_request", "player_names", "yes", "no",
	}
	return slices.Contains(callbacks, cb)
}

func (s *Session) IsDone() bool {
	return s.Done
}

func (s *Session) DispatchCallback(context tele.Context, cbUnique string, cbData string) error {
	switch cbUnique {
	case "start_master_request":
		return handleStartFlow(context, s)
	case "player_names":
		return handleResipient(context, s, cbData)
	case "yes":
		return handleYes(context, s)
	case "no":
		return handleNo(context, s)
	default:
		return fmt.Errorf("met unexpected callback unique: %s", cbUnique)
	}
}

func (s *Session) DispatchText(context tele.Context) error {
	switch s.UserState {
	case AwaitText:
		return handleText(context, s)
	case AwaitRoll:
		return handleRoll(context, s)
	default:
		return fmt.Errorf("met unsupported state while handling master request: %d", s.UserState)
	}
}

func NewSession(db *sql.DB) *Session {
	return &Session{
		DB:           db,
		UserState:    FlowStart,
		Done:         false,
		RequestData:  &bot.MasterRequest{},
		RollRequests: make([]*bot.RollRequest, 0),
	}
}

type State int

const (
	FlowStart         State = 0
	AwaitResipient    State = 1
	AwaitText         State = 2
	AwaitRollDecision State = 3
	AwaitRoll         State = 4
)
