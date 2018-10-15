package model

import "time"

type BaseModel struct {
	ID        uint       `gorm:"primay_key"`
	CreatedAt time.Time  `form:"-"`
	UpdatedAt time.Time  `form:"-"`
	DeletedAt *time.Time `sql:"index" form:"-"`
}
