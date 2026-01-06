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
				if _, ok := b.UserSessionState[chatID]; !ok {
					b.UserSessionState[chatID] = bot.UserStateAwaitPassword
				}
				return registration.HandleRegisterUser(context, b)
			}

			if _, ok := b.UserSessionState[chatID]; !ok {
				b.UserSessionState[chatID] = bot.UserStateDefault
			}
			return next(context)
		}
	}
}
