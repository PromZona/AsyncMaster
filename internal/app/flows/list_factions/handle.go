package list_faction

import (
	"fmt"
	"strings"

	"github.com/PromZona/AsyncMaster/internal/app/db"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func handleListFactions(context runtime.Context, s *Session) error {

	data, err := db.GetUsersAll(s.DB)
	if err != nil {
		return err
	}

	var text strings.Builder

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

	s.Done = true
	return context.Send(text.String())
}
