package middleware

import (
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func ErrorRecovery(b *bot.BotData) runtime.Middleware {
	return func(next runtime.Handler) runtime.Handler {
		return func(context runtime.Context) error {
			err := next(context)
			if err != nil {
				log.Print("ERROR: ", err)
				b.ClearUserCache(context.ChatID())
				context.Send("Возникла техническая проблема. Срочно обратись к Мастеру")
			}
			return err
		}
	}
}
