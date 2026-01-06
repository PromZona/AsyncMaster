package ui

import (
	"fmt"
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	tele "gopkg.in/telebot.v4"
)

func MainMenuKeyboard(context tele.Context, b *bot.BotData) error {
	user := db.GetUser(b.DB, context.Chat().ID)
	var menu *tele.ReplyMarkup
	if user.Role == bot.RoleMaster {
		menu = b.MasterMenu
	} else {
		menu = b.PlayerMenu
	}
	return context.Send("Main menu", menu)
}

func PlayerNamesKeyboard(b *bot.BotData) *tele.ReplyMarkup {
	result := &tele.ReplyMarkup{}

	playerNames, chatIDs, err := db.GetUserPlayerNamesAndChatID(b.DB)
	if err != nil {
		log.Print("Error while creating keyboard, ", err)
		return nil
	}
	if len(playerNames) != len(chatIDs) {
		log.Print("Error while creating keyboard, playerNames are not the same size as chatIDs: ", len(playerNames), "; ", len(chatIDs))
		return nil
	}

	var btnPlayerNames []tele.Btn

	for i, name := range playerNames {
		dataString := fmt.Sprintf("%s:%d", name, chatIDs[i])
		btnPlayerNames = append(btnPlayerNames, result.Data(name, "player_names", dataString))
	}

	result.Inline(
		result.Row(btnPlayerNames...),
		result.Row(b.BtnCancel),
	)

	return result
}

func TitleQuestionKeyboard(bot *bot.BotData) *tele.ReplyMarkup {
	result := &tele.ReplyMarkup{}

	btnNo := result.Data("No", "no_title")
	btnYes := result.Data("Yes", "yes_title")

	result.Inline(
		result.Row(btnNo, btnYes),
		result.Row(bot.BtnCancel),
	)

	return result
}
