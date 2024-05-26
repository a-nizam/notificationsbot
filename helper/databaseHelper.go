package helper

import (
	"database/sql"
	"errors"
	"log"
)

type NotificationsDB struct {
	filename string
	conn     *sql.DB
}

type Chat struct {
	Id     int64 `json:"id"`
	ChatId int64 `json:"chat_id"`
}

type Notification struct {
	Id           int64  `json:"id"`
	Notification string `json:"notification"`
	Datetime     string `json:"datetime"`
	ChatId       int64  `json:"chat_id"`
}

func (db *NotificationsDB) Init(filename string) error {
	db.filename = filename
	var err error
	db.conn, err = sql.Open("sqlite3", db.filename)
	return err
}

func (db *NotificationsDB) Close() {
	err := db.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (db *NotificationsDB) InsertChat(chat Chat) (int64, error) {
	err := db.conn.QueryRow("SELECT id FROM chat WHERE chat_id=?", chat.ChatId).Scan(&chat.Id)
	if err == nil {
		return chat.Id, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return -1, err
	}
	res, err := db.conn.Exec("INSERT INTO chat(chat_id) VALUES(?)", chat.ChatId)
	if err != nil {
		return -1, err
	}
	lastId, err := res.LastInsertId()
	return lastId, nil
}

func (db *NotificationsDB) InsertNotification(notification Notification) (int64, error) {
	stmt, err := db.conn.Prepare("INSERT INTO notification (notification, datetime, chat_id) VALUES (?, ?, ?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(notification.Notification, notification.Datetime, notification.ChatId)
	if err != nil {
		return -1, err
	}
	lastId, err := res.LastInsertId()
	return lastId, nil
}
