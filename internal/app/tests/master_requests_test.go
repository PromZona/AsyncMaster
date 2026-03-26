package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func TestMasterRequest_Create_HappyPath(t *testing.T) {
	resetState()

	//
	// SETUP: 1 master, 1 player
	//
	input := []string{
		`/server create_user Master MasterTG`,
		`/server create_user John JohnTG`,

		// Master
		`/user Master "MasterPega26"`,
		`/user Master "GameMaster"`,

		// Player
		`/user John "pepega26"`,
		`/user John "John"`,
		`/user John "FactionJ"`,
		`/user John "DescJ"`,
	}

	for _, cmd := range input {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// ACTION: Master creates request
	//
	actions := []string{
		`/user Master master_request_send`,
		`/user Master master_request_resipient|John:1`,
		`/user Master "Do something important"`,
		`/user Master master_request_dice_no`,
	}

	for _, cmd := range actions {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// VERIFY DB
	//
	mr, err := db.GetFirstUnansweredMasterRequest(application.DB, 1)
	if err != nil {
		t.Fatal(err)
	}

	if mr.TextRequest != "Do something important" {
		t.Fatal("unexpected master request text")
	}
}

func TestMasterRequest_Create_And_Answer(t *testing.T) {
	resetState()

	//
	// SETUP
	//
	input := []string{
		`/server create_user Master MasterTG`,
		`/server create_user John JohnTG`,

		// Master
		`/user Master "MasterPega26"`,
		`/user Master "GameMaster"`,

		// Player
		`/user John "pepega26"`,
		`/user John "John"`,
		`/user John "FactionJ"`,
		`/user John "DescJ"`,
	}

	for _, cmd := range input {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// CREATE REQUEST
	//
	create := []string{
		`/user Master master_request_send`,
		`/user Master master_request_resipient|John:1`,
		`/user Master "Attack enemy base"`,
		`/user Master master_request_dice_no`,
	}

	for _, cmd := range create {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// GET request ID
	//
	mr, err := db.GetFirstUnansweredMasterRequest(application.DB, 1)
	if err != nil {
		t.Fatal(err)
	}

	//
	// PLAYER ANSWERS
	//
	answerCmd := fmt.Sprintf(`/user John master_request_answer|%d`, mr.ID)
	actions := []string{
		answerCmd,
		`/user John "We attack at dawn"`,
	}

	for _, cmd := range actions {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// VERIFY DB
	//
	updated, err := db.GetMasterRequestByID(application.DB, mr.ID)
	if err != nil {
		t.Fatal(err)
	}

	if updated.TextResponse != "We attack at dawn" {
		t.Fatal("response not saved")
	}

	if updated.State != bot.MRAnswered {
		t.Fatal("state not updated to answered")
	}
}

func TestMasterRequest_FullFlow(t *testing.T) {
	resetState()

	//
	// SETUP
	//
	input := []string{
		`/server create_user Master MasterTG`,
		`/server create_user John JohnTG`,

		// Master
		`/user Master "MasterPega26"`,
		`/user Master "GameMaster"`,

		// Player
		`/user John "pepega26"`,
		`/user John "John"`,
		`/user John "FactionJ"`,
		`/user John "DescJ"`,
	}

	for _, cmd := range input {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// CREATE
	//
	create := []string{
		`/user Master master_request_send`,
		`/user Master master_request_resipient|John:1`,
		`/user Master "Scout area"`,
		`/user Master master_request_dice_no`,
	}

	for _, cmd := range create {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	mr, err := db.GetFirstUnansweredMasterRequest(application.DB, 1)
	if err != nil {
		t.Fatal(err)
	}

	//
	// ANSWER
	//
	answerCmd := fmt.Sprintf(`/user John master_request_answer|%d`, mr.ID)
	answer := []string{
		answerCmd,
		`/user John "Area is clear"`,
	}

	for _, cmd := range answer {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// MASTER GET ANSWERED
	//
	err, _ = runtime.ExecuteCommand(rt, `/user Master master_request_master_get`)
	if err != nil {
		t.Fatal(err)
	}

	last := rt.MessageManager.Messages[len(rt.MessageManager.Messages)-1]

	if !strings.Contains(last.Text, "Area is clear") {
		t.Fatal("master did not receive correct response")
	}

	//
	// FIND mark_read callback
	//
	var callback string
	found := false

	for _, row := range last.Keyboard {
		for _, btn := range row {
			if strings.Contains(btn.Unique, "master_request_master_mark_read") {
				callback = btn.Unique + "|" + btn.Data
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Fatal("mark read button not found")
	}

	//
	// MARK AS READ
	//
	cmd := fmt.Sprintf(`/user Master %s`, callback)
	err, _ = runtime.ExecuteCommand(rt, cmd)
	if err != nil {
		t.Fatal(err)
	}

	//
	// VERIFY DB STATE
	//
	final, err := db.GetMasterRequestByID(application.DB, mr.ID)
	if err != nil {
		t.Fatal(err)
	}

	if final.State != bot.MRCheckedByMaster {
		t.Fatal("master request not marked as read")
	}
}
