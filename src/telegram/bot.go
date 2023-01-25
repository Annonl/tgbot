package telegram

import (
	"awesomeProject/repository"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
	"log"
	"strings"
	"time"
)

type TelegramBot struct {
	bot  *tgbotapi.BotAPI
	repo *repository.NotificationRepository
}

func CreateTelegramBot(token string) (*TelegramBot, error) {
	if token == "" {
		return nil, errors.New("передан пустой токен")
	}

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	log.Printf("Авторизация в аккаунт %s", bot.Self.UserName)
	repo, errRepo := repository.CreateRepository()

	if errRepo != nil {
		return nil, errRepo
	}
	return &TelegramBot{bot: bot, repo: repo}, err
}

func (b *TelegramBot) Start() {
	c := cron.New()
	c.AddFunc("@daily", b.checkDates)
	c.Start()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 600

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			b.bot.Send(msg)

			params := strings.Split(msg.Text, " ")

			notificationTime, errTime := time.Parse("02.01.2006", params[1])
			log.Println(notificationTime)
			if errTime != nil {
				log.Panic(errTime)
			}

			errResult := b.repo.SetNotification(repository.Notification{
				ChatId:           update.Message.Chat.ID,
				Message:          params[0],
				NotificationDate: notificationTime})

			if errResult != nil {
				log.Panic(errResult)
			}
		}
	}
}

func (b *TelegramBot) checkDates() {
	chatId, err := b.repo.GetTodayChatId()

	if err != nil {
		log.Panic(err)
	}

	for _, id := range chatId {
		notifications, err := b.repo.GetTodayNotifications(id)
		if err != nil {
			log.Panic(err)
		}

		b.sendNotifications(notifications, id)
	}

}

func (b *TelegramBot) sendNotifications(notifications []repository.Notification, chatId int64) {
	var result string
	result += "Сегодня:\n"
	for _, notification := range notifications {
		result += notification.Message + "\n"
	}
	msg := tgbotapi.NewMessage(chatId, result)

	b.bot.Send(msg)
	log.Printf("[%s] %s", chatId, result)
}
