package repository

import (
	"gorm.io/gorm"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
)

// RoleRepository database structure
type RoleRepository struct {
	db     lib.Database
	logger lib.Logger
}

// WithTrx enables repository with transaction
func (r RoleRepository) WithTrx(trxHandle *gorm.DB) RoleRepository {
	if trxHandle == nil {
		r.logger.Zap.Error("Transaction Database not found in echo context.")
		return r
	}

	r.db.ORM = trxHandle
	return r
}

func (r RoleRepository) Query(param *models.RoleQueryParam) (*models.RoleQueryResult, error) {
	db := r.db.ORM.Model(&models.Role{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Name; v != "" {
		db = db.Where("name = ?", v)
	}

	if v := param.UserID; v != "" {
		subQuery := r.db.ORM.Model(&models.UserRole{}).
			Where("user_id = ?", v).
			Select("role_id")

		db = db.Where("id IN (?)", subQuery)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR remark LIKE ?", v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Roles, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.RoleQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (r RoleRepository) Get(id string) (*models.Role, error) {
	role := new(models.Role)

	if ok, err := QueryOne(r.db.ORM.Model(role).Where("id = ?", id), role); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return role, nil
}

func (r RoleRepository) Create(role *models.Role) error {
	result := r.db.ORM.Model(role).Create(role)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r RoleRepository) Update(id string, role *models.Role) error {
	result := r.db.ORM.Model(role).Where("id = ?", id).Updates(role)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r RoleRepository) Delete(id string) error {
	role := new(models.Role)

	result := r.db.ORM.Model(role).Where("id = ?", id).Delete(role)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r RoleRepository) UpdateStatus(id string, status int) error {
	role := new(models.Role)

	result := r.db.ORM.Model(role).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db lib.Database, logger lib.Logger) RoleRepository {
	return RoleRepository{
		db:     db,
		logger: logger,
	}
}
