package model

import (
	"time"

	_ "github.com/jinzhu/gorm"
)

// This is a test struct for gin.ShouldBind
// Not work on created_at, updated_at and deleted_at fields.

type Test struct {
	ID        uint       `gorm:"primay_key"`
	Name      string     `form:"name"`
	CreatedAt time.Time  `form:"-"`
	UpdatedAt time.Time  `form:"-"`
	DeletedAt *time.Time `sql:"index" form:"-"`
}
