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

func Init() (*tele.Bot, *sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
		return nil, nil, err
	}

	psqlInfo := os.Getenv("DB_CONNECTION_STRING")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("failed to open db connection")
		return nil, nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping the server")
		return nil, nil, err
	}
	log.Print("Database successfully connected!")

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal("failed to create tele.Bot")
		return nil, nil, err
	}

	botData := bot.BotInit(db)
	router.Register(b, botData)

	return b, db, nil
}

func Start(b *tele.Bot) {
	b.Start()
}
