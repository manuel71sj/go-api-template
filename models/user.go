package models

import (
	"manuel71sj/go-api-template/models/database"
	"manuel71sj/go-api-template/models/dto"
)

type User struct {
	database.Model
	ID        string    `gorm:"column:id;size:36;index;not null;" json:"id"`
	Username  string    `gorm:"column:username;size:64;not null;index;" json:"username" validate:"required"`
	Realname  string    `gorm:"column:realname;size:64;not null;" json:"realname" validate:"required"`
	Password  string    `gorm:"column:password;not null;" json:"password" json:"phone"`
	Email     string    `gorm:"column:email;default:'';" json:"email"`
	Phone     string    `gorm:"column:phone;default:'';" json:"phone"`
	Status    int       `gorm:"column:status;not null;default:0;" json:"status" validate:"required,max=1,min=-1"`
	CreatedBy string    `gorm:"column:created_by;not null;" json:"created_by"`
	UserRoles UserRoles `gorm:"-" json:"user_roles"`
}

func (u *User) CleanSecure() *User {
	u.Password = ""
	return u
}

type Users []*User

func (u Users) ToIDs() []string {
	ids := make([]string, len(u))
	for i, item := range u {
		ids[i] = item.ID
	}

	return ids
}

type UserInfo struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	Realname string `json:"realname"`
	Roles    Roles  `json:"roles"`
}

type UserQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	QueryPassword bool
	Username      string   `query:"username"`
	Realname      string   `query:"realname"`
	QueryValue    string   `query:"query_value"`
	Status        int      `query:"status" validate:"max=1,min=-1"`
	RoleIDs       []string `query:"-"`
}

type UserQueryResult struct {
	List       Users           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
