package router

import (
	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/handlers"
	"github.com/PromZona/AsyncMaster/internal/app/middleware"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
)

func Register(rt runtime.Runtime, botData *bot.BotData) {

	rt.Use(middleware.ErrorRecovery(botData))
	rt.Use(middleware.RegistrationCheck(botData))

	rt.HandleCommand("/elevate", func(ctx runtime.Context) error { return handlers.HandleElevateToMaster(ctx, botData) })

	rt.HandleText(func(ctx runtime.Context) error { return DispatchText(ctx, botData) })
	rt.HandleCallback(func(ctx runtime.Context) error { return DispatchCallback(ctx, botData) })
}
