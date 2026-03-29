package test

import (
	"log"
	"os"
	"testing"

	"github.com/PromZona/AsyncMaster/internal/app"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	goose "github.com/pressly/goose/v3"
)

var (
	application *app.App
	rt          *runtime.MockRuntime
	err         error
)

func TestMain(m *testing.M) {
	application, err = app.InitTesting()
	if err != nil {
		log.Printf("Testing: Failed initialize application\n%s", err)
		os.Exit(1)
	}

	// Clean DB
	err := DropTablesDB(application.DB)
	if err != nil {
		log.Printf("Testing: Failed at dropping db\n%s", err)
		os.Exit(1)
	}

	// Goose migration
	err = goose.Up(application.DB, "./../../../migrations")
	if err != nil {
		log.Printf("Testing: Failed to apply migrations.\n%s", err)
		os.Exit(1)
	}

	rt = application.Runtime.(*runtime.MockRuntime)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func resetState() {
	err := TruncateAllTablesDB(application.DB)
	if err != nil {
		panic("Failed to truncate Database")
	}
	rt.UserManager.DeleteAllUsers()
	for k := range application.BotData.Sessions {
		delete(application.BotData.Sessions, k)
	}
}
