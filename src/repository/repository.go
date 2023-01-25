package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

var (
	host            = os.Getenv("host")
	port            = os.Getenv("port")
	user            = os.Getenv("user")
	password        = os.Getenv("password")
	dbname          = os.Getenv("dbname")
	dbConnectString = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)
)

type NotificationRepository struct {
	db *sql.DB
}

func CreateRepository() (*NotificationRepository, error) {
	db, err := sql.Open("postgres", dbConnectString)

	return &NotificationRepository{db: db}, err
}

func (r *NotificationRepository) SetNotification(notification Notification) error {
	query := `INSERT INTO "notifications"("chatid", "message", "notificationdate") values($1, $2, $3)`

	_, e := r.db.Exec(query, notification.ChatId, notification.Message, notification.NotificationDate)

	return e
}

func (r *NotificationRepository) GetTodayChatId() ([]int64, error) {
	rows, err := r.db.Query(`SELECT DISTINCT chatid FROM notifications WHERE notificationdate = CURRENT_DATE`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var result []int64

	for rows.Next() {
		var chatId int64

		if errRead := rows.Scan(&chatId); errRead != nil {
			return nil, errRead
		}
		result = append(result, chatId)
	}

	return result, nil
}

func (r *NotificationRepository) GetTodayNotifications(chatId int64) ([]Notification, error) {
	rows, err := r.db.Query(`SELECT * FROM notifications WHERE notificationdate = CURRENT_DATE AND chatid = $1`, chatId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var notifications []Notification

	for rows.Next() {
		var notification Notification

		if errRead := rows.Scan(&notification.Id, &notification.ChatId, &notification.Message, &notification.NotificationDate); errRead != nil {
			return nil, errRead
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
