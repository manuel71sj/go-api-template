package services

import (
	"gorm.io/gorm"
	"manuel71sj/go-api-template/api/repository"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
	"manuel71sj/go-api-template/models/dto"
	"manuel71sj/go-api-template/pkg/uuid"
)

// MenuService Service layer
type MenuService struct {
	logger                       lib.Logger
	menuRepository               repository.MenuRepository
	menuActionRepository         repository.MenuActionRepository
	menuActionResourceRepository repository.MenuActionResourceRepository
}

// WithTrx delegates transaction to repository database
func (s MenuService) WithTrx(trxHandle *gorm.DB) MenuService {
	s.menuRepository = s.menuRepository.WithTrx(trxHandle)
	s.menuActionRepository = s.menuActionRepository.WithTrx(trxHandle)
	s.menuActionResourceRepository = s.menuActionResourceRepository.WithTrx(trxHandle)

	return s
}

func (s MenuService) Check(item *models.Menu) error {
	result, err := s.menuRepository.Query(&models.MenuQueryParam{
		Name:     item.Name,
		ParentID: item.ParentID,
	})

	if err != nil {
		return err
	} else if len(result.List) > 0 {
		return errors.MenuAlreadyExists
	}

	return nil
}

func (s MenuService) Query(param *models.MenuQueryParam) (*models.MenuQueryResult, error) {
	menuQR, err := s.menuRepository.Query(param)
	if err != nil {
		return nil, err
	}

	if !param.IncludeActions {
		return menuQR, nil
	}

	menuActionQR, err := s.menuActionRepository.Query(&models.MenuActionQueryParam{
		PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})
	if err != nil {
		return nil, err
	}

	menuResourceQR, err := s.menuActionResourceRepository.Query(&models.MenuActionResourceQueryParam{
		MenuIDs: menuQR.List.ToIDs(), PaginationParam: dto.PaginationParam{PageSize: 999, Current: 1},
	})
	if err != nil {
		return nil, err
	}

	menuQR.List.FillMenuAction(menuActionQR.List.ToMenuIDMap(), menuResourceQR.List.ToActionIDMap())

	return menuQR, nil
}

func (s MenuService) GetMenuActions(id string) (models.MenuActions, error) {
	paginationParam := dto.PaginationParam{PageSize: 999, Current: 1}

	menuActionQR, err := s.menuActionRepository.Query(&models.MenuActionQueryParam{
		MenuID: id, PaginationParam: paginationParam,
	})
	if err != nil {
		return nil, err
	} else if len(menuActionQR.List) == 0 {
		return nil, nil
	}

	menuResourceQR, err := s.menuActionResourceRepository.Query(&models.MenuActionResourceQueryParam{
		MenuID: id, PaginationParam: paginationParam,
	})
	if err != nil {
		return nil, err
	}

	menuActionQR.List.FillResources(menuResourceQR.List.ToActionIDMap())
	return menuActionQR.List, nil
}

func (s MenuService) Get(id string) (*models.Menu, error) {
	menu, err := s.menuRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func (s MenuService) Create(menu *models.Menu) (id string, err error) {
	if err = s.Check(menu); err != nil {
		return
	}

	if menu.ParentPath, err = s.GetParentPath(menu.ParentID); err != nil {
		return
	}

	menu.ID = uuid.MustString()
	if err = s.menuRepository.Create(menu); err != nil {
		return
	}

	return menu.ID, nil
}

func (s MenuService) CreateMenus(parentID string, mTrees models.MenuTrees) error {
	for _, mTree := range mTrees {
		menu := &models.Menu{
			Name:      mTree.Name,
			Sequence:  mTree.Sequence,
			Icon:      mTree.Icon,
			Router:    mTree.Router,
			Component: mTree.Component,
			ParentID:  mTree.ParentID,
			Status:    1,
			Hidden:    -1,
		}

		if v := mTree.Hidden; v != 0 {
			menu.Hidden = v
		}

		menuID, err := s.Create(menu)
		if err != nil {
			return err
		}

		if err := s.CreateActions(menu.ID, mTree.Actions); err != nil {
			return err
		}

		if mTree.Children != nil && len(mTree.Children) > 0 {
			err := s.CreateMenus(menuID, mTree.Children)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s MenuService) CreateActions(menuID string, menuActions models.MenuActions) error {
	for _, menuAction := range menuActions {
		menuAction.ID = uuid.MustString()
		menuAction.MenuID = menuID

		if err := s.menuActionRepository.Create(menuAction); err != nil {
			return err
		}

		for _, resource := range menuAction.Resources {
			resource.ID = uuid.MustString()
			resource.ActionID = menuAction.ID

			if err := s.menuActionResourceRepository.Create(resource); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s MenuService) Update(id string, menu *models.Menu) error {
	if id == menu.ParentID {
		return errors.MenuInvalidParent
	}

	// get old menu
	oMenu, err := s.Get(id)
	if err != nil {
		return err
	} else if oMenu.Name != menu.Name {
		if err = s.Check(menu); err != nil {
			return err
		}
	}

	menu.ID = oMenu.ID
	menu.CreatedBy = oMenu.CreatedBy
	menu.CreatedAt = oMenu.CreatedAt

	if menu.ParentID != oMenu.ParentID {
		parentPath, err := s.GetParentPath(menu.ParentID)
		if err != nil {
			return err
		}

		menu.ParentPath = parentPath
	} else {
		menu.ParentPath = oMenu.ParentPath
	}

	if err = s.UpdateChildParentPath(oMenu, menu); err != nil {
		return err
	}

	if err = s.menuRepository.Update(id, menu); err != nil {
		return err
	}

	return nil
}

func (s MenuService) UpdateActions(menuId string, actions models.MenuActions) error {
	oActions, err := s.GetMenuActions(menuId)
	if err != nil {
		return err
	}

	aActions, dActions, uActions := s.CompareActions(oActions, actions)

	err = s.CreateActions(menuId, aActions)
	if err != nil {
		return err
	}

	for _, dAction := range dActions {
		if err = s.menuActionRepository.Delete(dAction.ID); err != nil {
			return err
		}

		if err = s.menuActionResourceRepository.DeleteByActionID(dAction.ID); err != nil {
			return err
		}
	}

	oMap := oActions.ToMap()
	for _, uAction := range uActions {
		// old menu action
		oAction := oMap[uAction.Code]

		// update action name
		if uAction.Name != oAction.Name {
			oAction.Name = uAction.Name
			if err = s.menuActionRepository.Update(uAction.ID, oAction); err != nil {
				return err
			}
		}

		// compare resources to update
		aResources, dResources := s.CompareResources(oAction.Resources, uAction.Resources)
		for _, aResource := range aResources {
			aResource.ID = uuid.MustString()
			aResource.ActionID = oAction.ID

			err := s.menuActionResourceRepository.Create(aResource)
			if err != nil {
				return err
			}
		}

		for _, dResource := range dResources {
			err := s.menuActionResourceRepository.Delete(dResource.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s MenuService) Delete(id string) error {
	_, err := s.menuRepository.Get(id)
	if err != nil {
		return err
	}

	menuQR, err := s.menuRepository.Query(&models.MenuQueryParam{
		ParentID: id,
	})
	if err != nil {
		return err
	} else if menuQR.Pagination.Total > 0 {
		return errors.MenuNotAllowDeleteWithChild
	}

	if err = s.menuActionResourceRepository.DeleteByMenuID(id); err != nil {
		return err
	}

	if err = s.menuActionRepository.DeleteByMenuID(id); err != nil {
		return err
	}

	if err = s.menuRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s MenuService) UpdateStatus(id string, status int) error {
	_, err := s.menuRepository.Get(id)
	if err != nil {
		return err
	}

	return s.menuRepository.UpdateStatus(id, status)

}

func (s MenuService) GetParentPath(parentID string) (string, error) {
	if parentID == "" {
		return "", nil
	}

	parentMenu, err := s.menuRepository.Get(parentID)
	if err != nil {
		return "", err
	}

	return s.JoinParentPath(parentMenu.ParentPath, parentMenu.ID), nil

}

func (s MenuService) JoinParentPath(parent, id string) string {
	if parent != "" {
		return parent + "/" + id
	}

	return id
}

func (s MenuService) CompareActions(oActions, nActions models.MenuActions) (aList, dList, uList models.MenuActions) {
	oMap := oActions.ToMap()
	nMap := nActions.ToMap()

	for k, item := range nMap {
		if _, ok := oMap[k]; ok {
			uList = append(uList, item)
			delete(oMap, k)

			continue
		}

		aList = append(aList, item)
	}

	for _, item := range oMap {
		dList = append(dList, item)
	}

	return
}

func (s MenuService) CompareResources(oResources, nResources models.MenuActionResources) (aList, dList models.MenuActionResources) {
	oMap := oResources.ToMap()
	nMap := nResources.ToMap()

	for k, item := range nMap {
		if _, ok := oMap[k]; ok {
			delete(oMap, k)
			continue
		}

		aList = append(aList, item)
	}

	for _, item := range oMap {
		dList = append(dList, item)
	}

	return
}

func (s MenuService) UpdateChildParentPath(oMenu, nMenu *models.Menu) error {
	if oMenu.ParentID == nMenu.ParentID {
		return nil
	}

	oPath := s.JoinParentPath(oMenu.ParentPath, oMenu.ID)
	menuQR, err := s.menuRepository.Query(&models.MenuQueryParam{
		PrefixParentPath: oPath,
	})
	if err != nil {
		return err
	}

	nPath := s.JoinParentPath(nMenu.ParentPath, nMenu.ID)
	for _, menu := range menuQR.List {
		err = s.menuRepository.UpdateParentPath(menu.ID, nPath+menu.ParentPath[len(oPath):])
		if err != nil {
			return err
		}
	}

	return nil
}

// NewMenuService creates a new menu service
func NewMenuService(
	logger lib.Logger,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
	menuActionResourceRepository repository.MenuActionResourceRepository,
) MenuService {
	return MenuService{
		logger:                       logger,
		menuRepository:               menuRepository,
		menuActionRepository:         menuActionRepository,
		menuActionResourceRepository: menuActionResourceRepository,
	}
}
