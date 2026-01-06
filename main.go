package main

import "github.com/PromZona/AsyncMaster/internal/app"

func main() {
	application, err := app.Init()
	if err != nil {
		panic(err)
	}
	application.Start()
}
