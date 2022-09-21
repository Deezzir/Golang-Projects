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
	Topics      []*Topic    `gorm:"many2many:reps_topics;ForeignKey:id;References:name" json:"topics"`
	Languages   []*Language `gorm:"many2many:reps_langs;ForeignKey:id;References:name"  json:"languages"`

	ProgrammerID uint        `json:"programmerID"`
	Programmer   *Programmer `gorm:"constraint:OnUpdate:CASCADE;" json:"programmer"`
}

func (r *Repository) BeforeDelete(tx *gorm.DB) (err error) {
	if err := tx.Model(&r).Association("Topics").Clear(); err != nil {
		return err
	}
	if err := tx.Model(&r).Association("Languages").Clear(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) BeforeCreate(tx *gorm.DB) (err error) {
	if programmer, _ := GetProgrammerByID(r.ProgrammerID); programmer == nil {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func GetRepositoriesByProgrammerID(id uint) []Repository {
	var repositories []Repository

	result := DB.Model(&Repository{}).Where("programmer_id = ?", id).Find(&repositories)
	if result.Error != nil {
		utils.ErrLogger.Printf("Failed to retrieve Repository records with ProgrammerID(%d) from the DB\n", id)
		return nil
	} else {
		utils.InfoLogger.Printf("Retrieving (%d) Repository records with ProgrammerID(%d) from the DB\n", result.RowsAffected, id)
	}
	return repositories
}

func (r *Repository) CreateRepository() (*Repository, bool) {
	result := DB.Create(&r)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, true
	} else if result.Error != nil {
		return nil, false
	} else {
		utils.InfoLogger.Printf("Creating a new Repository record with ID(%d) in the DB\n", r.ID)
		return r, true
	}
}

func UploadRepositories(reps []Repository) []Repository {
	result := DB.Create(&reps)
	if result.Error != nil {
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
		return nil
	} else {
		utils.InfoLogger.Printf("Retrieving (%d) Repository records from the DB\n", result.RowsAffected)
		return repositories
	}
}

func GetRepositoryByID(id uint) (*Repository, *gorm.DB) {
	var repository Repository

	result := DB.Model(&Repository{}).Preload("Languages").Preload("Topics").Preload("Programmer").First(&repository, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result
	} else {
		utils.InfoLogger.Printf("Querying a Repository record with ID(%d) from the DB\n", id)
		return &repository, result
	}
}

func DeleteRepository(id uint) (*Repository, bool) {
	var repository Repository

	result := DB.Model(&Repository{}).Preload("Languages").Preload("Topics").First(&repository, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, true
	} else {
		result := DB.Delete(&repository)
		if result.Error != nil {
			return nil, false
		}
		utils.InfoLogger.Printf("Deleting a Repository record with ID(%d) from the DB\n", id)
		repository.Programmer = nil
		return &repository, true
	}
}
