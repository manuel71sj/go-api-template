package models

import (
	"manuel71sj/go-api-template/models/database"
	"manuel71sj/go-api-template/models/dto"
)

type MenuAction struct {
	database.Model
	ID        string              `gorm:"column:id;size:36;not null;index;" json:"id" yaml:"-"`
	MenuID    string              `gorm:"column:menu_id;size:36;not null;index;" json:"menu_id" yaml:"-"`
	Code      string              `gorm:"column:code;not null;" json:"code" validate:"required" yaml:"code"`
	Name      string              `gorm:"column:name;not null;" json:"name" validate:"required" yaml:"name"`
	Resources MenuActionResources `gorm:"-" json:"resources" yaml:"resources"`
}

type MenuActions []*MenuAction

func (ma MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, v := range ma {
		m[v.Code] = v
	}
	return m
}

func (ma MenuActions) FillResources(maResources map[string]MenuActionResources) {
	for i, v := range ma {
		ma[i].Resources = maResources[v.ID]
	}
}

func (ma MenuActions) ToMenuIDMap() map[string]MenuActions {
	m := make(map[string]MenuActions)
	for _, v := range ma {
		m[v.MenuID] = append(m[v.MenuID], v)
	}
	return m
}

type MenuActionQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	MenuID string
	IDs    []string
}

type MenuActionQueryResult struct {
	List       MenuActions     `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
