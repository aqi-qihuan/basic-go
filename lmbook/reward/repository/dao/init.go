package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	// TODO
	return db.AutoMigrate(&Reward{})
}
