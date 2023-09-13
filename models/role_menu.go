package models

import (
	"manuel71sj/go-api-template/models/database"
	"manuel71sj/go-api-template/models/dto"
)

type RoleMenu struct {
	database.Model
	ID       string `gorm:"column:id;size:36;not null;" json:"id"`
	RoleID   string `gorm:"column:role_id;size:36;not null;index;" json:"role_id" validate:"required"`
	MenuID   string `gorm:"column:menu_id;size:36;not null;index;" json:"menu_id" validate:"required"`
	ActionID string `gorm:"column:action_id;size:36;not null;index;" json:"action_id" validate:"required"`
}

type RoleMenus []*RoleMenu

type RoleMenuQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	RoleID  string
	RoleIDs []string
}

type RoleMenuQueryResult struct {
	List       RoleMenus       `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (r RoleMenus) ToMap() map[string]*RoleMenu {
	m := make(map[string]*RoleMenu)
	for _, item := range r {
		m[item.MenuID+"-"+item.ActionID] = item
	}

	return m
}

func (r RoleMenus) ToRoleIDMap() map[string]RoleMenus {
	m := make(map[string]RoleMenus)
	for _, item := range r {
		m[item.RoleID] = append(m[item.RoleID], item)
	}

	return m
}

func (r RoleMenus) ToMenuIDs() []string {
	var idList []string
	m := make(map[string]struct{})

	for _, item := range r {
		if _, ok := m[item.MenuID]; ok {
			continue
		}
		idList = append(idList, item.MenuID)
		m[item.MenuID] = struct{}{}
	}

	return idList
}

func (r RoleMenus) ToActionIDs() []string {
	idList := make([]string, len(r))

	m := make(map[string]struct{})
	for i, item := range r {
		if _, ok := m[item.ActionID]; ok {
			continue
		}
		idList[i] = item.ActionID
		m[item.ActionID] = struct{}{}
	}

	return idList
}
