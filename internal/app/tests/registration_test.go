package test

import (
	"testing"

	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func TestRegistration_CannotSkipSteps(t *testing.T) {
	//
	// TEST: CLEANING
	//
	resetState()

	//
	// TEST: EXECUTING LOGIC
	//
	input := []string{
		`/server create_user John Vitus`,
		// trying to jump directly to faction
		`/user John "Wolves"`,
	}

	for _, cmd := range input {
		_, _ = runtime.ExecuteCommand(rt, cmd)
	}

	//
	// TEST: Evaluationg
	//
	_, err = db.GetUserByID(application.DB, 0)
	if err == nil {
		t.Fatal("expected user not registering")
	}
}

func TestRegistration_UnexpectedCommand(t *testing.T) {

	//
	// TEST: CLEANING
	//
	resetState()

	//
	// TEST: EXECUTING LOGIC
	//
	input := []string{
		`/server create_user John Vitus`,

		`/user John "pepega26"`,
		`/user John list_factions`, // unexpected in middle
		`/user John "Vitus"`,
		`/user John "Wolves"`,
		`/user John "Strong wolves"`,
	}

	for _, cmd := range input {
		_, _ = runtime.ExecuteCommand(rt, cmd)
	}

	//
	// TEST: EVALUATION
	//
	john, err := db.GetUserByID(application.DB, 0)
	if err != nil {
		t.Fatal(err)
	}

	if john.Faction.Name != "Wolves" {
		t.Fatal("flow broken by unexpected command")
	}
}
