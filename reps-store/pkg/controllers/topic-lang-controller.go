package controllers

import (
	"encoding/json"
	"net/http"
	"reps-store/pkg/models"
	"reps-store/pkg/utils"
)

var Topic models.Topic
var Language models.Language

func GetTopics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	topics := models.GetAllTopics()
	if topics != nil {
		res, err := json.Marshal(topics)
		if err != nil {
			utils.ErrLogger.Printf("Failed to encode the Topic records into JSON - %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			utils.InfoLogger.Println("Sent all Topic records ('GET':'/topics')")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		}
	} else {
		utils.ErrLogger.Println("Failed to retrieve Topic records from the DB")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ServerErrorMsg)
	}
}

func GetLanguages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	languages := models.GetAllLanguages()
	if languages != nil {
		res, err := json.Marshal(languages)
		if err != nil {
			utils.ErrLogger.Printf("Failed to encode the Language records into JSON - %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.ServerErrorMsg)
		} else {
			utils.InfoLogger.Println("Sent all Language records ('GET':'/languages')")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		}
	} else {
		utils.ErrLogger.Println("Failed to retrieve Language records from the DB")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ServerErrorMsg)
	}
}
