package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	// gorm.Model
	BaseModel
	Name     string     `gorm:"type:varchar(100);not null;unique_index" json:"name" form:"name" validate:"min:4,max:20"`
	Mobile   string     `json:"mobile" form:"mobile"`
	Tel      string     `json:"tel" form:"tel"`
	Birthday *time.Time `json:"birthday" form:"birthday"`
	Email    string     `gorm:"type:varchar(100);unique_index" json:"email" form:"email"`
	Password string     `gorm:"type:varchar(200)" json:"-" form:"password"`
	IsAdmin  bool       `json:"is_admin" form:"is_admin"`
}

func (this *User) BeforeDelete(scope *gorm.Scope) (err error) {
	if this.IsAdmin {
		err = errors.New("Cannot delete admin user.")
	}
	return
}
