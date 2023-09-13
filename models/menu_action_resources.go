package models

import (
	"manuel71sj/go-api-template/models/database"
	"manuel71sj/go-api-template/models/dto"
)

type MenuActionResource struct {
	database.Model
	ID       string `gorm:"column:id;size:36;index;not null;" json:"-" yaml:"-"`
	ActionID string `gorm:"column:action_id;size:36;index;not null;" json:"-" yaml:"-"`
	Method   string `gorm:"column:method;not null;" json:"method" validate:"required" yaml:"method"`
	Path     string `gorm:"column:path;not null;" json:"path" validate:"required" yaml:"path"`
}

type MenuActionResources []*MenuActionResource

func (mars MenuActionResources) ToMap() map[string]*MenuActionResource {
	m := make(map[string]*MenuActionResource)
	for _, item := range mars {
		m[item.Method+item.Path] = item
	}

	return m
}

func (mars MenuActionResources) ToActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)
	for _, item := range mars {
		m[item.ActionID] = append(m[item.ActionID], item)
	}

	return m
}

type MenuActionResourceQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	MenuID  string
	MenuIDs []string
}

type MenuActionResourceQueryResult struct {
	List       MenuActionResources `json:"list"`
	Pagination *dto.Pagination     `json:"pagination"`
}
