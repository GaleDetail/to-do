package migrations

import (
	"gorm.io/gorm"
	"to-do/models"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		return err
	}
	return nil
}
