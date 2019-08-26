package model

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Group struct {
	BaseModel
	Name        string `gorm:"type:varchar(200);not null;unique_index" json:"name" form:"name" validate:"min:4;max:200"`
	ParentId    uint   `gorm:"json:parent_id" form:"parent_id"`
	Description string `gorm:"type:varchar(400);" json:"description" form:"description"`
	CreatedUid  uint   `json:"created_uid"`
	UpdatedUid  uint   `json:"updated_uid"`
}

// Hooks Here.

func (this *Group) BeforeCreate(scope *gorm.Scope) (err error) {
	this.CreatedUid = Uid
	return nil
}

func (this *Group) BeforeSave(scope *gorm.Scope) (err error) {
	this.UpdatedUid = Uid
	return nil
}

func (this *Group) BeforeDelete(scope *gorm.Scope) (err error) {
	var group Group
	scope.DB().Where("parent_id = ?", this.ParentId).First(&group)
	if group.ID > 0 {
		err = errors.New("The group have sub group, delete them first.")
	}
	return nil
}
