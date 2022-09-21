package models

type Topic struct {
	Name string `gorm:"primaryKey" json:"name"`
}

type Language struct {
	Name string `gorm:"primaryKey" json:"name"`
}
