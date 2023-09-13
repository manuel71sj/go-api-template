package database

import (
	"database/sql"
	"gorm.io/gorm"
)

type Model struct {
	RecordID  uint           `gorm:"column:record_id;primaryKey;autoIncrement;" json:"-"`
	CreatedAt sql.NullTime   `gorm:"column:created_at;autoCreateTime;" json:"created_at"`
	UpdatedAt sql.NullTime   `gorm:"column:updated_at;autoUpdateTime;" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;" json:"-"`
	Deleted   bool           `gorm:"column:deleted;default:false;" json:"deleted"`
}
