package ui

import (
	"fmt"
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"

	answrmstrc "github.com/PromZona/AsyncMaster/internal/app/flows/answer_master/contract"
	listmstrreqc "github.com/PromZona/AsyncMaster/internal/app/flows/list_master_requests/contract"
	listmsgc "github.com/PromZona/AsyncMaster/internal/app/flows/list_messages/contract"
	mstrreqc "github.com/PromZona/AsyncMaster/internal/app/flows/master_request/contract"
	sendmsgc "github.com/PromZona/AsyncMaster/internal/app/flows/send_message/contract"

	tele "gopkg.in/telebot.v4"
)

func MainMenuKeyboard(context tele.Context, user *bot.UserData) error {

	if user.Role == bot.RoleMaster {
		return context.Send("Some master info here", masterMenu())
	}

	if user.Role == bot.RolePlayer {
		return context.Send("Some player info here", playerMenu())
	}

	return context.Send("Your role is not supported, contact administrator")
}

func PlayerNamesKeyboard(playerNames []string, chatIDs []int64) *tele.ReplyMarkup {
	if len(playerNames) != len(chatIDs) {
		log.Print("Error while creating keyboard, playerNames are not the same size as chatIDs: ", len(playerNames), "; ", len(chatIDs))
		return nil
	}

	result := &tele.ReplyMarkup{}
	var btnPlayerNames []tele.Btn

	for i, name := range playerNames {
		dataString := fmt.Sprintf("%s:%d", name, chatIDs[i])
		btnPlayerNames = append(btnPlayerNames, result.Data(name, sendmsgc.CBPlayerNames, dataString))
	}

	result.Inline(
		result.Row(btnPlayerNames...),
		result.Row(cancelButton()),
	)

	return result
}

func YesNoKeyboard() *tele.ReplyMarkup {
	result := &tele.ReplyMarkup{}

	btnNo := result.Data("No", "no")
	btnYes := result.Data("Yes", "yes")

	result.Inline(
		result.Row(btnNo, btnYes),
		result.Row(cancelButton()),
	)

	return result
}

func AnswerMasterKeyboard(masterRequest *bot.MasterRequest) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	allRows := make([]tele.Row, 0, len(masterRequest.RollRequests)+1)

	btnReply := menu.Data("Reply to Master", answrmstrc.CBReplyToMaster, fmt.Sprintf("%d", masterRequest.ID))
	allRows = append(allRows, menu.Row(btnReply))

	for _, roll := range masterRequest.RollRequests {
		text := fmt.Sprintf("%dd%d: %s", roll.DiceCount, roll.DiceSides, roll.Title)
		data := fmt.Sprintf("%d", roll.ID)
		btnRoll := menu.Data(text, answrmstrc.CBRollRequest, data)
		allRows = append(allRows, menu.Row(btnRoll))
	}
	menu.Inline(allRows...)

	return menu
}

func UserMessagesKeyboard(transactions []*bot.MessageTransaction) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{
		ResizeKeyboard: true,
	}

	allRows := make([]tele.Row, 0, len(transactions)+1)
	for _, t := range transactions {
		text := fmt.Sprintf("%s", t.Message.Title)
		data := fmt.Sprintf("%d", t.ID)
		btnMessage := menu.Data(text, listmsgc.CBGetMessage, data)
		allRows = append(allRows, menu.Row(btnMessage))
	}

	allRows = append(allRows, menu.Row(cancelButton()))
	menu.Inline(allRows...)
	return menu
}

func cancelButton() tele.Btn {
	btnCancel := tele.Btn{
		Unique: "cancel",
		Text:   "Cancel",
	}
	return btnCancel
}

func masterMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	btnSendMasters := menu.Data("Send Message", sendmsgc.CBSend)
	btnMasterRequest := menu.Data("Master Request", mstrreqc.CBStartMasterRequest)
	menu.Inline(
		menu.Row(btnSendMasters, btnMasterRequest),
	)
	return menu
}

func playerMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	btnSend := menu.Data("Send Message", sendmsgc.CBSend)
	btnMessages := menu.Data("My Messages", listmsgc.CBGetMessageList)
	btnMasterRequests := menu.Data("Last Master Request", listmstrreqc.CBGetMasterRequests)
	menu.Inline(
		menu.Row(btnSend),
		menu.Row(btnMessages),
		menu.Row(btnMasterRequests),
	)
	return menu
}
