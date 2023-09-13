package services

import (
	"gorm.io/gorm"
	"manuel71sj/go-api-template/api/repository"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
	"manuel71sj/go-api-template/pkg/uuid"
)

// RoleService service layer
type RoleService struct {
	logger               lib.Logger
	casbinService        CasbinService
	userRepository       repository.UserRepository
	roleRepository       repository.RoleRepository
	roleMenuRepository   repository.RoleMenuRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
}

// WithTrx delegates transaction to repository database
func (s RoleService) WithTrx(trxHandle *gorm.DB) RoleService {
	s.roleRepository = s.roleRepository.WithTrx(trxHandle)
	s.userRepository = s.userRepository.WithTrx(trxHandle)
	s.roleMenuRepository = s.roleMenuRepository.WithTrx(trxHandle)

	return s
}

func (s RoleService) Query(param *models.RoleQueryParam) (roleQR *models.RoleQueryResult, err error) {
	return s.roleRepository.Query(param)
}

func (s RoleService) QueryRoleMenus(roleID string) (models.RoleMenus, error) {
	roleMenuQR, err := s.roleMenuRepository.Query(&models.RoleMenuQueryParam{
		RoleID: roleID,
	})
	if err != nil {
		return nil, err
	}

	return roleMenuQR.List, nil
}

func (s RoleService) Get(id string) (*models.Role, error) {
	role, err := s.roleRepository.Get(id)
	if err != nil {
		return nil, err
	}

	roleMenus, err := s.QueryRoleMenus(id)
	if err != nil {
		return nil, err
	}

	role.RoleMenus = roleMenus
	return role, nil
}

func (s RoleService) Check(item *models.Role) error {
	qr, err := s.roleRepository.Query(&models.RoleQueryParam{Name: item.Name})
	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.RoleAlreadyExists
	}

	return nil
}

func (s RoleService) CheckRoleMenu(rMenu *models.RoleMenu) error {
	if _, err := s.menuRepository.Get(rMenu.MenuID); err != nil {
		return errors.Wrap(err, "menu id")
	}

	if _, err := s.menuActionRepository.Get(rMenu.ActionID); err != nil {
		return errors.Wrap(err, "menu action id")
	}

	return nil
}

func (s RoleService) CompareRoleMenus(oRoleMenus, nRoleMenus models.RoleMenus) (aList, dList models.RoleMenus) {
	oMap := oRoleMenus.ToMap()
	nMap := nRoleMenus.ToMap()

	for k, nRoleMenu := range nMap {
		if _, ok := oMap[k]; ok {
			delete(oMap, k)
			continue
		}

		aList = append(aList, nRoleMenu)
	}

	for _, oRoleMenu := range oMap {
		dList = append(dList, oRoleMenu)
	}

	return
}

func (s RoleService) Create(role *models.Role) (id string, err error) {
	if err = s.Check(role); err != nil {
		return
	}

	role.ID = uuid.MustString()
	for _, roleMenu := range role.RoleMenus {
		roleMenu.ID = uuid.MustString()
		roleMenu.RoleID = role.ID

		if err = s.CheckRoleMenu(roleMenu); err != nil {
			return
		}

		if err = s.roleMenuRepository.Create(roleMenu); err != nil {
			return
		}
	}

	if err = s.roleRepository.Create(role); err != nil {
		return
	}

	_ = s.casbinService.Enforcer.LoadPolicy()
	return role.ID, nil
}

func (s RoleService) Update(id string, role *models.Role) error {
	oRole, err := s.Get(id)
	if err != nil {
		return err
	} else if role.Name != oRole.Name {
		if err = s.Check(role); err != nil {
			return err
		}
	}

	role.ID = oRole.ID
	role.CreatedBy = oRole.CreatedBy
	role.CreatedAt = oRole.CreatedAt

	aRoleMenus, dRoleMenus := s.CompareRoleMenus(oRole.RoleMenus, role.RoleMenus)
	for _, aRoleMenu := range aRoleMenus {
		aRoleMenu.ID = uuid.MustString()
		aRoleMenu.RoleID = id

		if err := s.CheckRoleMenu(aRoleMenu); err != nil {
			return err
		}

		if err := s.roleMenuRepository.Create(aRoleMenu); err != nil {
			return err
		}
	}

	for _, dRoleMenu := range dRoleMenus {
		if err := s.roleMenuRepository.Delete(dRoleMenu.ID); err != nil {
			return err
		}
	}

	if err := s.roleRepository.Update(id, role); err != nil {
		return err
	}

	_ = s.casbinService.Enforcer.LoadPolicy()
	return nil
}

func (s RoleService) Delete(id string) error {
	_, err := s.roleRepository.Get(id)
	if err != nil {
		return err
	}

	userQR, err := s.userRepository.Query(&models.UserQueryParam{
		RoleIDs: []string{id},
	})
	if err != nil {
		return err
	} else if userQR.Pagination.Total > 0 {
		return errors.RoleNotAllowDeleteWithUser
	}

	if err := s.roleMenuRepository.DeleteByRoleID(id); err != nil {
		return err
	}

	if err := s.roleRepository.Delete(id); err != nil {
		return err
	}

	_ = s.casbinService.Enforcer.LoadPolicy()
	return nil
}

func (s RoleService) UpdateStatus(id string, status int) error {
	_, err := s.roleRepository.Get(id)
	if err != nil {
		return err
	}

	if err := s.roleRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	_ = s.casbinService.Enforcer.LoadPolicy()
	return nil
}

// NewRoleService creates a new role service
func NewRoleService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	roleRepository repository.RoleRepository,
	roleMenuRepository repository.RoleMenuRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) RoleService {
	return RoleService{
		logger:               logger,
		casbinService:        casbinService,
		userRepository:       userRepository,
		roleRepository:       roleRepository,
		roleMenuRepository:   roleMenuRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
	}
}
