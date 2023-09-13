package models

import (
	"manuel71sj/go-api-template/models/database"
	"manuel71sj/go-api-template/models/dto"
)

type Role struct {
	database.Model
	ID        string    `gorm:"column:id;size:36;index;not null;" json:"id"`
	Name      string    `gorm:"column:name;not null;" json:"name" validate:"required"`
	Remark    string    `gorm:"column:remark;default:'';" json:"remark" validate:"required"`
	Sequence  int       `gorm:"column:sequence;not null;index;" json:"sequence" validate:"required"`
	Status    int       `gorm:"column:status;not null;default:0;" json:"status" validate:"required,max=1,min=-1"`
	CreatedBy string    `gorm:"column:created_by;not null;" json:"created_by"`
	RoleMenus RoleMenus `gorm:"-" json:"role_menus"`
}

type Roles []*Role

func (r Roles) ToNames() []string {
	names := make([]string, len(r))
	for i, item := range r {
		names[i] = item.Name
	}

	return names
}

func (r Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range r {
		m[item.ID] = item
	}

	return m
}

type RoleQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	Name       string   `query:"name"`
	QueryValue string   `query:"query_value"`
	UserID     string   `query:"user_id"`
	Status     int      `query:"status" validate:"max=1,min=-1"`
}

type RoleQueryResult struct {
	List       Roles           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
