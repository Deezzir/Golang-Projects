package models

import (
// "log"

// "gorm.io/gorm"
)

type Programmer struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Login string `gorm:"primaryKey,unique" json:"login"`
	Bio   string `json:"bio"`
	Email string `json:"email"`
}
