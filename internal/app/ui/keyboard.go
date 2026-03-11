package ui

import (
	"fmt"
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"

	answrmstrc "github.com/PromZona/AsyncMaster/internal/app/flows/answer_master/contract"
	listfctns "github.com/PromZona/AsyncMaster/internal/app/flows/list_factions/contract"
	listmstrreqc "github.com/PromZona/AsyncMaster/internal/app/flows/list_master_requests/contract"
	listmsgc "github.com/PromZona/AsyncMaster/internal/app/flows/list_messages/contract"
	mstrreqc "github.com/PromZona/AsyncMaster/internal/app/flows/master_request/contract"
	sendmsgc "github.com/PromZona/AsyncMaster/internal/app/flows/send_message/contract"
)

func MainMenuPlayerKeyboard(context runtime.Context, user *bot.UserData, menuData *PlayerMenu) error {
	menuText := fmt.Sprintf(`Menu
		Player: %s
		Faction: %s
		%s
		---
		Resources:
		%s
		---
		You have %d unanswered Master Requests`,
		menuData.PlayerName,
		menuData.FactionName,
		menuData.FactionDescription,
		menuData.FactionResources,
		menuData.UnansweredMasterRequests)

	return context.Send(menuText, playerMenu(menuData.UnansweredMasterRequests))
}

func MainMenuMasterKeyboard(context runtime.Context, user *bot.UserData, menuData *MasterMenu) error {
	menuText := fmt.Sprintf("Menu:\nMaster")

	return context.Send(menuText, masterMenu())
}

func PlayerNamesKeyboard(playerNames []string, chatIDs []int64) runtime.Keyboard {
	if len(playerNames) != len(chatIDs) {
		log.Print("Error while creating keyboard, playerNames are not the same size as chatIDs: ", len(playerNames), "; ", len(chatIDs))
		return nil
	}

	var btnPlayerNames []runtime.Button

	for i, name := range playerNames {
		dataString := fmt.Sprintf("%s:%d", name, chatIDs[i])
		btnPlayerNames = append(
			btnPlayerNames,
			runtime.Button{Text: name, Unique: sendmsgc.CBPlayerNames, Data: dataString})
	}

	rows := []runtime.Row{}
	if len(btnPlayerNames) > 0 {
		rows = append(rows, runtime.Row(btnPlayerNames))
	}
	rows = append(rows, runtime.Row{cancelButton()})

	keyboard := runtime.Keyboard(rows)
	return keyboard
}

func YesNoKeyboard() runtime.Keyboard {
	btnNo := runtime.Button{Text: "No", Unique: "no"}
	btnYes := runtime.Button{Text: "Yes", Unique: "yes"}

	rows := []runtime.Row{}
	rows = append(rows, runtime.Row{
		btnNo,
		btnYes,
	})
	rows = append(rows, runtime.Row{cancelButton()})

	keyboard := runtime.Keyboard(rows)
	return keyboard
}

func AnswerMasterKeyboard(masterRequest *bot.MasterRequest) runtime.Keyboard {
	allRows := make([]runtime.Row, 0, len(masterRequest.RollRequests)+1)

	btnReply := runtime.Button{Text: "Reply to Master", Unique: answrmstrc.CBReplyToMaster, Data: fmt.Sprintf("%d", masterRequest.ID)}
	allRows = append(allRows, runtime.Row{btnReply})

	for _, roll := range masterRequest.RollRequests {
		text := fmt.Sprintf("%dd%d: %s", roll.DiceCount, roll.DiceSides, roll.Title)
		data := fmt.Sprintf("%d", roll.ID)
		btnRoll := runtime.Button{Text: text, Unique: answrmstrc.CBRollRequest, Data: data}
		allRows = append(allRows, runtime.Row{btnRoll})
	}

	keyboard := runtime.Keyboard(allRows)
	return keyboard
}

func UserMessagesKeyboard(transactions []*bot.MessageTransaction) runtime.Keyboard {
	allRows := make([]runtime.Row, 0, len(transactions)+1)
	for _, t := range transactions {
		text := fmt.Sprintf("%s", t.Message.Title)
		data := fmt.Sprintf("%d", t.ID)
		btnMessage := runtime.Button{Text: text, Unique: listmsgc.CBGetMessage, Data: data}
		allRows = append(allRows, runtime.Row{btnMessage})
	}

	allRows = append(allRows, runtime.Row{cancelButton()})
	keyboard := runtime.Keyboard(allRows)
	return keyboard
}

func MasterRequestMarkRead(requestID int64) runtime.Keyboard {
	data := fmt.Sprintf("%d", requestID)
	btnMarkRead := runtime.Button{Text: "Mark as Read", Unique: listmstrreqc.CBMarkAsRead, Data: data}

	keyboard := runtime.Keyboard{runtime.Row{btnMarkRead}}
	return keyboard
}

func cancelButton() runtime.Button {
	btnCancel := runtime.Button{
		Unique: "cancel",
		Text:   "Cancel",
	}
	return btnCancel
}

func masterMenu() runtime.Keyboard {

	btnSendMasters := runtime.Button{
		Text:   "Send Message",
		Unique: sendmsgc.CBSend,
	}
	btnSendEveryone := runtime.Button{
		Text:   "Send Message to Everyone",
		Unique: sendmsgc.CBSendEveryone,
	}
	btnMasterRequest := runtime.Button{
		Text:   "Master Request",
		Unique: mstrreqc.CBStartMasterRequest,
	}
	btnMasterRequestEveryone := runtime.Button{
		Text:   "Master Request to Everyone",
		Unique: mstrreqc.CBStartMasterRequestEveryone,
	}
	btnCheckAnsweredMasterRequest := runtime.Button{
		Text:   "Answered Request",
		Unique: listmstrreqc.CBGetAnsweredMasterRequest,
	}

	keyboard := runtime.Keyboard{
		{btnSendMasters, btnSendEveryone},
		{btnMasterRequest, btnMasterRequestEveryone},
		{btnCheckAnsweredMasterRequest},
	}
	return keyboard
}

func playerMenu(unansweredMRCount int) runtime.Keyboard {
	btnSend := runtime.Button{Text: "Send Message", Unique: sendmsgc.CBSend}
	btnMessages := runtime.Button{Text: "My Messages", Unique: listmsgc.CBGetMessageList}
	btnFactions := runtime.Button{Text: "Factions", Unique: listfctns.CBListFactions}

	masterRequestEmoji := "🟢"
	if unansweredMRCount > 0 {
		masterRequestEmoji = "🔴"
	}
	masterRequestText := fmt.Sprintf("%s Answer Master Request (%d unanswered) %s",
		masterRequestEmoji,
		unansweredMRCount,
		masterRequestEmoji)
	btnMasterRequests := runtime.Button{Text: masterRequestText, Unique: listmstrreqc.CBGetMasterRequests}

	rows := runtime.Row{
		btnMasterRequests,
		btnSend,
		btnMessages,
		btnFactions,
	}

	keyboard := runtime.Keyboard{rows}
	return keyboard
}
