package routes

import (
	"github.com/gorilla/mux"
	"reps-store/pkg/controllers"
)

var RegisterRepsStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/repository/upload", controllers.UploadRepositories).Methods("POST")
	router.HandleFunc("/repository", controllers.GetRepositories).Methods("GET")
	router.HandleFunc("/repository", controllers.CreateRepository).Methods("POST")
	router.HandleFunc("/repository/{repID}", controllers.GetRepositoryByID).Methods("GET")
	router.HandleFunc("/repository/{repID}", controllers.UpdateRepository).Methods("PUT")
	router.HandleFunc("/repository/{repID}", controllers.DeleteRepository).Methods("DELETE")

	router.HandleFunc("/programmer", controllers.GetProgrammers).Methods("GET")
	router.HandleFunc("/programmer", controllers.CreateProgrammer).Methods("POST")
	router.HandleFunc("/programmer/{progID}", controllers.GetProgrammerByID).Methods("GET")
	router.HandleFunc("/programmer/{progID}", controllers.UpdateProgrammer).Methods("PUT")
	router.HandleFunc("/programmer/{progID}", controllers.DeleteProgrammer).Methods("DELETE")

	router.HandleFunc("/topics", controllers.GetTopics).Methods("GET")
	router.HandleFunc("/languages", controllers.GetLanguages).Methods("GET")
}
