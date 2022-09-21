package models

import (
	"reps-store/pkg/config"

	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	config.SetDataBase()
	DB = config.GetDataBase()
	DB.AutoMigrate(&Repository{})
}
