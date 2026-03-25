package test

import (
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
		`/user John list_factions`,
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
