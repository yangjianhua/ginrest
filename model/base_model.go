package model

import "time"

var Uid uint = 0 // The Uid to created_uid and updated_uid

// For gin.ShouldBind Event could cause error, so all models cannot just from gorm.Model
// Should inherit the BaseModel, not gorm.Model
// Because the BaseModel have `form:"-"`
// The Model Must import "gorm", or the Hooks will not work (if the hooks does not have the param like scope *gorm.Scope)

type BaseModel struct {
	ID        uint       `gorm:"primay_key" json:"id"`
	CreatedAt time.Time  `form:"-" json:"created_at"`
	UpdatedAt time.Time  `form:"-" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" form:"-" json:"deleted_at"`
}
