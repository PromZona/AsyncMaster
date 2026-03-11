package app

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/router"
	"github.com/PromZona/AsyncMaster/internal/app/runtime"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

type App struct {
	Runtime runtime.Runtime
	BotData *bot.BotData
	DB      *sql.DB
}

func Init() (*App, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
		return nil, err
	}

	isMock := flag.Bool("mock", false, "mocking runtime")
	flag.Parse()

	psqlInfo := os.Getenv("DB_CONNECTION_STRING")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("failed to open db connection")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping the server")
		return nil, err
	}
	log.Print("Database successfully connected!")

	var rt runtime.Runtime
	if *isMock {
		panic("Not yet")
	} else {
		pref := tele.Settings{
			Token:  os.Getenv("BOT_TOKEN"),
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}

		b, err := tele.NewBot(pref)
		if err != nil {
			log.Fatal("failed to create tele.Bot")
			return nil, err
		}

		rt = &runtime.TelebotRuntime{
			Bot: b,
		}
	}

	botData := bot.BotInit(db)
	router.Register(rt, botData)

	app := &App{
		Runtime: rt,
		BotData: botData,
		DB:      db,
	}

	return app, nil
}

func (app *App) Start() {
	defer app.DB.Close()
	app.Runtime.Start()
}
