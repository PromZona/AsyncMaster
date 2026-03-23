package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/PromZona/AsyncMaster/internal/app"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"

	goose "github.com/pressly/goose/v3"
)

var (
	application *app.App
	err         error
)

func TestMain(m *testing.M) {
	application, err = app.InitTesting()
	if err != nil {
		fmt.Printf("Testing: Failed initialize application\n%s", err)
		os.Exit(1)
	}

	// Clean DB
	err := DropTablesDB(application.DB)
	if err != nil {
		fmt.Printf("Testing: Failed at dropping db\n%s", err)
		os.Exit(1)
	}

	// Goose migration
	err = goose.Up(application.DB, "./../../../migrations")
	if err != nil {
		fmt.Printf("Testing: Failed to apply migrations.\n%s", err)
		os.Exit(1)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestFactionCreation(t *testing.T) {

	//
	// TEST: CLEANING BEFORE
	//
	err := TruncateAllTablesDB(application.DB)
	if err != nil {
		t.Fatal("Can not clean DB and start testing")
	}
	rt := application.Runtime.(*runtime.MockRuntime)
	rt.UserManager.DeleteAllUsers()
	for k := range application.BotData.Sessions {
		delete(application.BotData.Sessions, k)
	}

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
	err := TruncateAllTablesDB(application.DB)
	if err != nil {
		t.Fatal("Can not clean DB and start testing")
	}
	rt := application.Runtime.(*runtime.MockRuntime)
	rt.UserManager.DeleteAllUsers()
	for k := range application.BotData.Sessions {
		delete(application.BotData.Sessions, k)
	}

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
