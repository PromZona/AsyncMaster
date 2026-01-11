package ui

import (
	"fmt"
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
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
		btnPlayerNames = append(btnPlayerNames, result.Data(name, "player_names", dataString))
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

	btnReply := menu.Data("Reply to Master", "reply_to_master", fmt.Sprintf("%d", masterRequest.ID))
	allRows = append(allRows, menu.Row(btnReply))

	for _, roll := range masterRequest.RollRequests {
		text := fmt.Sprintf("%dd%d: %s", roll.DiceCount, roll.DiceSides, roll.Title)
		data := fmt.Sprintf("%d", roll.ID)
		btnRoll := menu.Data(text, "roll_request", data)
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
	btnSendMasters := menu.Data("Send To Player", "send")
	btnMasterRequest := menu.Data("Master Request", "start_master_request")
	menu.Inline(
		menu.Row(btnSendMasters, btnMasterRequest),
	)
	return menu
}

func playerMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	btnSend := menu.Data("Send To Player", "send")
	menu.Inline(
		menu.Row(btnSend),
	)
	return menu
}
