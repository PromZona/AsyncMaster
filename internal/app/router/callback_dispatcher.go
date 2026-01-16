package router

import (
	"database/sql"
	"log"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/flows/common"
	"github.com/PromZona/AsyncMaster/internal/app/ui"

	answermaster "github.com/PromZona/AsyncMaster/internal/app/flows/answer_master"
	listmasterrequests "github.com/PromZona/AsyncMaster/internal/app/flows/list_master_requests"
	listmessages "github.com/PromZona/AsyncMaster/internal/app/flows/list_messages"
	masterrequest "github.com/PromZona/AsyncMaster/internal/app/flows/master_request"
	sendmessage "github.com/PromZona/AsyncMaster/internal/app/flows/send_message"

	answrmstrc "github.com/PromZona/AsyncMaster/internal/app/flows/answer_master/contract"
	listmstrreqc "github.com/PromZona/AsyncMaster/internal/app/flows/list_master_requests/contract"
	listmsgc "github.com/PromZona/AsyncMaster/internal/app/flows/list_messages/contract"
	mstrreqc "github.com/PromZona/AsyncMaster/internal/app/flows/master_request/contract"
	sendmsgc "github.com/PromZona/AsyncMaster/internal/app/flows/send_message/contract"

	tele "gopkg.in/telebot.v4"
)

func DispatchCallback(context tele.Context, b *bot.BotData) error {
	context.Respond()

	chatID := context.Chat().ID
	rawCallbackData := context.Callback().Data
	cbUnique, cbData := parseCallbackDataString(rawCallbackData)

	if cbUnique == common.CBCancel {
		return common.HandleCancelButton(context, b)
	}

	session := b.GetUserSession(chatID)
	if session == nil {
		factory, ok := UniqueToSessionFactory[cbUnique]
		if !ok {
			log.Printf("Met not start flow cbUnique: %s. Have you forgot to register it?", cbUnique)
			context.Send("Please start the action properly by pressing button from the menu")
			return ui.MainMenuKeyboard(context, db.GetUserByID(b.DB, chatID).Role)
		}
		session = factory(b.DB)
		b.UserActiveSessions[chatID] = session
		log.Printf("Session created for chatID: %d, %s", chatID, session.Name())
	}

	if !session.IsSupportedCallback(cbUnique) {
		return context.Send("You are performaing a different action right now, please finish it first")
	}

	err := session.DispatchCallback(context, cbUnique, cbData)

	if session.IsDone() {
		delete(b.UserActiveSessions, chatID)
		log.Printf("Session delete for chatID: %d", chatID)
	}

	return err
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

type SessionFactory func(db *sql.DB) bot.FlowSession

var UniqueToSessionFactory = map[string]SessionFactory{
	sendmsgc.CBSend: func(db *sql.DB) bot.FlowSession {
		return sendmessage.NewSession(db)
	},
	mstrreqc.CBStartMasterRequest: func(db *sql.DB) bot.FlowSession {
		return masterrequest.NewSession(db)
	},
	answrmstrc.CBReplyToMaster: func(db *sql.DB) bot.FlowSession {
		return answermaster.NewSession(db)
	},
	answrmstrc.CBRollRequest: func(db *sql.DB) bot.FlowSession {
		return answermaster.NewSession(db)
	},
	listmsgc.CBGetMessageList: func(db *sql.DB) bot.FlowSession {
		return listmessages.NewSession(db)
	},
	listmstrreqc.CBGetMasterRequests: func(db *sql.DB) bot.FlowSession {
		return listmasterrequests.NewSession(db)
	},
}
