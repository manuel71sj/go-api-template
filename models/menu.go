package models

import (
	"manuel71sj/go-api-template/models/database"
	"manuel71sj/go-api-template/models/dto"
	"strings"
)

type Menu struct {
	database.Model
	ID         string      `gorm:"column:id;size:36;not null;index;" json:"id"`
	Name       string      `gorm:"column:name;not null;index;" json:"name" validate:"required"`
	Sequence   int         `gorm:"column:sequence;not null;index;" json:"sequence" validate:"required"`
	Icon       string      `gorm:"column:icon;" json:"icon" validate:"required"`
	Router     string      `gorm:"column:router;" json:"router"`
	Component  string      `gorm:"column:component;" json:"component"`
	ParentID   string      `gorm:"column:parent_id;size:36;index;" json:"parent_id"`
	ParentPath string      `gorm:"column:parent_path;" json:"parent_path"`
	Hidden     int         `gorm:"column:hidden;not null;" json:"hidden" validate:"required,max=1,min=-1"`
	Status     int         `gorm:"column:status;not null;" json:"status" validate:"required,max=1,min=-1"`
	Remark     string      `gorm:"column:remark;" json:"remark" validate:"required"`
	CreatedBy  string      `gorm:"column:created_by;not null;" json:"created_by"`
	Actions    MenuActions `gorm:"-" json:"actions,omitempty"`
}

type MenuTree struct {
	ID         string      `yaml:"-" json:"id"`
	Name       string      `yaml:"name" json:"name"`
	Icon       string      `yaml:"icon" json:"icon"`
	Router     string      `yaml:"router,omitempty" json:"router"`
	Component  string      `yaml:"component,omitempty" json:"component"`
	ParentID   string      `yaml:"-" json:"parent_id"`
	ParentPath string      `yaml:"-" json:"parent_path"`
	Sequence   int         `yaml:"sequence" json:"sequence"`
	Hidden     int         `yaml:"-" json:"hidden"`
	Status     int         `yaml:"-" json:"status"`
	Actions    MenuActions `yaml:"actions,omitempty" json:"actions"`
	Children   MenuTrees   `yaml:"children,omitempty" json:"children,omitempty"`
}

type Menus []*Menu

func (ms Menus) Len() int {
	return len(ms)
}

func (ms Menus) Less(i, j int) bool {
	return ms[i].Sequence < ms[j].Sequence
}

func (ms Menus) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms Menus) ToMap() map[string]*Menu {
	m := make(map[string]*Menu)
	for _, menu := range ms {
		m[menu.ID] = menu
	}

	return m
}

func (ms Menus) SplitParentIDs() []string {
	idList := make([]string, len(ms))
	mIDList := make(map[string]struct{})

	for _, v := range ms {
		if _, ok := mIDList[v.ID]; ok || v.ParentPath == "" {
			continue
		}

		for _, pp := range strings.Split(v.ParentPath, "/") {
			if _, ok := mIDList[pp]; ok {
				continue
			}

			idList = append(idList, pp)
			mIDList[pp] = struct{}{}
		}
	}

	return idList
}

func (ms Menus) ToIDs() []string {
	ids := make([]string, len(ms))
	for i, v := range ms {
		ids[i] = v.ID
	}

	return ids
}

func (ms Menus) ToMenuTrees() MenuTrees {
	menuTrees := make(MenuTrees, len(ms))
	for i, v := range ms {
		menuTrees[i] = &MenuTree{
			ID:         v.ID,
			Name:       v.Name,
			Icon:       v.Icon,
			Router:     v.Router,
			Component:  v.Component,
			ParentID:   v.ParentID,
			ParentPath: v.ParentPath,
			Sequence:   v.Sequence,
			Hidden:     v.Hidden,
			Status:     v.Status,
			Actions:    v.Actions,
		}
	}

	return menuTrees.ToTree()
}

func (ms Menus) FillMenuAction(mActions map[string]MenuActions, mResources map[string]MenuActionResources) Menus {
	for _, item := range ms {
		if v, ok := mActions[item.ID]; ok {
			item.Actions = v
			item.Actions.FillResources(mResources)
		}
	}

	return ms
}

type MenuTrees []*MenuTree

func (ms MenuTrees) ToTree() MenuTrees {
	menuTreeMap := make(map[string]*MenuTree)
	for _, v := range ms {
		menuTreeMap[v.ID] = v
	}

	menuTrees := make(MenuTrees, 0)
	for _, v := range ms {
		if v.ParentID == "" {
			menuTrees = append(menuTrees, v)
			continue
		}

		if parentMenuTree, ok := menuTreeMap[v.ParentID]; ok {
			if parentMenuTree.Children == nil {
				parentMenuTree.Children = MenuTrees{v}
				continue
			}

			parentMenuTree.Children = append(parentMenuTree.Children, v)
		}
	}

	return menuTrees
}

type MenuQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs              []string `query:"ids"`
	Name             string   `query:"name"`
	PrefixParentPath string   `query:"prefix_parent_path"`
	QueryValue       string   `query:"query_value"`
	ParentID         string   `query:"parent_id"`
	Hidden           int      `query:"hidden" validate:"max=1,min=-1"`
	Status           int      `query:"status" validate:"max=1,min=-1"`
	Tree             bool     `query:"tree"`
	IncludeActions   bool     `query:"include_actions"`
}

type MenuQueryResult struct {
	List       Menus           `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
