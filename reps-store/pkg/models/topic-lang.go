package models

import "reps-store/pkg/utils"

type Topic struct {
	Name string `gorm:"primaryKey" json:"name"`
}

type Language struct {
	Name string `gorm:"primaryKey" json:"name"`
}

func GetAllTopics() []Topic {
	var topics []Topic
	result := DB.Find(&topics)
	if result.Error != nil {
		return nil
	} else {
		utils.InfoLogger.Printf("Retrieving (%d) Topic records from the DB\n", result.RowsAffected)
		return topics
	}
}

func GetAllLanguages() []Language {
	var languages []Language
	result := DB.Find(&languages)
	if result.Error != nil {
		return nil
	} else {
		utils.InfoLogger.Printf("Retrieving (%d) Language records from the DB\n", result.RowsAffected)
		return languages
	}
}
