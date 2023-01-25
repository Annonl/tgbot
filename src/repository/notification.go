package repository

import (
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	Id               uuid.UUID
	ChatId           int64
	Message          string
	NotificationDate time.Time
}
