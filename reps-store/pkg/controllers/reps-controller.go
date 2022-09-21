package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"reps-store/pkg/models"
	"reps-store/pkg/utils"
)

var Repository models.Repository

func GetRepositories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	repositories := models.GetAllRepositories()
	res, err := json.Marshal(repositories)
	if err != nil {
		utils.ErrLogger.Printf("Failed to encode the Repository records into JSON - %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		utils.InfoLogger.Println("Sent all Repository records ('GET':'/repository')")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
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
				utils.ErrLogger.Printf("Failed to encode the Repository record into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
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
		rep := repository.CreateRepository()
		if rep == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			res, err := json.Marshal(rep)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode the Repository record into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				utils.InfoLogger.Printf("Created a Repository record with ID(%d) ('POST':'/repository')\n", rep.ID)
				w.WriteHeader(http.StatusOK)
				w.Write(res)
			}
		}
	} else {
		utils.ErrLogger.Println("Failed to decode the Repository record from JSON")
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			res, err := json.Marshal(reps)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode the Repository records into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				utils.InfoLogger.Println("Uploaded Repository records ('POST':'/repository/upload')")
				w.WriteHeader(http.StatusOK)
				w.Write(res)
			}
		}
	} else {
		utils.ErrLogger.Println("Failed to decode the Repository records from JSON")
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
		repository := models.DeleteRepository(uint(repID))
		if repository == nil {
			utils.InfoLogger.Printf("A Repository record with ID(%d) not found ('DELETE':'/repository/%d')\n", repID, repID)
			w.WriteHeader(http.StatusNotFound)
			w.Write(utils.NotFoundMsg)
		} else {
			res, err := json.Marshal(repository)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode the Repository record into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
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

	newRepository := &models.Repository{}
	utils.ParseBody(r, newRepository)

	vars := mux.Vars(r)
	repIDstr := vars["repID"]

	repID, err := strconv.ParseUint(repIDstr, 0, 0)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse string '%s' into RepID integer\n", repIDstr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.BadRequestMsg)
	} else {
		repository, db := models.GetRepositoryByID(uint(repID))
		updateRepository(newRepository, repository)

		db.Save(&repository)
		res, err := json.Marshal(repository)
		if err != nil {
			utils.ErrLogger.Printf("Failed to encode the Repository record into JSON - %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			utils.InfoLogger.Printf("Updated a Repository record with ID(%d) ('PUT':'/repository/%d')\n", repID, repID)
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		}
	}
}

func updateRepository(newRep, rep *models.Repository) {
	if newRep.Name != "" {
		rep.Name = newRep.Name
	}

	if newRep.Description != "" {
		rep.Description = newRep.Description
	}

	if newRep.License != "" {
		rep.License = newRep.License
	}

	if len(newRep.Topics) == 0 || newRep.Topics == nil {
		rep.Topics = []*models.Topic{}
	} else {
		rep.Topics = newRep.Topics
	}

	if len(newRep.Languages) == 0 || newRep.Languages == nil {
		rep.Languages = []*models.Language{}
	} else {
		rep.Languages = newRep.Languages
	}
}
