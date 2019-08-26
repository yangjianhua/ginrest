package model

import (
	"time"
)

type Session struct {
	BaseModel
	UserId         uint
	SessionId      string `gorm:"unique_index;not null"`
	Ip             string
	ExpireTime     *time.Time
	Authentication string
}
