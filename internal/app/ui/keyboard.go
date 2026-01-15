package ui

import (
	"fmt"
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"

	answrmstrc "github.com/PromZona/AsyncMaster/internal/app/flows/answer_master/contract"
	mstrreqc "github.com/PromZona/AsyncMaster/internal/app/flows/master_request/contract"
	sendmsgc "github.com/PromZona/AsyncMaster/internal/app/flows/send_message/contract"

	tele "gopkg.in/telebot.v4"
)

func MainMenuKeyboard(context tele.Context, role bot.UserRole) error {
	var menu *tele.ReplyMarkup
	if role == bot.RoleMaster {
		menu = masterMenu()
	} else {
		menu = playerMenu()
	}
	return context.Send("Main menu", menu)
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
	// btnMessages := menu.Data("My Messages", ..."user_list_messages")
	// btnMasterRequests := menu.Data("Requests From Master", ..."list_master_requests")
	menu.Inline(
		menu.Row(btnSend),
		// menu.Row(btnMessages),
		//menu.Row(btnMasterRequests),
	)
	return menu
}
