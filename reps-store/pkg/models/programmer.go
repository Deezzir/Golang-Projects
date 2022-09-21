package models

import (
	"errors"
	"reps-store/pkg/utils"

	"gorm.io/gorm"
)

type Programmer struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Login string `gorm:"primaryKey;not null" json:"login"`
	Bio   string `json:"bio"`
	Email string `gorm:"not null;unique" json:"email"`
}

type repository struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	License     string `json:"license"`
}

func (p *Programmer) CreateProgrammer() *Programmer {
	result := DB.Create(&p)
	if result.Error != nil {
		return nil
	} else {
		utils.InfoLogger.Printf("Creating a new Programmer record with ID(%d) in the DB\n", p.ID)
		return p
	}
}

func GetAllProgrammers() []Programmer {
	var programmers []Programmer
	result := DB.Model(&Programmer{}).Find(&programmers)
	if result.Error != nil {
		return nil
	} else {
		utils.InfoLogger.Printf("Retrieving (%d) Programmer records from the DB\n", result.RowsAffected)
		return programmers
	}
}

func GetProgrammerByIDFull(id uint) interface{} {
	var programmer Programmer

	record := struct {
		Programmer Programmer   `json:"programmer"`
		Reps       []repository `json:"repositories"`
	}{}

	result := DB.Model(&Programmer{}).First(&programmer, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		utils.InfoLogger.Printf("Querying a Programmer record with ID(%d) from the DB", id)
		record.Programmer = programmer
		reps := GetRepositoriesByProgrammerID(id)
		for _, rep := range reps {
			record.Reps = append(record.Reps, repository{
				ID:          rep.ID,
				Name:        rep.Name,
				Description: rep.Description,
				License:     rep.License,
			})
		}
	}
	return &record
}

func GetProgrammerByID(id uint) (*Programmer, *gorm.DB) {
	var programmer Programmer

	result := DB.Model(&Programmer{}).First(&programmer, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result
	} else {
		utils.InfoLogger.Printf("Querying a Programmer record with ID(%d) from the DB\n", id)
		return &programmer, result
	}
}

func DeleteProgrammer(id uint) (*Programmer, bool) {
	var programmer Programmer

	result := DB.Model(&Programmer{}).First(&programmer, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, true
	} else {
		repositories := GetRepositoriesByProgrammerID(id)
		for _, repository := range repositories {
			_, ok := DeleteRepository(repository.ID)
			if !ok {
				return nil, false
			}
		}

		result = DB.Delete(&programmer)
		if result.Error != nil {
			return nil, false
		} else {
			utils.InfoLogger.Printf("Deleting a Programmer record with ID(%d) from the DB\n", id)
			return &programmer, true
		}
	}
}
