package main

import (
	"github.com/gorilla/mux"

	"net/http"
	"reps-store/pkg/routes"
	"reps-store/pkg/utils"
)

// ======================================
// Main
// ======================================
var server_port = "0.0.0.0:8080"

func main() {
	r := mux.NewRouter()
	routes.RegisterRepsStoreRoutes(r)

	http.Handle("/", r)

	utils.InfoLogger.Println("Starting server at", server_port)
	if err := http.ListenAndServe(server_port, r); err != nil {
		utils.ErrLogger.Fatalln(err)
	}
}
