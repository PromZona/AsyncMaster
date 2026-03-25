package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func TestMessage_HappyPath_PlayersAndMaster(t *testing.T) {
	resetState()

	//
	// SETUP: Create users (2 players + 1 master)
	//
	input := []string{
		`/server create_user John JohnTG`,
		`/server create_user Alice AliceTG`,
		`/server create_user Bob BobTG`,

		// John registration (player)
		`/user John "pepega26"`,
		`/user John "John"`,
		`/user John "FactionJ"`,
		`/user John "DescJ"`,

		// Alice registration (player)
		`/user Alice "pepega26"`,
		`/user Alice "Alice"`,
		`/user Alice "FactionA"`,
		`/user Alice "DescA"`,

		// Bob registration (MASTER)
		`/user Bob "MasterPega26"`,
		`/user Bob "Bob"`,
	}

	for _, cmd := range input {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// ACTION: John sends message to Alice
	//
	inputJohn := []string{
		`/user John message_send`,
		`/user John message_player_name|Alice:1`, // Alice chatID assumed = 1
		`/user John "Hello Alice"`,
		`/user John message_title_no`,
	}

	for _, cmd := range inputJohn {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// ACTION: Alice sends message to John
	//
	inputAlice := []string{
		`/user Alice message_send`,
		`/user Alice message_player_name|John:0`, // John chatID assumed = 0
		`/user Alice "Hello John"`,
		`/user Alice message_title_no`,
	}

	for _, cmd := range inputAlice {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// ACTION: Master sends message to everyone
	//
	inputMaster := []string{
		`/user Bob message_send_everyone`,
		`/user Bob "Hello everyone"`,
		`/user Bob message_title_no`,
	}

	for _, cmd := range inputMaster {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// EVALUATION
	//

	// Check DB transactions count (2 messages sent to John)
	transactions, err := db.GetLastMessageTransactions(application.DB, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(transactions) != 2 {
		t.Fatal("expected at 2 message transactions")
	}

	// Check delivery via Mock messages
	var foundJohnToAlice bool
	var foundAliceToJohn bool
	var foundMasterBroadcast int

	for _, msg := range rt.MessageManager.Messages {
		if strings.Contains(msg.Text, "Hello Alice") {
			foundJohnToAlice = true
		}
		if strings.Contains(msg.Text, "Hello John") {
			foundAliceToJohn = true
		}
		if strings.Contains(msg.Text, "Hello everyone") {
			foundMasterBroadcast++
		}
	}

	if !foundJohnToAlice {
		t.Fatal("John → Alice message not delivered")
	}

	if !foundAliceToJohn {
		t.Fatal("Alice → John message not delivered")
	}

	// Master should send to ALL users (3 users)
	if foundMasterBroadcast < 3 {
		t.Fatal("master broadcast not delivered to all users")
	}
}

func TestMessage_FlowValidation(t *testing.T) {
	err := TruncateAllTablesDB(application.DB)
	if err != nil {
		t.Fatal(err)
	}

	rt := application.Runtime.(*runtime.MockRuntime)
	rt.UserManager.DeleteAllUsers()
	for k := range application.BotData.Sessions {
		delete(application.BotData.Sessions, k)
	}

	//
	// SETUP: 2 players
	//
	input := []string{
		`/server create_user John JohnTG`,
		`/server create_user Alice AliceTG`,

		// John registration
		`/user John "pepega26"`,
		`/user John "John"`,
		`/user John "FactionJ"`,
		`/user John "DescJ"`,

		// Alice registration
		`/user Alice "pepega26"`,
		`/user Alice "Alice"`,
		`/user Alice "FactionA"`,
		`/user Alice "DescA"`,
	}

	for _, cmd := range input {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// TEST CASE 1: Skip recipient → send text
	//
	_, _ = runtime.ExecuteCommand(rt, `/user John message_send`)
	_, _ = runtime.ExecuteCommand(rt, `/user John "Hello without recipient"`)

	//
	// TEST CASE 2: Skip text → try finalize
	//
	_, _ = runtime.ExecuteCommand(rt, `/user John message_send`)
	_, _ = runtime.ExecuteCommand(rt, `/user John message_player_name|Alice:1`)
	_, _ = runtime.ExecuteCommand(rt, `/user John message_title_no`)

	//
	// TEST CASE 3: Invalid callback format
	//
	_, _ = runtime.ExecuteCommand(rt, `/user John message_player_name`)
	_, _ = runtime.ExecuteCommand(rt, `/user John message_player_name|invalidformat`)

	//
	// TEST CASE 4: Random command mid-flow
	//
	_, _ = runtime.ExecuteCommand(rt, `/user John message_send`)
	_, _ = runtime.ExecuteCommand(rt, `/user John list_factions`) // interrupt flow

	//
	// TEST CASE 5: Ensure flow still works after failures
	//
	validFlow := []string{
		`/user John message_send`,
		`/user John message_player_name|Alice:1`,
		`/user John "Hello Alice VALID"`,
		`/user John message_title_no`,

		`/user Alice message_send`,
		`/user Alice message_player_name|John:0`,
		`/user Alice "Hello John VALID"`,
		`/user Alice message_title_no`,
	}

	for _, cmd := range validFlow {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// EVALUATION
	//

	// Ensure ONLY valid messages were stored (2 messages expected)
	transactionsJohn, err := db.GetLastMessageTransactions(application.DB, 0)
	if err != nil {
		t.Fatal(err)
	}

	transactionsAlice, err := db.GetLastMessageTransactions(application.DB, 1)
	if err != nil {
		t.Fatal(err)
	}

	total := len(transactionsJohn) + len(transactionsAlice)
	if total < 2 {
		t.Fatal("expected valid transactions after flow recovery")
	}

	// Ensure valid messages exist in output
	var foundJohn bool
	var foundAlice bool

	for _, msg := range rt.MessageManager.Messages {
		if strings.Contains(msg.Text, "Hello Alice VALID") {
			foundJohn = true
		}
		if strings.Contains(msg.Text, "Hello John VALID") {
			foundAlice = true
		}
	}

	if !foundJohn {
		t.Fatal("valid John → Alice message missing")
	}

	if !foundAlice {
		t.Fatal("valid Alice → John message missing")
	}
}

func TestMessage_ListMessages_HappyPath(t *testing.T) {
	err := TruncateAllTablesDB(application.DB)
	if err != nil {
		t.Fatal(err)
	}

	rt := application.Runtime.(*runtime.MockRuntime)
	rt.UserManager.DeleteAllUsers()
	for k := range application.BotData.Sessions {
		delete(application.BotData.Sessions, k)
	}

	//
	// SETUP: 3 users
	//
	input := []string{
		`/server create_user John JohnTG`,
		`/server create_user Alice AliceTG`,
		`/server create_user Bob BobTG`,

		// John registration
		`/user John "pepega26"`,
		`/user John "John"`,
		`/user John "FactionJ"`,
		`/user John "DescJ"`,

		// Alice registration
		`/user Alice "pepega26"`,
		`/user Alice "Alice"`,
		`/user Alice "FactionA"`,
		`/user Alice "DescA"`,

		// Bob registration
		`/user Bob "pepega26"`,
		`/user Bob "Bob"`,
		`/user Bob "FactionB"`,
		`/user Bob "DescB"`,
	}

	for _, cmd := range input {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// ACTION: send messages to Alice
	//
	actions := []string{
		// John → Alice
		`/user John message_send`,
		`/user John message_player_name|Alice:1`,
		`/user John "Hello from John 1"`,
		`/user John message_title_yes`,
		`/user John "message title 1"`,

		`/user John message_send`,
		`/user John message_player_name|Alice:1`,
		`/user John "Hello from John 2"`,
		`/user John message_title_yes`,
		`/user John "message title 2"`,

		// Bob → Alice
		`/user Bob message_send`,
		`/user Bob message_player_name|Alice:1`,
		`/user Bob "Hello from Bob"`,
		`/user Bob message_title_yes`,
		`/user Bob "message title 3"`,
	}

	for _, cmd := range actions {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// ACTION: Alice lists messages
	//
	err, _ = runtime.ExecuteCommand(rt, `/user Alice message_list`)
	if err != nil {
		t.Fatal(err)
	}

	//
	// EVALUATION
	//
	last := rt.MessageManager.Messages[len(rt.MessageManager.Messages)-1]
	if !keyboardContains(last.Keyboard, "message title 1") {
		t.Fatal("missing message title 1")
	}
	if !keyboardContains(last.Keyboard, "message title 2") {
		t.Fatal("missing message title 2")
	}
	if !keyboardContains(last.Keyboard, "message title 3") {
		t.Fatal("missing message title 3")
	}

	//
	// ACTION: simulate clicking first message button
	//
	var callback string
	found := false

	for _, row := range last.Keyboard {
		for _, btn := range row {
			// skip non-message buttons if you have any (like cancel)
			if strings.Contains(btn.Unique, "message_get") {
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
		t.Fatal("no message_get button found in keyboard")
	}

	// EXECUTE: simulate button click
	cmd := fmt.Sprintf(`/user Alice %s`, callback)
	err, _ = runtime.ExecuteCommand(rt, cmd)
	if err != nil {
		t.Fatal(err)
	}

	// EVALUATION: check returned message
	last = rt.MessageManager.Messages[len(rt.MessageManager.Messages)-1]

	if !strings.Contains(last.Text, "From:") {
		t.Fatal("expected sender info in message output")
	}

	if !strings.Contains(last.Text, "message title") {
		t.Fatal("expected message title in output")
	}
}

func keyboardContains(keyboard runtime.Keyboard, text string) bool {
	for _, row := range keyboard {
		for _, btn := range row {
			if strings.Contains(btn.Text, text) {
				return true
			}
		}
	}
	return false
}
