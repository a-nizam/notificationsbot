package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"main/helper"
	"time"
)

type tgBot struct {
	bot      *tgbotapi.BotAPI
	keyboard tgbotapi.InlineKeyboardMarkup
	action   Action
}

type Action struct {
	chatId int64
	action string
}

func (bot *tgBot) Init() {
	sendMessage := func(id int64, text string) {
		msg := tgbotapi.NewMessage(id, text)
		_, err := bot.bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
	clearAction := func() {
		bot.action.chatId = 0
		bot.action.action = ""
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
		if update.Message != nil && update.Message.IsCommand() {
			switch update.Message.Command() {
			case "menu":
				// Создаем клавиатуру
				bot.keyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Меню напоминаний", "menu"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Создать", "create"),
						tgbotapi.NewInlineKeyboardButtonData("Список", "list"),
						tgbotapi.NewInlineKeyboardButtonData("Удалить", "remove"),
					),
				)

				// Отправляем сообщение с клавиатурой
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите пункт меню:")
				msg.ReplyMarkup = bot.keyboard
				_, err := bot.bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			switch data {
			case "create":
				bot.action = Action{update.CallbackQuery.Message.Chat.ID, "create"}
				sendMessage(update.CallbackQuery.Message.Chat.ID, "Введите текст напоминания")
			}
		} else {
			if bot.action.action == "create" && bot.action.chatId == update.Message.Chat.ID {
				clearAction()
				notificationDb := helper.NotificationsDB{}
				err = notificationDb.Init("data.db")
				if err != nil {
					log.Fatal(err)
				}
				id, err := notificationDb.InsertChat(helper.Chat{ChatId: update.Message.Chat.ID})
				if err != nil {
					log.Fatal(err)
				}
				_, err = notificationDb.InsertNotification(helper.Notification{Notification: update.Message.Text, Datetime: time.Now().String(), ChatId: id})
				if err != nil {
					log.Fatal(err)
				}
				sendMessage(update.Message.Chat.ID, "Напоминание успешно добавлено")
			} else {
				sendMessage(update.Message.Chat.ID, "Введите \"/menu\" для вызова меню бота")
			}
		}
	}
}

func main() {
	myTgBot := &tgBot{}
	myTgBot.Init()
}
