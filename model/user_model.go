package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	BaseModel
	Name       string     `gorm:"type:varchar(100);not null;unique_index" json:"name" form:"name" validate:"min:4,max:20"`
	Mobile     string     `json:"mobile" form:"mobile"`
	Tel        string     `json:"tel" form:"tel"`
	Birthday   *time.Time `json:"birthday" form:"birthday"`
	Email      string     `gorm:"type:varchar(100);unique_index" json:"email" form:"email"`
	Password   string     `gorm:"type:varchar(200)" json:"-" form:"password"`
	Avator     string     `gorm:"type:varchar(400)" json:"avator" form:"avator"`
	IsAdmin    bool       `json:"is_admin" form:"is_admin"`
	Access     string     `gorm:"-" json:"access"`
	CreatedUid uint       `json:"created_uid" form:"created_uid"`
	UpdatedUid uint       `json:"updated_uid" form:"updated_uid"`
}

func (this *User) BeforeCreate() (err error) {
	this.CreatedUid = Uid
	return nil
}

func (this *User) BeforeSave() (err error) {
	this.UpdatedUid = Uid
	return nil
}

func (this *User) BeforeDelete(scope *gorm.Scope) (err error) {
	if this.IsAdmin {
		err = errors.New("Cannot delete admin user.")
	}
	return
}

// In this event to get the "access" of the user's group
func (this *User) AfterFind(scope *gorm.Scope) (err error) {
	this.Access = "admin,super_admin"
	return
}
