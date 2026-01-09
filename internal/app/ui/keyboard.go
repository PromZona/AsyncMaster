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

func TitleQuestionKeyboard() *tele.ReplyMarkup {
	result := &tele.ReplyMarkup{}

	btnNo := result.Data("No", "no_title")
	btnYes := result.Data("Yes", "yes_title")

	result.Inline(
		result.Row(btnNo, btnYes),
		result.Row(cancelButton()),
	)

	return result
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
	btnMasterRequest := menu.Data("Master Request", "master_request")
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
