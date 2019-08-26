package model

type GroupUser struct {
	BaseModel
	UserId      uint   `grom:"json:user_id"`
	GroupId     uint   `gorm:"json:group_id"`
	Group       Group  `gorm:"foreignkey:group_id"`
	User        User   `gorm:"foreignkey:user_id"`
	IsAdmin     bool   `json:"is_admin" form:"is_admin"`
	IsLeader    bool   `json:"is_leader" form:"is_leader"`
	Description string `gorm:"type:varchar(400)" json:"description"`
	CreatedUid  uint   `json:"created_uid" form:"created_uid"`
	UpdatedUid  uint   `json:"updated_uid" form:"updated_uid"`
}

func (this *GroupUser) BeforeCreate() (err error) {
	this.CreatedUid = Uid
	return nil
}

func (this *GroupUser) BeforeSave() (err error) {
	this.UpdatedUid = Uid
	return nil
}
