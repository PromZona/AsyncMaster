package middleware

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/flows/registration"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func RegistrationCheck(b *bot.BotData) runtime.Middleware {
	return func(next runtime.Handler) runtime.Handler {
		return func(context runtime.Context) error {
			chatID := context.ChatID()

			if !db.EnsureUserExist(b.DB, context.ChatID()) {
				session := b.GetUserSession(chatID)
				if session == nil {
					session = &registration.Session{
						DB:        b.DB,
						UserState: registration.AwaitPassword,
						Done:      false,
					}

					b.UserActiveSessions[chatID] = session
				}
				err := session.DispatchText(context)
				if !session.IsDone() {
					return err
				}
				delete(b.UserActiveSessions, chatID)
			}
			return next(context)
		}
	}
}
