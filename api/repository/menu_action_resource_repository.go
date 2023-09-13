package repository

import (
	"gorm.io/gorm"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
)

// MenuActionResourceRepository database structure
type MenuActionResourceRepository struct {
	db     lib.Database
	logger lib.Logger
}

// WithTrx enables repository with transaction
func (r MenuActionResourceRepository) WithTrx(trxHandle *gorm.DB) MenuActionResourceRepository {
	if trxHandle == nil {
		r.logger.Zap.Error("Transaction Database not found in echo context.")
		return r
	}

	r.db.ORM = trxHandle
	return r
}

func (r MenuActionResourceRepository) Query(param *models.MenuActionResourceQueryParam) (*models.MenuActionResourceQueryResult, error) {
	db := r.db.ORM.Model(&models.MenuActionResource{})

	if v := param.MenuID; v != "" {
		subQuery := r.db.ORM.Model(&models.MenuAction{}).
			Where("menu_id = ?", v).
			Select("id")

		db = db.Where("action_id IN (?)", subQuery)
	}

	if v := param.MenuIDs; len(v) > 0 {
		subQuery := r.db.ORM.Model(&models.MenuAction{}).
			Where("menu_id IN (?)", v).
			Select("id")

		db = db.Where("action_id IN (?)", subQuery)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.MenuActionResources, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.MenuActionResourceQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (r MenuActionResourceRepository) Get(id string) (*models.MenuActionResource, error) {
	menuActionResource := new(models.MenuActionResource)

	if ok, err := QueryOne(r.db.ORM.Model(menuActionResource).Where("id = ?", id), menuActionResource); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return menuActionResource, nil
}

func (r MenuActionResourceRepository) Create(menuActionResource *models.MenuActionResource) error {
	result := r.db.ORM.Model(menuActionResource).Create(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r MenuActionResourceRepository) Update(id string, menuActionResource *models.MenuActionResource) error {
	result := r.db.ORM.Model(menuActionResource).Where("id=?", id).Updates(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r MenuActionResourceRepository) Delete(id string) error {
	menuActionResource := new(models.MenuActionResource)

	result := r.db.ORM.Model(menuActionResource).Where("id=?", id).Delete(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r MenuActionResourceRepository) DeleteByActionID(actionID string) error {
	menuActionResource := new(models.MenuActionResource)

	result := r.db.ORM.Model(menuActionResource).Where("action_id=?", actionID).Delete(menuActionResource)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (r MenuActionResourceRepository) DeleteByMenuID(menuID string) error {
	menuAction := new(models.MenuAction)
	menuActionResource := new(models.MenuActionResource)

	subQuery := r.db.ORM.Model(menuAction).
		Where("menu_id=?", menuID).Select("id")

	result := r.db.ORM.Model(menuActionResource).
		Where("action_id IN (?)", subQuery).Delete(menuActionResource)

	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

// NewMenuActionResourceRepository creates a new menu action resource repository
func NewMenuActionResourceRepository(db lib.Database, logger lib.Logger) MenuActionResourceRepository {
	return MenuActionResourceRepository{
		db:     db,
		logger: logger,
	}
}
