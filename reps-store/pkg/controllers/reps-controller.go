package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"reps-store/pkg/models"
	"reps-store/pkg/utils"
)

var Repository models.Repository

func GetRepositories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repositories := models.GetAllRepositories()
	if repositories != nil {
		res, err := json.Marshal(repositories)
		if err != nil {
			utils.ErrLogger.Printf("Failed to encode the Repository records into JSON - %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			utils.InfoLogger.Println("Sent all Repository records ('GET':'/repository')")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		}
	} else {
		utils.ErrLogger.Println("Failed to retrieve all Repository records from DB")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ServerErrorMsg)
	}

}

func GetRepositoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	repIDstr := vars["repID"]

	repID, err := strconv.ParseUint(repIDstr, 0, 0)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse string '%s' into RepID integer\n", repIDstr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.BadRequestMsg)
	} else {
		repository, _ := models.GetRepositoryByID(uint(repID))
		if repository == nil {
			utils.InfoLogger.Printf("A Repository record with ID(%d) not found ('GET':'/repository/%d')\n", repID, repID)
			w.WriteHeader(http.StatusNotFound)
			w.Write(utils.NotFoundMsg)
		} else {
			res, err := json.Marshal(repository)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode a Repository record into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(utils.ServerErrorMsg)
			} else {
				utils.InfoLogger.Printf("Sent a Repository record with ID(%d) ('GET':'/repository%d')\n", repID, repID)
				w.WriteHeader(http.StatusOK)
				w.Write(res)
			}
		}
	}
}

func CreateRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repository := &models.Repository{}
	if ok := utils.ParseBody(r, repository); ok {
		rep, ok := repository.CreateRepository()
		if rep == nil && ok {
			utils.ErrLogger.Printf("Provided ProgrammerID(%d) is invalid for creating a new Repository record ('POST':'/repository')\n", repository.ProgrammerID)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.BadRequestMsg)
		} else if rep == nil && !ok {
			utils.ErrLogger.Println("Failed to create a new Repository record in DB")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			res, err := json.Marshal(rep)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode a Repository record into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(utils.ServerErrorMsg)
			} else {
				utils.InfoLogger.Printf("Created a Repository record with ID(%d) ('POST':'/repository')\n", rep.ID)
				w.WriteHeader(http.StatusOK)
				w.Write(res)
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.InvalidJSONMsg)
	}
}

func UploadRepositories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repositories := []models.Repository{}
	if ok := utils.ParseBody(r, &repositories); ok {
		reps := models.UploadRepositories(repositories)
		if reps == nil {
			utils.ErrLogger.Println("Failed to upload Repository records into DB")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			res, err := json.Marshal(reps)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode the Repository records into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(utils.ServerErrorMsg)
			} else {
				utils.InfoLogger.Println("Uploaded Repository records ('POST':'/repository/upload')")
				w.WriteHeader(http.StatusOK)
				w.Write(res)
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.InvalidJSONMsg)
	}
}

func DeleteRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	repIDstr := vars["repID"]

	repID, err := strconv.ParseUint(repIDstr, 0, 0)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse string '%s' into RepID integer\n", repIDstr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.BadRequestMsg)
	} else {
		repository, ok := models.DeleteRepository(uint(repID))
		if repository == nil && ok {
			utils.InfoLogger.Printf("A Repository record with ID(%d) not found ('DELETE':'/repository/%d')\n", repID, repID)
			w.WriteHeader(http.StatusNotFound)
			w.Write(utils.NotFoundMsg)
		} else if repository == nil && !ok {
			utils.InfoLogger.Printf("Failed to delete a Repository record with ID(%d)\n", repID)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			res, err := json.Marshal(repository)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode a Repository record into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(utils.ServerErrorMsg)
			} else {
				utils.InfoLogger.Printf("Deleted a Repository record with ID(%d) ('DELETE':'/repository/%d')\n", repID, repID)
				w.WriteHeader(http.StatusOK)
				w.Write(res)
			}
		}
	}
}

func UpdateRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	repIDstr := vars["repID"]

	repID, err := strconv.ParseUint(repIDstr, 0, 0)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse string '%s' into RepID integer\n", repIDstr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.BadRequestMsg)
	} else {
		newRepository := &models.Repository{}

		if ok := utils.ParseBody(r, newRepository); ok {
			repository, db := models.GetRepositoryByID(uint(repID))
			if repository == nil {
				utils.InfoLogger.Printf("A Repository record with ID(%d) not found ('PUT':'/repository/%d')\n", repID, repID)
				w.WriteHeader(http.StatusNotFound)
				w.Write(utils.NotFoundMsg)
			} else {
				if ok := updateRepository(newRepository, repository, db); ok {
					db.Save(&repository)
					res, err := json.Marshal(repository)
					if err != nil {
						utils.ErrLogger.Printf("Failed to encode the Repository record into JSON - %s\n", err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(utils.ServerErrorMsg)
					} else {
						utils.InfoLogger.Printf("Updated a Repository record with ID(%d) ('PUT':'/repository/%d')\n", repID, repID)
						w.WriteHeader(http.StatusOK)
						w.Write(res)
					}
				} else {
					utils.ErrLogger.Printf("Failed to update a Repository record with ID(%d) ('PUT':'/repository/%d')\n", repID, repID)
					w.WriteHeader(http.StatusBadRequest)
					w.Write(utils.BadRequestMsg)
				}
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.InvalidJSONMsg)
		}
	}
}

func updateRepository(newRep, rep *models.Repository, db *gorm.DB) bool {
	if newRep.Name != "" {
		rep.Name = newRep.Name
	}

	if newRep.Description != "" {
		rep.Description = newRep.Description
	}

	if newRep.License != "" {
		rep.License = newRep.License
	}

	if len(newRep.Topics) == 1 && newRep.Topics[0] == nil {
		models.DB.Model(rep).Association("Topics").Clear()
	} else if len(newRep.Topics) > 0 {
		for _, topic := range newRep.Topics {
			if topic.Name == "" {
				return false
			}
		}
		models.DB.Model(rep).Association("Topics").Replace(newRep.Topics)
	}

	if len(newRep.Languages) == 1 && newRep.Languages[0] == nil {
		models.DB.Model(rep).Association("Languages").Clear()
	} else if len(newRep.Languages) > 0 {
		for _, lang := range newRep.Languages {
			if lang.Name == "" {
				return false
			}
		}
		models.DB.Model(rep).Association("Languages").Append(newRep.Languages)
	}

	if newRep.ProgrammerID != 0 {
		programmer, _ := models.GetProgrammerByID(newRep.ProgrammerID)
		if programmer == nil {
			return false
		}
		models.DB.Model(&rep).Association("Programmer").Replace(programmer)
	}

	return true
}
