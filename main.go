package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"main/helper"
	"strings"
	"time"
)

type tgBot struct {
	bot *tgbotapi.BotAPI
}

func (bot *tgBot) Init() {
	sendMessage := func(id int64, text string) {
		msg := tgbotapi.NewMessage(id, text)
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	var err error
	bot.bot, err = tgbotapi.NewBotAPI("243190873:AAGQZxq2Vttvdbo2egkFjNVIqxEgYaDvf-Y")
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
		words := strings.Split(update.Message.Text, " ")
		command := words[0]
		text := strings.Join(words[1:], " ")
		switch command {
		case "напомнить":
			notificationDb := helper.NotificationsDB{}
			err = notificationDb.Init("data.db")
			if err != nil {
				log.Fatal(err)
			}
			id, err := notificationDb.InsertChat(helper.Chat{ChatId: update.Message.Chat.ID})
			if err != nil {
				log.Fatal(err)
			}
			_, err = notificationDb.InsertNotification(helper.Notification{Notification: text, Datetime: time.Now().String(), ChatId: id})
			if err != nil {
				log.Fatal(err)
			}
			sendMessage(update.Message.Chat.ID, "Напоминание успешно добавлено")
		default:
			sendMessage(update.Message.Chat.ID, "Введите \"напомнить текст\"")
		}
	}
}

func main() {
	myTgBot := &tgBot{}
	myTgBot.Init()
}
