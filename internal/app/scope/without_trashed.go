package scope

import "gorm.io/gorm"

func WithoutTrashed(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}
