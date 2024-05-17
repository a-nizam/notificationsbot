package main

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type notificationsDB struct {
	filename string
	conn     *sql.DB
}

type tgbot struct {
	bot *tgbotapi.BotAPI
}

func (db *notificationsDB) Init(filename string) error {
	db.filename = filename
	var err error
	db.conn, err = sql.Open("sqlite3", db.filename)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (bot *tgbot) Init() {
	var err error
	bot.bot, err = tgbotapi.NewBotAPI("6431929745:AAFcT3VHztJ5gSK5CXTjFTlB1x3H1UlQAC0")
	if err != nil {
		log.Fatal(err)
	}
	bot.bot.Debug = true

	u := tgbotapi.NewUpdate(60)
	u.Timeout = 60
	updates := bot.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		//update.Message.Chat.ID
	}
}

func main() {
	myTgbot := &tgbot{}
	myTgbot.Init()
}
