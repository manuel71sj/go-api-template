package repository

import (
	"gorm.io/gorm"
	"manuel71sj/go-api-template/errors"
	"manuel71sj/go-api-template/lib"
	"manuel71sj/go-api-template/models"
)

// MenuRepository database structure
type MenuRepository struct {
	db     lib.Database
	logger lib.Logger
}

// WithTrx enables repository with transaction
func (m MenuRepository) WithTrx(trxHandle *gorm.DB) MenuRepository {
	if trxHandle == nil {
		m.logger.Zap.Error("Transaction Database not found in echo context.")
		return m
	}

	m.db.ORM = trxHandle
	return m
}

func (m MenuRepository) Query(param *models.MenuQueryParam) (*models.MenuQueryResult, error) {
	db := m.db.ORM.Model(&models.Menu{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Name; v != "" {
		db = db.Where("name = ?", v)
	}

	if v := param.ParentID; v != "" {
		db = db.Where("parent_id = ?", v)
	}

	if v := param.PrefixParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}

	if v := param.Hidden; v != 0 {
		//db = db.Where("show_status = ?", v)
		db = db.Where("hidden = ?", v)
	}

	if v := param.Status; v != 0 {
		db = db.Where("status = ?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR remark LIKE ?", v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Menus, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.MenuQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (m MenuRepository) Get(id string) (*models.Menu, error) {
	menu := new(models.Menu)

	if ok, err := QueryOne(m.db.ORM.Model(menu).Where("id = ?", id), menu); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return menu, nil
}

func (m MenuRepository) Create(menu *models.Menu) error {
	result := m.db.ORM.Model(menu).Create(menu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (m MenuRepository) Update(id string, menu *models.Menu) error {
	result := m.db.ORM.Model(menu).Where("id = ?", id).Updates(menu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (m MenuRepository) Delete(id string) error {
	menu := new(models.Menu)

	result := m.db.ORM.Model(menu).Where("id = ?", id).Delete(menu)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (m MenuRepository) UpdateStatus(id string, status int) error {
	menu := new(models.Menu)

	result := m.db.ORM.Model(menu).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (m MenuRepository) UpdateParentPath(id string, parentPath string) error {
	menu := new(models.Menu)

	result := m.db.ORM.Model(menu).Where("id = ?", id).Update("parent_path", parentPath)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

// NewMenuRepository creates a new menu repository
func NewMenuRepository(db lib.Database, logger lib.Logger) MenuRepository {
	return MenuRepository{
		db:     db,
		logger: logger,
	}
}
