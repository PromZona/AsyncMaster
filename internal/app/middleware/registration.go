package middleware

import (
	tele "gopkg.in/telebot.v4"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/flows/registration"
)

func RegistrationCheck(b *bot.BotData) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(context tele.Context) error {
			chatID := context.Chat().ID

			if !db.EnsureUserExist(b.DB, context.Chat().ID) {
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
