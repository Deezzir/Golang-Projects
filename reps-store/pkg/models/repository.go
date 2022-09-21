package models

import (
	"errors"
	"reps-store/pkg/utils"

	"gorm.io/gorm"
)

type Repository struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `gorm:"primaryKey" json:"name"`
	Description string      `json:"description"`
	License     string      `json:"license"`
	Topics      []*Topic    `gorm:"many2many:reps_topics;" json:"topics"`
	Languages   []*Language `gorm:"many2many:reps_langs;"  json:"languages"`

	ProgrammerID uint       `json:"programmerID"`
	Programmer   Programmer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"programmer"`
}

func (r *Repository) CreateRepository() *Repository {
	result := DB.Create(&r)
	if result.Error != nil {
		utils.ErrLogger.Println("Failed to create a new Repository record in the DB")
		return nil
	} else {
		utils.InfoLogger.Printf("Creating a new Repository record with ID(%d) in the DB\n", r.ID)
		return r
	}
}

func UploadRepositories(reps []Repository) []Repository {
	result := DB.Create(&reps)
	if result.Error != nil {
		utils.ErrLogger.Println("Failed to upload Repository records to the DB")
		return nil
	} else {
		utils.InfoLogger.Printf("Uploading (%d) Repository record to the DB\n", len(reps))
		return reps
	}
}

func GetAllRepositories() []Repository {
	var repositories []Repository
	result := DB.Model(&Repository{}).Preload("Languages").Preload("Topics").Preload("Programmer").Find(&repositories)
	if result.Error != nil {
		utils.ErrLogger.Println("Failed to retrieve Repository records from the DB")
	} else {
		utils.InfoLogger.Printf("Retrieving (%d) Repository records from the DB\n", result.RowsAffected)
	}
	return repositories
}

func GetRepositoryByID(id uint) (*Repository, *gorm.DB) {
	var repository Repository

	utils.InfoLogger.Printf("Querying a Repository record with ID(%d) from the DB\n", id)
	result := DB.Model(&Repository{}).Preload("Languages").Preload("Topics").Preload("Programmer").First(&repository, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result
	} else {
		return &repository, result
	}
}

func DeleteRepository(id uint) *Repository {
	var repository Repository

	utils.InfoLogger.Printf("Deleting a Repository record with ID(%d) from the DB\n", id)
	result := DB.Delete(&repository, id)
	if result.Error != nil || repository.ID == 0 {
		return &repository
	} else {
		return &repository
	}

}
