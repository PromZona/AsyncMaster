package router

import (
	"fmt"
	"log"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/handlers"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

var Routes = []*bot.Route{
	{
		Callback:              bot.HID_factions_list,
		Handler:               handlers.ListFactions,
		NextPossibleCallbacks: []bot.HandlerID{},
	},
	{
		Callback:              bot.HID_factions_update,
		Handler:               handlers.FactionsUpdate,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_factions_update_player},
	},
	{
		Callback:              bot.HID_factions_update_player,
		Handler:               handlers.FactionUpdatePickWhat,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_factions_update_resources}, // can add more options to update here
	},
	{
		Callback:              bot.HID_factions_update_resources,
		Handler:               handlers.FactionUpdateResources,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_factions_update_resources_text},
	},
	{
		Callback:              bot.HID_factions_update_resources_text,
		Handler:               handlers.FactionUpdateResourcesText,
		NextPossibleCallbacks: []bot.HandlerID{},
	},

	// SEND MESSAGE SCENARIO
	{
		Callback:              bot.HID_message_send,
		Handler:               handlers.InitialSendMessage,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_message_player_name},
	},
	{
		Callback:              bot.HID_message_send_everyone,
		Handler:               handlers.InitialSendEveryone,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_message_text},
	},
	{
		Callback:              bot.HID_message_player_name,
		Handler:               handlers.PlayerNameProcess,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_message_text},
	},
	{
		Callback:              bot.HID_message_text,
		Handler:               handlers.MessageTextProcess,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_message_title_yes, bot.HID_message_title_no},
	},
	{
		Callback:              bot.HID_message_title_yes,
		Handler:               handlers.MessageYesTitle,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_message_title_text},
	},
	{
		Callback:              bot.HID_message_title_no,
		Handler:               handlers.MessageNoTitle,
		NextPossibleCallbacks: []bot.HandlerID{},
	},
	{
		Callback:              bot.HID_message_title_text,
		Handler:               handlers.MessageTitleProcess,
		NextPossibleCallbacks: []bot.HandlerID{},
	},
	// END SEND MESSAGE SCENARIO

	{
		Callback:              bot.HID_message_list,
		Handler:               handlers.ListMessages,
		NextPossibleCallbacks: []bot.HandlerID{},
	},
	{
		Callback:              bot.HID_message_get,
		Handler:               handlers.GetMessage,
		NextPossibleCallbacks: []bot.HandlerID{},
	},

	// MASTER REQUEST SCENARIO
	{
		Callback:              bot.HID_mr_send,
		Handler:               handlers.InitialMasterRequest,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_mr_resipient},
	},
	{
		Callback:              bot.HID_mr_send_everyone,
		Handler:               handlers.InitialMasterRequestEveryone,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_mr_text},
	},
	{
		Callback:              bot.HID_mr_resipient,
		Handler:               handlers.MasterRequestResipient,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_mr_text},
	},
	{
		Callback:              bot.HID_mr_text,
		Handler:               handlers.MasterRequestText,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_mr_dice_yes, bot.HID_mr_dice_no},
	},
	{
		Callback:              bot.HID_mr_dice_yes,
		Handler:               handlers.MasterRequestDiceYes,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_mr_dice_text},
	},
	{
		Callback:              bot.HID_mr_dice_no,
		Handler:               handlers.MasterRequestDiceNo,
		NextPossibleCallbacks: []bot.HandlerID{},
	},
	{
		Callback:              bot.HID_mr_dice_text,
		Handler:               handlers.MasterRequestAddRoll,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_mr_dice_yes, bot.HID_mr_dice_no},
	},
	// END MASTER REQUEST SCENARIO

	// MASTER REQUEST ANSWER SCENARIO
	{
		Callback:              bot.HID_mr_answer,
		Handler:               handlers.MasterRequestAnswer,
		NextPossibleCallbacks: []bot.HandlerID{bot.HID_mr_answer_text},
	},
	{
		Callback:              bot.HID_mr_answer_text,
		Handler:               handlers.MasterRequestAnswerText,
		NextPossibleCallbacks: []bot.HandlerID{},
	},
	{
		Callback:              bot.HID_mr_answer_roll,
		Handler:               handlers.MasterRequestAnswerRoll,
		NextPossibleCallbacks: []bot.HandlerID{},
	},

	// END MASTER REQUEST ANSWER SCENARIO

	// MASTER REQUEST GET SCENARIO
	{
		Callback:              bot.HID_mr_player_get,
		Handler:               handlers.MasterRequestFirstUnanswered,
		NextPossibleCallbacks: []bot.HandlerID{},
	},
	{
		Callback:              bot.HID_mr_master_get,
		Handler:               handlers.MasterRequestGetFirstAnswered,
		NextPossibleCallbacks: []bot.HandlerID{},
	},
	{
		Callback:              bot.HID_mr_master_mark_read,
		Handler:               handlers.MasterRequestMarkAsRead,
		NextPossibleCallbacks: []bot.HandlerID{},
	},
	// END MASTER REQUEST GET SCENARIO
}

func DispatchText(context runtime.Context, b *bot.BotData) error {
	chatID := context.ChatID()

	session, ok := b.Sessions[chatID]
	if !ok {
		session = bot.NewSession(b.DB)
		b.Sessions[chatID] = session
		log.Printf("Created new session for %d", chatID)
	}

	prevRoute := session.PreviousRoute
	if prevRoute == nil || len(prevRoute.NextPossibleCallbacks) == 0 {
		user, err := db.GetUserByID(b.DB, chatID)
		if err != nil {
			return err
		}
		return handlers.GetMainMenuByRole(context, b.DB, user)
	}

	if len(prevRoute.NextPossibleCallbacks) > 1 {
		panic("Expected to be only one callback for text. If there are options, then what to choose...")
	}

	nextHID := prevRoute.NextPossibleCallbacks[0]

	var route *bot.Route
	for _, r := range Routes {
		if r.Callback == bot.HandlerID(nextHID) {
			route = r
		}
	}

	if route == nil {
		panic("Routes dependency is wrongly set. Expected Route from available routes for given HID")
	}

	log.Printf("Received Text from %d. Going to process by %s route", chatID, route.Callback)
	err := route.Handler(context, session)
	if err != nil {
		return err
	}

	session.PreviousRoute = route
	return nil
}

func DispatchCallback(context runtime.Context, b *bot.BotData) error {
	context.Respond()

	chatID := context.ChatID()
	rawCallbackData := context.Callback()
	cbUnique, _ := handlers.ParseCallbackDataString(rawCallbackData)

	if bot.HandlerID(cbUnique) == bot.HID_cancel {
		return handlers.HandleCancelButton(context, b)
	}

	session, ok := b.Sessions[chatID]
	if !ok {
		session = bot.NewSession(b.DB)
		b.Sessions[chatID] = session
		log.Printf("Created new session for %d", chatID)
	}

	var route *bot.Route
	for _, r := range Routes {
		if r.Callback == bot.HandlerID(cbUnique) {
			route = r
		}
	}

	if route == nil {
		return fmt.Errorf("met unsupported callback %s", cbUnique)
	}

	if session.PreviousRoute != nil && len(session.PreviousRoute.NextPossibleCallbacks) != 0 {
		var isAllowedCB bool
		for _, allowedCB := range session.PreviousRoute.NextPossibleCallbacks {
			if allowedCB == bot.HandlerID(cbUnique) {
				isAllowedCB = true
				break
			}
		}
		if !isAllowedCB {
			return context.Send("Эта руна сейчас не доступна. Закончи делать что делал - и вернись к этой руне позже..")
		}
	}

	log.Printf("Received Callback from %d. Going to process by %s route", chatID, route.Callback)
	err := route.Handler(context, session)
	if err != nil {
		return err
	}

	session.PreviousRoute = route
	return nil
}
