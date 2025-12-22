package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v4"

	"github.com/PromZona/AsyncMaster/internal/app"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	psqlInfo := os.Getenv("DB_CONNECTION_STRING")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Print("Database successfully connected!")

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	botData := app.BotData{DB: db, Users: make(map[int64]*app.UserData)}

	b.Handle("/start", func(ctx tele.Context) error { return app.HandleStartMessage(ctx, &botData) })
	b.Handle(tele.OnText, func(ctx tele.Context) error { return app.HandleText(ctx, &botData) })
	b.Handle("/sendAll", func(ctx tele.Context) error { return app.MasterSendMessageToAll(ctx, &botData) })

	b.Start()
}
