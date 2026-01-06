package router

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/dispatcher"
	"github.com/PromZona/AsyncMaster/internal/app/flows/master"
	"github.com/PromZona/AsyncMaster/internal/app/flows/registration"
	"github.com/PromZona/AsyncMaster/internal/app/middleware"
	"github.com/PromZona/AsyncMaster/internal/app/ui"
	tele "gopkg.in/telebot.v4"
)

func Register(b *tele.Bot, botData *bot.BotData) {

	b.Use(middleware.ErrorRecovery(botData))
	b.Use(middleware.RegistrationCheck(botData))

	b.Handle("/start", func(ctx tele.Context) error { return registration.StartMessage(ctx) })
	b.Handle("/elevate", func(ctx tele.Context) error { return master.HandleElevateToMaster(ctx, botData) })

	b.Handle(tele.OnText, func(ctx tele.Context) error { return dispatcher.Text(ctx, botData) })
	b.Handle(tele.OnCallback, func(ctx tele.Context) error { return dispatcher.Callbacks(ctx, botData) })

	b.Handle(&botData.BtnCancel, func(ctx tele.Context) error { return ui.HandleCancelButton(ctx, botData) })
}
