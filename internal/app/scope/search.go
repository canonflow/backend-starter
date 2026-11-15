package scope

import "gorm.io/gorm"

func Search(column string, value any, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(map[string]any{
			column: value,
		})
	}
}
