package dispatcher

import (
	"fmt"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	sendmessage "github.com/PromZona/AsyncMaster/internal/app/flows/send_message"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

func Text(context tele.Context, b *bot.BotData) error {
	chatID := context.Chat().ID
	state := b.UserSessionState[chatID]

	switch state {
	case bot.UserStateDefault:
		return ui.MainMenu(context, b)
	case bot.UserStateAwaitMessage, bot.UserStateAwaitTitle:
		return sendmessage.DispatchText(context, b)
	default:
		return fmt.Errorf("dispatcher text met unsupported state: %d", state)
	}
}
