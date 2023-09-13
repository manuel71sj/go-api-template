package repository

import (
	"gorm.io/gorm"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
)

// UserRoleRepository database structure
type UserRoleRepository struct {
	db     lib.Database
	logger lib.Logger
}

// WithTrx enables repository with transaction
func (r UserRoleRepository) WithTrx(trxHandle *gorm.DB) UserRoleRepository {
	if trxHandle == nil {
		r.logger.Zap.Error("Transaction Database not found in echo context.")
		return r
	}

	r.db.ORM = trxHandle
	return r
}

func (r UserRoleRepository) Query(param *models.UserRoleQueryParam) (*models.UserRoleQueryResult, error) {
	db := r.db.ORM.Model(models.UserRole{})

	if v := param.UserID; v != "" {
		db = db.Where("user_id = ?", v)
	}

	if v := param.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.UserRoles, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.UserRoleQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (r UserRoleRepository) Get(id string) (*models.UserRole, error) {
	userRole := new(models.UserRole)

	if ok, err := QueryOne(r.db.ORM.Model(userRole).Where("id = ?", id), userRole); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return userRole, nil
}

func (r UserRoleRepository) Create(userRole *models.UserRole) error {
	result := r.db.ORM.Model(userRole).Create(userRole)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r UserRoleRepository) Update(id string, userRole *models.UserRole) error {
	result := r.db.ORM.Model(userRole).Where("id = ?", id).Updates(userRole)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r UserRoleRepository) Delete(id string) error {
	userRole := new(models.UserRole)

	result := r.db.ORM.Model(userRole).Where("id = ?", id).Delete(userRole)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r UserRoleRepository) DeleteByUserID(userID string) error {
	userRole := new(models.UserRole)

	result := r.db.ORM.Model(userRole).Where("user_id = ?", userID).Delete(userRole)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

// NewUserRoleRepository creates a new user role repository
func NewUserRoleRepository(db lib.Database, logger lib.Logger) UserRoleRepository {
	return UserRoleRepository{
		db:     db,
		logger: logger,
	}
}
