package middleware

import (
	"fmt"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/handlers"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func RegistrationCheck(b *bot.BotData) runtime.Middleware {
	return func(next runtime.Handler) runtime.Handler {
		return func(context runtime.Context) error {
			chatID := context.ChatID()

			if !db.EnsureUserExist(b.DB, context.ChatID()) {
				session := b.GetSession(chatID)
				if session == nil {
					session = bot.NewSession(b.DB)
					session.RegistrationState = bot.RegistrationAwaitPassword
					b.Sessions[chatID] = session
				}

				var err error
				switch session.RegistrationState {
				case bot.RegistrationAwaitPassword:
					err = handlers.RegistrationPassword(context, session)
				case bot.RegistrationAwaitCodename:
					err = handlers.RegistrationPlayerName(context, session)
				case bot.RegistrationAwaitFactionName:
					err = handlers.RegistrationFactionName(context, session)
				case bot.RegistrationAwaitFactionDescription:
					err = handlers.RegistrationFactionDescription(context, session)
				case bot.RegistrationNotActive:
					fallthrough
				case bot.RegistrationFinished:
					fallthrough
				default:
					return fmt.Errorf("Met unexpected state while registering user: %d", session.RegistrationState)
				}

				if session.RegistrationState != bot.RegistrationFinished {
					return err
				}

			} else {
				session := b.GetSession(chatID)
				if session == nil {
					session = bot.NewSession(b.DB)
					b.Sessions[chatID] = session
				}
				session.RegistrationState = bot.RegistrationFinished
			}
			return next(context)
		}
	}
}
