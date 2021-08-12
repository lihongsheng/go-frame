package model

import (
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	table string
}
type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(20);not null " json:"user_name" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500);not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Status   int    `gorm:"type:int;DEFAULT:2" json:"status" validate:"required,gte=2" label:"角色码"`
	//CreateT  int    `gorm:"type:int;"`
	//UpdateT  int    `gorm:"type:int;"`
}

func NewUserModel() (*UserModel) {
	return &UserModel{
		table: "`user`",
	}
}

func (u *UserModel)GetUser(id int) (User, error) {
	var user User
	err := DB.Model(&user).Limit(1).Where("id = ?", id).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, err
}

func (u *UserModel)FindMany(name string,page int,pageSize int) ([]User, error) {
	var user []User
	err := DB.Limit(pageSize).Offset((page - 1) * pageSize).Where("user_name LIKE ?", "%"+name+"%").Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, err
}
