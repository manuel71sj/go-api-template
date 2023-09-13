package repository

import (
	"gorm.io/gorm"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
)

// MenuActionRepository database structure
type MenuActionRepository struct {
	db     lib.Database
	logger lib.Logger
}

// WithTrx enables repository with transaction
func (r MenuActionRepository) WithTrx(trxHandle *gorm.DB) MenuActionRepository {
	if trxHandle == nil {
		r.logger.Zap.Error("Transaction Database not found in echo context.")
		return r
	}

	r.db.ORM = trxHandle

	return r
}

func (r MenuActionRepository) Query(param *models.MenuActionQueryParam) (*models.MenuActionQueryResult, error) {
	db := r.db.ORM.Model(&models.MenuAction{})

	if v := param.MenuID; v != "" {
		db = db.Where("menu_id = ?", v)
	}

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.MenuActions, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.MenuActionQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (r MenuActionRepository) Get(id string) (*models.MenuAction, error) {
	menuAction := new(models.MenuAction)

	if ok, err := QueryOne(r.db.ORM.Model(menuAction).Where("id = ?", id), menuAction); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return menuAction, nil
}

func (r MenuActionRepository) Create(menuAction *models.MenuAction) error {
	result := r.db.ORM.Model(menuAction).Create(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r MenuActionRepository) Update(id string, menuAction *models.MenuAction) error {
	result := r.db.ORM.Model(menuAction).Where("id = ?", id).Updates(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r MenuActionRepository) Delete(id string) error {
	menuAction := new(models.MenuAction)

	result := r.db.ORM.Model(menuAction).Where("id = ?", id).Delete(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r MenuActionRepository) DeleteByMenuID(menuID string) error {
	menuAction := new(models.MenuAction)

	result := r.db.ORM.Model(menuAction).Where("menu_id = ?", menuID).Delete(menuAction)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

// NewMenuActionRepository creates a new menu action repository
func NewMenuActionRepository(db lib.Database, logger lib.Logger) MenuActionRepository {
	return MenuActionRepository{
		db:     db,
		logger: logger,
	}
}
