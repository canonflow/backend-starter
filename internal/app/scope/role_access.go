package scope

import "gorm.io/gorm"

func RoleWithPermission(db *gorm.DB) *gorm.DB {
	return db.Preload("Permissions").
		Preload("Permissions.Actions", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "name")
		})
}

func PermissionWithAction(db *gorm.DB) *gorm.DB {
	return db.Preload("Actions")
}

func ActionWithPermission(db *gorm.DB) *gorm.DB {
	return db.Preload("Permissions")
}
