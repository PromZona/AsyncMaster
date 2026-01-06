package router

import (
	"fmt"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	sendmessage "github.com/PromZona/AsyncMaster/internal/app/flows/send_message"
	tele "gopkg.in/telebot.v4"
)

func DispatchCallback(context tele.Context, b *bot.BotData) error {
	context.Respond()

	rawCallbackData := context.Callback().Data
	cbUnique, cbData := parseCallbackDataString(rawCallbackData)

	// log.Printf("callback unique   = %q", cbUnique)
	// log.Printf("callback data     = %q", cbData)
	// log.Printf("raw callback data = %q", rawCallbackData)

	switch cbUnique {
	case "send", "player_names", "no_title", "yes_title":
		return sendmessage.DispatchCallback(context, b, cbUnique, cbData)
	default:
		return fmt.Errorf("error, met unknown callback unique while dispatching callback: %s", cbUnique)
	}
}

func parseCallbackDataString(callbackData string) (unique, data string) {
	trimmed := strings.Trim(callbackData, "\f")
	parts := strings.SplitN(trimmed, "|", 2)
	count := len(parts)
	if count == 2 {
		return parts[0], parts[1]
	}
	if count == 1 {
		return parts[0], ""
	}
	return "", ""
}
