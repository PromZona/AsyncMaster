package test

import (
	"fmt"
	"os"
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
	fmt.Printf("main")

	application, err = app.InitTesting()
	if err != nil {
		fmt.Printf("Testing: Failed initialize application\n%s", err)
		os.Exit(1)
	}

	// Clean DB
	_, err := application.DB.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;`)
	if err != nil {
		fmt.Printf("Testing: Failed at dropping db\n%s", err)
		os.Exit(1)
	}

	// Goose migration
	err = goose.Up(application.DB, "./../../../migrations")
	if err != nil {
		fmt.Printf("Testing: failed to apply migrations.\n%s", err)
		os.Exit(1)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestFactionCreation(t *testing.T) {
	rt := application.Runtime.(*runtime.MockRuntime)

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
