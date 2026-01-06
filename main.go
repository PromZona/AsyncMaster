package main

import "github.com/PromZona/AsyncMaster/internal/app"

func main() {
	b, db, err := app.Init()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app.Start(b)
}
