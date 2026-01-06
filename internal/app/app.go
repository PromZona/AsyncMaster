package app

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/PromZona/AsyncMaster/internal/app/bot"
	"github.com/PromZona/AsyncMaster/internal/app/router"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

type App struct {
	TeleBot *tele.Bot
	BotData *bot.BotData
	DB      *sql.DB
}

func Init() (*App, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
		return nil, err
	}

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

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal("failed to create tele.Bot")
		return nil, err
	}

	botData := bot.BotInit(db)
	router.Register(b, botData)

	app := &App{
		TeleBot: b,
		BotData: botData,
		DB:      db,
	}

	return app, nil
}

func (app *App) Start() {
	defer app.DB.Close()
	app.TeleBot.Start()
}
