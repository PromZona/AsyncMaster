package router

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/flows/master"
	"github.com/PromZona/AsyncMaster/internal/app/flows/registration"
	"github.com/PromZona/AsyncMaster/internal/app/middleware"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func Register(rt runtime.Runtime, botData *bot.BotData) {

	rt.Use(middleware.ErrorRecovery(botData))
	rt.Use(middleware.RegistrationCheck(botData))

	rt.HandleCommand("/start", func(ctx runtime.Context) error { return registration.StartMessage(ctx) })
	rt.HandleCommand("/elevate", func(ctx runtime.Context) error { return master.HandleElevateToMaster(ctx, botData) })

	rt.HandleText(func(ctx runtime.Context) error { return DispatchText(ctx, botData) })
	rt.HandleCallback(func(ctx runtime.Context) error { return DispatchCallback(ctx, botData) })
}
