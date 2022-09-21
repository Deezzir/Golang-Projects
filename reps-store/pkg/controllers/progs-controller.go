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

var Programmer models.Programmer

func GetProgrammers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	programmers := models.GetAllProgrammers()
	if programmers != nil {
		res, err := json.Marshal(programmers)
		if err != nil {
			utils.ErrLogger.Printf("Failed to encode the Programmer records into JSON - %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			utils.InfoLogger.Println("Sent all Programmer records ('GET':'/programmer')")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		}
	} else {
		utils.ErrLogger.Println("Failed to retrieve Programmer records from the DB")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ServerErrorMsg)
	}
}

func GetProgrammerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	progIDstr := vars["progID"]

	progID, err := strconv.ParseUint(progIDstr, 0, 0)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse string '%s' into ProgID integer\n", progIDstr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.BadRequestMsg)
	} else {
		programmer := models.GetProgrammerByIDFull(uint(progID))
		if programmer == nil {
			utils.InfoLogger.Printf("A Programmer record with ID(%d) not found ('GET':'/programmer/%d')\n", progID, progID)
			w.WriteHeader(http.StatusNotFound)
			w.Write(utils.NotFoundMsg)
		} else {
			res, err := json.Marshal(programmer)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode a Programmer record into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(utils.ServerErrorMsg)
			} else {
				utils.InfoLogger.Printf("Sent a Programmer record with ID(%d) ('GET':'/programmer%d')\n", progID, progID)
				w.WriteHeader(http.StatusOK)
				w.Write(res)
			}
		}
	}
}

func CreateProgrammer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	programmer := &models.Programmer{}
	if ok := utils.ParseBody(r, programmer); ok {
		prog := programmer.CreateProgrammer()
		if prog == nil {
			utils.ErrLogger.Println("Failed to create a new Programmer record in the DB")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			res, err := json.Marshal(prog)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode a Programmer record into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(utils.ServerErrorMsg)
			} else {
				utils.InfoLogger.Printf("Created a Programmer record with ID(%d) ('POST':'/programmer')\n", prog.ID)
				w.WriteHeader(http.StatusOK)
				w.Write(res)
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.InvalidJSONMsg)
	}
}

func DeleteProgrammer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	progIDstr := vars["progID"]

	progID, err := strconv.ParseUint(progIDstr, 0, 0)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse string '%s' into ProgID integer\n", progIDstr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.BadRequestMsg)
	} else {
		programmer, ok := models.DeleteProgrammer(uint(progID))
		if programmer == nil && ok {
			utils.InfoLogger.Printf("A Programmer record with ID(%d) not found ('DELETE':'/programmer/%d')\n", progID, progID)
			w.WriteHeader(http.StatusNotFound)
			w.Write(utils.NotFoundMsg)
		} else if programmer == nil && !ok {
			utils.InfoLogger.Printf("Failed to delete a Programmer record with ID(%d)\n", progID)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			res, err := json.Marshal(programmer)
			if err != nil {
				utils.ErrLogger.Printf("Failed to encode a Programmer record into JSON - %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(utils.ServerErrorMsg)
			} else {
				utils.InfoLogger.Printf("Deleted a Programmer record with ID(%d) ('DELETE':'/programmer/%d')\n", progID, progID)
				w.WriteHeader(http.StatusOK)
				w.Write(res)
			}
		}
	}
}

func UpdateProgrammer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	progIDstr := vars["progID"]

	progID, err := strconv.ParseUint(progIDstr, 0, 0)
	if err != nil {
		utils.ErrLogger.Printf("Failed to parse string '%s' into ProgID integer\n", progIDstr)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.BadRequestMsg)
	} else {
		newProgrammer := &models.Programmer{}

		if ok := utils.ParseBody(r, newProgrammer); ok {
			programmer, db := models.GetProgrammerByID(uint(progID))
			if programmer == nil {
				utils.InfoLogger.Printf("A Programmer record with ID(%d) not found ('PUT':'/programmer/%d')\n", progID, progID)
				w.WriteHeader(http.StatusNotFound)
				w.Write(utils.NotFoundMsg)
			} else {
				if ok := updateProgrammer(newProgrammer, programmer, db); ok {
					db.Save(&programmer)
					res, err := json.Marshal(programmer)
					if err != nil {
						utils.ErrLogger.Printf("Failed to encode the Programmer record into JSON - %s\n", err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						w.Write(utils.ServerErrorMsg)
					} else {
						utils.InfoLogger.Printf("Updated a Programmer record with ID(%d) ('PUT':'/programmer/%d')\n", progID, progID)
						w.WriteHeader(http.StatusOK)
						w.Write(res)
					}
				} else {
					utils.ErrLogger.Printf("Failed to update a Programmer record with ID(%d) ('PUT':'/programmer/%d')\n", progID, progID)
					w.WriteHeader(http.StatusBadRequest)
					w.Write(utils.InvalidJSONMsg)
				}
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.InvalidJSONMsg)
		}
	}
}

func updateProgrammer(newProg, prog *models.Programmer, db *gorm.DB) bool {
	if newProg.Bio != "" {
		prog.Bio = newProg.Bio
	}

	if newProg.Email != "" {
		prog.Email = newProg.Email
	}

	return true
}
