package models

import (
	"manuel71sj/go-api-template/models/database"
	"manuel71sj/go-api-template/models/dto"
)

type UserRole struct {
	database.Model
	ID     string `gorm:"column:id;size:36;not null;" json:"id"`
	UserID string `gorm:"column:user_id;size:36;index;not null;" json:"user_id"`
	RoleID string `gorm:"column:role_id;size:36;index;not null;" json:"role_id"`
}

type UserRoles []UserRole

func (u UserRoles) ToMap() map[string]*UserRole {
	m := make(map[string]*UserRole)
	for _, item := range u {
		m[item.RoleID] = &item
	}

	return m
}

func (u UserRoles) ToRoleIDs() []string {
	list := make([]string, len(u))
	for i, item := range u {
		list[i] = item.RoleID
	}

	return list
}

func (u UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, item := range u {
		m[item.UserID] = append(m[item.UserID], item)
	}

	return m
}

type UserRoleQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	UserID  string
	UserIDs []string
}

type UserRoleQueryResult struct {
	List       UserRoles       `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
