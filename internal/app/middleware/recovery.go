package middleware

import (
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	tele "gopkg.in/telebot.v4"
)

func ErrorRecovery(b *bot.BotData) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(context tele.Context) error {
			err := next(context)
			if err != nil {
				log.Print("ERROR: ", err)
				b.UserSessionState[context.Chat().ID] = bot.UserStateDefault
				context.Send("Error occured while proccesing your request. Please, contact administrator. Returning you to main menu")
			}
			return err
		}
	}
}
