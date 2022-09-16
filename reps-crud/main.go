package main

import (
	"encoding/json"
	"log"
	// 	"math/random"
	// 	"net/http"
	"os"
	// 	"strconv"
	// "github.com/gorilla/mux"
)

type Repository struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	License     string      `json:"license"`
	Topics      []string    `json:"topics"`
	Languages   []string    `json:"languages"`
	Programmer  *Programmer `json:"programmer"`
}

type Programmer struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Bio   string `json:"bio"`
	Email string `json:"email"`
}

var reps []Repository

func readJSON(filename string) []Repository {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("[ERROR]: failed to open the provided file - '%s'\n", filename)
	}

	var reps []Repository
	if err := json.NewDecoder(file).Decode(&reps); err != nil {
		log.Fatalf("[ERROR]: Failed to decode repositories in the provided file - '%s'\n%s\n", filename, err.Error())
	}
	return reps
}

func main() {
	reps = readJSON("mock_data.json")
}
