package ui

import (
	"fmt"
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
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
			runtime.Button{Text: name, Unique: string(bot.HID_message_player_name), Data: dataString})
	}

	rows := []runtime.Row{}
	if len(btnPlayerNames) > 0 {
		rows = append(rows, runtime.Row(btnPlayerNames))
	}
	rows = append(rows, runtime.Row{cancelButton()})

	keyboard := runtime.Keyboard(rows)
	return keyboard
}

func YesNoKeyboard(yes bot.HandlerID, no bot.HandlerID) runtime.Keyboard {
	btnNo := runtime.Button{Text: "No", Unique: string(no)}
	btnYes := runtime.Button{Text: "Yes", Unique: string(yes)}

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

	btnReply := runtime.Button{Text: "Reply to Master", Unique: string(bot.HID_mr_answer_text), Data: fmt.Sprintf("%d", masterRequest.ID)}
	allRows = append(allRows, runtime.Row{btnReply})

	for _, roll := range masterRequest.RollRequests {
		text := fmt.Sprintf("%dd%d: %s", roll.DiceCount, roll.DiceSides, roll.Title)
		data := fmt.Sprintf("%d", roll.ID)
		btnRoll := runtime.Button{Text: text, Unique: string(bot.HID_mr_answer_roll), Data: data}
		allRows = append(allRows, runtime.Row{btnRoll})
	}

	keyboard := runtime.Keyboard(allRows)
	return keyboard
}

func UserMessagesKeyboard(transactions []*bot.MessageTransaction) runtime.Keyboard {
	allRows := make([]runtime.Row, 0, len(transactions)+1)
	for _, t := range transactions {
		text := t.Message.Title
		data := fmt.Sprintf("%d", t.ID)
		btnMessage := runtime.Button{Text: text, Unique: string(bot.HID_message_get), Data: data}
		allRows = append(allRows, runtime.Row{btnMessage})
	}

	allRows = append(allRows, runtime.Row{cancelButton()})
	keyboard := runtime.Keyboard(allRows)
	return keyboard
}

func MasterRequestMarkRead(requestID int64) runtime.Keyboard {
	data := fmt.Sprintf("%d", requestID)
	btnMarkRead := runtime.Button{Text: "Mark as Read", Unique: string(bot.HID_mr_master_mark_read), Data: data}

	keyboard := runtime.Keyboard{runtime.Row{btnMarkRead}}
	return keyboard
}

func FactionsUpdateFactionsList(users []bot.UserData) runtime.Keyboard {
	allRows := make([]runtime.Row, 0, len(users)+1)

	for _, u := range users {
		text := fmt.Sprintf("%s (%s)", u.Faction.Name, u.PlayerName)
		data := fmt.Sprintf("%d", u.ChatID)
		btn := runtime.Button{Text: text, Unique: string(bot.HID_factions_update_player), Data: data}
		allRows = append(allRows, runtime.Row{btn})
	}

	allRows = append(allRows, runtime.Row{cancelButton()})
	keyboard := runtime.Keyboard(allRows)
	return keyboard
}

func FactionsUpdateWhatToUpdate() runtime.Keyboard {

	btnResources := runtime.Button{Text: "Resources", Unique: string(bot.HID_factions_update_resources), Data: ""}

	keyboard := runtime.Keyboard{
		{btnResources},
		{cancelButton()},
	}
	return keyboard
}

func CancelMenu() runtime.Keyboard {
	return runtime.Keyboard{
		{cancelButton()},
	}
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
		Unique: string(bot.HID_message_send),
	}
	btnSendEveryone := runtime.Button{
		Text:   "Send Message to Everyone",
		Unique: string(bot.HID_message_send_everyone),
	}
	btnMasterRequest := runtime.Button{
		Text:   "Master Request",
		Unique: string(bot.HID_mr_send),
	}
	btnMasterRequestEveryone := runtime.Button{
		Text:   "Master Request to Everyone",
		Unique: string(bot.HID_mr_send_everyone),
	}
	btnCheckAnsweredMasterRequest := runtime.Button{
		Text:   "Answered Request",
		Unique: string(bot.HID_mr_master_get),
	}
	btnFactionsUpdate := runtime.Button{
		Text:   "Factions Update",
		Unique: string(bot.HID_factions_update),
	}
	btnFactionsList := runtime.Button{
		Text:   "Factions List",
		Unique: string(bot.HID_factions_list),
	}

	keyboard := runtime.Keyboard{
		{btnSendMasters, btnSendEveryone},
		{btnMasterRequest, btnMasterRequestEveryone},
		{btnFactionsList, btnFactionsUpdate},
		{btnCheckAnsweredMasterRequest},
	}
	return keyboard
}

func playerMenu(unansweredMRCount int) runtime.Keyboard {
	btnSend := runtime.Button{Text: "Send Message", Unique: string(bot.HID_message_send)}
	btnMessages := runtime.Button{Text: "My Messages", Unique: string(bot.HID_message_list)}
	btnFactions := runtime.Button{Text: "Factions", Unique: string(bot.HID_factions_list)}

	masterRequestEmoji := "🟢"
	if unansweredMRCount > 0 {
		masterRequestEmoji = "🔴"
	}
	masterRequestText := fmt.Sprintf("%s Answer Master Request (%d unanswered)",
		masterRequestEmoji,
		unansweredMRCount)
	btnMasterRequests := runtime.Button{Text: masterRequestText, Unique: string(bot.HID_mr_player_get)}

	rows := []runtime.Row{
		{btnMasterRequests},
		{btnSend, btnMessages},
		{btnFactions},
	}

	keyboard := runtime.Keyboard(rows)
	return keyboard
}
