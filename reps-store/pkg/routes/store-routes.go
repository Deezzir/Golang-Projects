package routes

import (
	"github.com/gorilla/mux"
	"reps-store/pkg/controllers"
)

var RegisterRepsStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/repository/", controllers.CreateRepository).Methods("POST")
	router.HandleFunc("/repository/", controllers.GetRepository).Methods("GET")
	router.HandleFunc("/repository/{repID}", controllers.GetRepositoryByID).Methods("GET")
	router.HandleFunc("/repository/{repID}", controllers.UpdateRepository).Methods("PUT")
	router.HandleFunc("/repository/{repID}", controllers.DeleteRepository).Methods("DELETE")
}
