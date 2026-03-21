package handlers

import (
	"fmt"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func ListFactions(context runtime.Context, s *bot.Session) error {

	data, err := db.GetUsersAll(s.DB)
	if err != nil {
		return err
	}

	var text strings.Builder
	text.WriteString("Factions\n\n")

	for _, user := range data {
		faction := fmt.Sprintf(`
			Faction Name: %s
			Leader: %s
			Description: %s

			---
			`,
			user.Faction.Name,
			user.PlayerName,
			user.Faction.Description,
		)

		text.WriteString(faction)
	}

	return context.Send(text.String())
}
