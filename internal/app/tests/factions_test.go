package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func TestFactionCreation(t *testing.T) {

	//
	// TEST: CLEANING
	//
	resetState()

	//
	// TEST: EXECUTING LOGIC
	//
	input := []string{
		`/server create_user John Vitus`,
		`/server create_user Victor Horrow`,
		`/user John "pepega26"`,
		`/user John "Vitus"`,
		`/user John "Tigers"`,
		`/user John "Stripe tigers. Very powerful"`,
		`/user Victor "pepega26"`,
		`/user Victor "Horrow"`,
		`/user Victor "Necro Guild"`,
		`/user Victor "Super scary necromancers. And a lot of money"`,
	}
	for _, cmd := range input {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// TEST: EVALUATION
	//
	john, err := db.GetUserByID(application.DB, 0)
	if err != nil {
		t.Fatal(err)
	}

	victor, err := db.GetUserByID(application.DB, 1)
	if err != nil {
		t.Fatal(err)
	}

	if john.Faction.Name != "Tigers" {
		t.Fatal("Met unexpected faction name")
	}

	if victor.Faction.Name != "Necro Guild" {
		t.Fatal("Met unexpected faction name")
	}
}

func TestFactionGet(t *testing.T) {

	//
	// TEST: CLEANING BEFORE
	//
	resetState()

	//
	// TEST: EXECUTING LOGIC
	//
	input := []string{
		`/server create_user John Vitus`,
		`/server create_user Victor Horrow`,
		`/user John "pepega26"`,
		`/user John "Vitus"`,
		`/user John "Tigers"`,
		`/user John "Stripe tigers. Very powerful"`,
		`/user Victor "pepega26"`,
		`/user Victor "Horrow"`,
		`/user Victor "Necro Guild"`,
		`/user Victor "Super scary necromancers. And a lot of money"`,
		`/user John factions_list`,
	}
	for _, cmd := range input {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// TEST: EVALUATION
	//
	last := rt.MessageManager.Messages[len(rt.MessageManager.Messages)-1]
	if !strings.Contains(last.Text, "Necro Guild") {
		t.Fatal("failed to find guild in a message")
	}
}

func TestFaction_Update_Resources_HappyPath(t *testing.T) {
	resetState()

	//
	// SETUP: 1 master + 1 player with faction
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
		`/user John "Initial Resources"`,
	}

	for _, cmd := range input {
		err, _ = runtime.ExecuteCommand(rt, cmd)
		if err != nil {
			t.Fatal(err)
		}
	}

	//
	// ACTION: start faction update flow
	//
	err, _ = runtime.ExecuteCommand(rt, `/user Master factions_update`)
	if err != nil {
		t.Fatal(err)
	}

	//
	// PICK USER (John)
	//
	last := rt.MessageManager.Messages[len(rt.MessageManager.Messages)-1]

	var pickCallback string
	found := false

	for _, row := range last.Keyboard {
		for _, btn := range row {
			if strings.Contains(btn.Unique, "factions_update_player") &&
				strings.Contains(btn.Data, "1") { // John chatID = 1
				pickCallback = btn.Unique + "|" + btn.Data
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Fatal("factions_update_player button not found")
	}

	cmd := fmt.Sprintf(`/user Master %s`, pickCallback)
	err, _ = runtime.ExecuteCommand(rt, cmd)
	if err != nil {
		t.Fatal(err)
	}

	//
	// PICK "update resources"
	//
	last = rt.MessageManager.Messages[len(rt.MessageManager.Messages)-1]

	var resourcesCallback string
	found = false

	for _, row := range last.Keyboard {
		for _, btn := range row {
			if strings.Contains(btn.Unique, "factions_update_resources") {
				resourcesCallback = btn.Unique
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Fatal("factions_update_resources button not found")
	}

	cmd = fmt.Sprintf(`/user Master %s`, resourcesCallback)
	err, _ = runtime.ExecuteCommand(rt, cmd)
	if err != nil {
		t.Fatal(err)
	}

	//
	// INPUT NEW RESOURCES
	//
	newResources := "Gold: 1000, Army: 500"
	cmd = `/user Master "` + newResources + `"`
	err, _ = runtime.ExecuteCommand(rt, cmd)
	if err != nil {
		t.Fatal(err)
	}

	//
	// VERIFY RESPONSE
	//
	last = rt.MessageManager.Messages[len(rt.MessageManager.Messages)-1]
	if !strings.Contains(last.Text, "Updated successfuly") {
		t.Fatal("expected success message")
	}

	//
	// VERIFY DB
	//
	user, err := db.GetUserByID(application.DB, 1)
	if err != nil {
		t.Fatal(err)
	}

	if user.Faction.Resources != newResources {
		t.Fatalf("resources not updated, got: %s", user.Faction.Resources)
	}
}
