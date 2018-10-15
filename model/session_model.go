package model

import (
	"time"

	_ "github.com/jinzhu/gorm"
)

type Session struct {
	// gorm.Model
	BaseModel
	UserId         uint
	SessionId      string `gorm:"unique_index;not null"`
	Ip             string
	ExpireTime     *time.Time
	Authentication string
}
