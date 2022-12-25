package main

import (
	"database/sql"
	"funny-bot/internal/app"
	"funny-bot/internal/database"
	"funny-bot/internal/handler"
	"funny-bot/internal/scheduler"
	"funny-bot/internal/telegram"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"time"
)

type config struct {
	botToken        string
	databaseConnUrl string
}

func readConfig() config {
	return config{
		botToken:        os.Getenv("TELEGRAM_BOT_TOKEN"),
		databaseConnUrl: os.Getenv("DB_CONNECTION_URL"),
	}
}

func initDatabase(databaseConnUrl string) (*sql.DB, error) {
	db, err := database.OpenConnection(databaseConnUrl)
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return nil, err
	}
	err = m.Up()
	log.Printf("Database migration: %s\n", err)
	return db, nil
}

func main() {
	cfg := readConfig()
	bot, err := telegram.NewBot(cfg.botToken)
	if err != nil {
		log.Fatalln(err)
	}
	db, err := initDatabase(cfg.databaseConnUrl)
	if err != nil {
		log.Fatalln(err)
	}
	repository := database.NewUserRepository(db)
	fs := app.NewFunnyService(repository, bot)
	ns := scheduler.NewNotificationScheduler(fs, time.Second)

	bot.AddHandler(handler.NewStartHandler(repository))
	ns.Start()
	bot.Start()
}
