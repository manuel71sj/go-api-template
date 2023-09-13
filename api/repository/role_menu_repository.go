package repository

import (
	"gorm.io/gorm"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
)

// RoleMenuRepository database structure
type RoleMenuRepository struct {
	db     lib.Database
	logger lib.Logger
}

// WithTrx enables repository with transaction
func (r RoleMenuRepository) WithTrx(trxHandle *gorm.DB) RoleMenuRepository {
	if trxHandle == nil {
		r.logger.Zap.Error("Transaction Database not found in echo context.")
		return r
	}

	r.db.ORM = trxHandle
	return r
}

func (r RoleMenuRepository) Query(param *models.RoleMenuQueryParam) (*models.RoleMenuQueryResult, error) {
	db := r.db.ORM.Model(&models.RoleMenu{})

	if v := param.RoleID; v != "" {
		db = db.Where("role_id = ?", v)
	}

	if v := param.RoleIDs; len(v) > 0 {
		db = db.Where("role_id IN (?)", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make([]*models.RoleMenu, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.RoleMenuQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (r RoleMenuRepository) Get(id string) (*models.RoleMenu, error) {
	roleMenu := new(models.RoleMenu)

	if ok, err := QueryOne(r.db.ORM.Model(roleMenu).Where("id = ?", id), roleMenu); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return roleMenu, nil
}

func (r RoleMenuRepository) Create(roleMenu *models.RoleMenu) error {
	result := r.db.ORM.Model(roleMenu).Create(roleMenu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r RoleMenuRepository) Update(id string, roleMenu *models.RoleMenu) error {
	result := r.db.ORM.Model(roleMenu).Where("id = ?", id).Updates(roleMenu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r RoleMenuRepository) Delete(id string) error {
	roleMenu := new(models.RoleMenu)

	result := r.db.ORM.Model(roleMenu).Where("id = ?", id).Delete(roleMenu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r RoleMenuRepository) DeleteByRoleID(id string) error {
	roleMenu := new(models.RoleMenu)

	result := r.db.ORM.Model(roleMenu).Where("role_id = ?", id).Delete(roleMenu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

// NewRoleMenuRepository creates a new role menu repository
func NewRoleMenuRepository(db lib.Database, logger lib.Logger) RoleMenuRepository {
	return RoleMenuRepository{
		db:     db,
		logger: logger,
	}
}
