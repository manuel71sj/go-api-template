package services

import (
	"gorm.io/gorm"
	"manuel71sj/go-api-template/api/repository"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
	"manuel71sj/go-api-template/models/dto"
	"manuel71sj/go-api-template/pkg/hash"
	"manuel71sj/go-api-template/pkg/uuid"
	"sort"
)

// UserService service layer
type UserService struct {
	logger               lib.Logger
	config               lib.Config
	casbinService        CasbinService
	userRepository       repository.UserRepository
	userRoleRepository   repository.UserRoleRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
	roleRepository       repository.RoleRepository
	roleMenuRepository   repository.RoleMenuRepository
}

func (s UserService) GetSuperAdmin() *models.User {
	admin := s.config.SuperAdmin
	return &models.User{
		ID:       admin.Username,
		Username: admin.Username,
		Realname: admin.RealName,
		Password: admin.Password,
	}
}

// WithTrx delegates transaction to repository database
func (s UserService) WithTrx(trxHandle *gorm.DB) UserService {
	s.userRepository = s.userRepository.WithTrx(trxHandle)
	s.userRoleRepository = s.userRoleRepository.WithTrx(trxHandle)

	return s
}

func (s UserService) Query(param *models.UserQueryParam) (userQR *models.UserQueryResult, err error) {
	if userQR, err = s.userRepository.Query(param); err != nil {
		return
	}

	uRoleQR, err := s.userRoleRepository.Query(
		&models.UserRoleQueryParam{UserIDs: userQR.List.ToIDs()},
	)
	if err != nil {
		return
	}

	m := uRoleQR.List.ToUserIDMap()
	for _, user := range userQR.List {
		if uRoles, ok := m[user.ID]; ok {
			user.UserRoles = uRoles
		}
	}

	return
}

func (s UserService) Verify(username, password string) (*models.User, error) {
	// super admin user
	admin := s.GetSuperAdmin()
	if admin.Username == username && admin.Password == password {
		return admin, nil
	}

	user, err := s.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if user.Password != hash.SHA256(password) {
		return nil, errors.UserInvalidPassword
	} else if user.Status != 1 {
		return nil, errors.UserIsDisable
	}

	return user, nil
}

func (s UserService) Check(user *models.User) error {
	if user.Username == s.GetSuperAdmin().Username {
		return errors.UserInvalidUsername
	}

	if qr, err := s.Query(&models.UserQueryParam{
		Username: user.Username,
	}); err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.UserAlreadyExists
	}

	return nil
}

func (s UserService) GetUserInfo(ID string) (*models.UserInfo, error) {
	if s.GetSuperAdmin().ID == ID {
		user := s.GetSuperAdmin()
		return &models.UserInfo{
			ID:       user.Username,
			Username: user.Username,
			Realname: user.Realname,
		}, nil
	}

	user, err := s.Get(ID)
	if err != nil {
		return nil, err
	}

	userinfo := &models.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Realname: user.Realname,
	}

	userRoleQR, err := s.userRoleRepository.Query(&models.UserRoleQueryParam{
		UserID: ID,
	})
	if err != nil {
		return nil, err
	}

	if roleIDs := userRoleQR.List.ToRoleIDs(); len(roleIDs) > 0 {
		roleQR, err := s.roleRepository.Query(&models.RoleQueryParam{
			IDs:    roleIDs,
			Status: 1,
		})
		if err != nil {
			return nil, err
		}

		userinfo.Roles = roleQR.List
	}

	return userinfo, nil
}

func (s UserService) GetUserMenuTrees(ID string) (models.MenuTrees, error) {
	if s.GetSuperAdmin().ID == ID {
		menuQR, err := s.menuRepository.Query(&models.MenuQueryParam{
			Status:     1,
			OrderParam: dto.OrderParam{Key: "sequence", Direction: dto.OrderByASC},
		})
		if err != nil {
			return nil, err
		}

		return menuQR.List.ToMenuTrees(), nil
	}

	var (
		userRoleQR *models.UserRoleQueryResult
		roleMenuQR *models.RoleMenuQueryResult
		menuQR     *models.MenuQueryResult
		err        error
	)

	if userRoleQR, err = s.userRoleRepository.Query(&models.UserRoleQueryParam{
		UserID: ID,
	}); err != nil {
		return nil, err
	} else if len(userRoleQR.List) == 0 {
		return nil, errors.UserNoPermission
	}

	if roleMenuQR, err = s.roleMenuRepository.Query(&models.RoleMenuQueryParam{
		RoleIDs: userRoleQR.List.ToRoleIDs(),
	}); err != nil {
		return nil, err
	} else if len(roleMenuQR.List) == 0 {
		return nil, errors.UserNoPermission
	}

	if menuQR, err = s.menuRepository.Query(&models.MenuQueryParam{
		IDs:        roleMenuQR.List.ToMenuIDs(),
		Status:     1,
		OrderParam: dto.OrderParam{Key: "sequence", Direction: dto.OrderByASC},
	}); err != nil {
		return nil, err
	} else if len(menuQR.List) == 0 {
		return nil, errors.UserNoPermission
	}

	menuMap := menuQR.List.ToMap()
	var parentIDs []string
	for _, parentID := range menuQR.List.SplitParentIDs() {
		if _, ok := menuMap[parentID]; !ok {
			parentIDs = append(parentIDs, parentID)
		}
	}

	if len(parentIDs) > 0 {
		parentMenuQR, err := s.menuRepository.Query(&models.MenuQueryParam{
			IDs: parentIDs,
		})
		if err != nil {
			return nil, err
		}

		menuQR.List = append(menuQR.List, parentMenuQR.List...)
	}

	sort.Sort(menuQR.List)
	return menuQR.List.ToMenuTrees(), nil
}

func (s UserService) GetByUsername(username string) (*models.User, error) {
	userQR, err := s.Query(
		&models.UserQueryParam{Username: username, QueryPassword: true},
	)
	if err != nil {
		return nil, err
	} else if len(userQR.List) == 0 {
		return nil, errors.UserRecordNotFound
	}

	// set schema
	user := userQR.List[0]

	// get user roles
	userRoleQR, err := s.userRoleRepository.Query(
		&models.UserRoleQueryParam{UserID: user.ID},
	)
	if err != nil {
		return nil, err
	}

	user.UserRoles = userRoleQR.List
	return user, nil
}

func (s UserService) Get(id string) (*models.User, error) {
	user, err := s.userRepository.Get(id)
	if err != nil {
		return nil, err
	}

	userRoleQR, err := s.userRoleRepository.Query(
		&models.UserRoleQueryParam{UserID: id},
	)
	if err != nil {
		return nil, err
	}

	user.UserRoles = userRoleQR.List
	return user, nil
}

func (s UserService) Create(user *models.User) (id string, err error) {
	if err = s.Check(user); err != nil {
		return
	}

	user.Password = hash.SHA256(user.Password)
	user.ID = uuid.MustString()

	for _, userRole := range user.UserRoles {
		userRole.ID = uuid.MustString()
		userRole.UserID = user.ID
		if err = s.userRoleRepository.Create(&userRole); err != nil {
			return
		}
	}

	if err = s.userRepository.Create(user); err != nil {
		return
	}

	_ = s.casbinService.Enforcer.LoadPolicy()
	return user.ID, nil
}

func (s UserService) Update(id string, user *models.User) error {
	oUser, err := s.Get(id)
	if err != nil {
		return err
	} else if user.Username != oUser.Username {
		if err := s.Check(user); err != nil {
			return err
		}
	}

	if user.Password != "" {
		user.Password = hash.SHA256(user.Password)
	} else {
		user.Password = oUser.Password
	}

	user.ID = oUser.ID
	user.CreatedAt = oUser.CreatedAt
	user.CreatedBy = oUser.CreatedBy

	aUserRoles, dUserRoles := s.CompareUserRoles(oUser.UserRoles, user.UserRoles)
	for _, aUserRole := range aUserRoles {
		aUserRole.ID = uuid.MustString()
		aUserRole.UserID = id
		if err := s.userRoleRepository.Create(&aUserRole); err != nil {
			return err
		}
	}

	for _, dUserRole := range dUserRoles {
		if err := s.userRoleRepository.Delete(dUserRole.ID); err != nil {
			return err
		}
	}

	if err := s.userRepository.Update(id, user); err != nil {
		return err
	}

	_ = s.casbinService.Enforcer.LoadPolicy()
	return nil
}

func (s UserService) CompareUserRoles(oUserRoles, nUserRoles models.UserRoles) (aList, dList models.UserRoles) {
	oMap := oUserRoles.ToMap()
	nMap := nUserRoles.ToMap()

	for k, nUserRole := range nMap {
		if _, ok := oMap[k]; ok {
			delete(oMap, k)
			continue
		}

		aList = append(aList, *nUserRole)
	}

	for _, oUserRole := range oMap {
		dList = append(dList, *oUserRole)
	}

	return
}

func (s UserService) Delete(id string) error {
	_, err := s.userRepository.Get(id)
	if err != nil {
		return err
	}

	if err := s.userRoleRepository.DeleteByUserID(id); err != nil {
		return err
	}

	_ = s.casbinService.Enforcer.LoadPolicy()
	return s.userRepository.Delete(id)
}

func (s UserService) UpdateStatus(id string, status int) error {
	_, err := s.userRepository.Get(id)
	if err != nil {
		return err
	}

	if err = s.userRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	_ = s.casbinService.Enforcer.LoadPolicy()
	return nil
}

// NewUserService creates a new user service
func NewUserService(
	logger lib.Logger,
	userRepository repository.UserRepository,
	userRoleRepository repository.UserRoleRepository,
	roleRepository repository.RoleRepository,
	roleMenuRepository repository.RoleMenuRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
	casbinService CasbinService,
	config lib.Config,
) UserService {
	return UserService{
		logger:               logger,
		config:               config,
		userRepository:       userRepository,
		userRoleRepository:   userRoleRepository,
		roleRepository:       roleRepository,
		roleMenuRepository:   roleMenuRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
		casbinService:        casbinService,
	}
}
