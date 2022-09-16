package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// ======================================
// Variables and Structs
// ======================================

type Repository struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	License     string      `json:"license"`
	Topics      []string    `json:"topics"`
	Languages   []string    `json:"languages"`
	Programmer  *Programmer `json:"programmer"`
}

type Programmer struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Bio   string `json:"bio"`
	Email string `json:"email"`
}

var (
	reps []Repository
	id   int64 = 1
)

// ======================================
// Utils
// ======================================

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

func findMaxID(rs []Repository) int64 {
	var id int64 = 0
	for _, rep := range rs {
		if rep.ID > id {
			id = rep.ID
		}
	}

	return id
}

func getRepository(id int64) (repo *Repository, index int) {
	for i, rep := range reps {
		if rep.ID == id {
			return &rep, i
		}
	}

	return nil, -1
}

// ======================================
// Handlers
// ======================================

func getRepositories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reps); err != nil {
		log.Printf("[ERROR]: Failed to encode repositories: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Printf("[INFO]: Sending %d Repository objects\n", len(reps))
	}

}

func deleteRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	repIDstr := vars["repID"]

	repID, err := strconv.ParseInt(repIDstr, 0, 0)
	if err != nil {
		log.Printf("[ERROR]: Failed to parse string '%s' into RepID integer\n", repIDstr)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		rep, index := getRepository(repID)
		if index != -1 {
			reps = append(reps[:index], reps[index+1:]...)
			if err := json.NewEncoder(w).Encode(rep); err != nil {
				log.Printf("[ERROR]: Failed to encode repository: %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				log.Printf("[INFO]: Sending deleted Repository object with ID(%d)\n", repID)
			}
			return
		}
		log.Printf("[WARNING]: A repository object with ID(%d) not found\n", repID)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getRepositoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	repIDstr := vars["repID"]

	repID, err := strconv.ParseInt(repIDstr, 0, 0)
	if err != nil {
		log.Printf("[ERROR]: Failed to parse string '%s' into RepID integer\n", repIDstr)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		rep, _ := getRepository(repID)
		if rep != nil {
			if err := json.NewEncoder(w).Encode(rep); err != nil {
				log.Printf("[ERROR]: Failed to encode repository: %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				log.Printf("[INFO]: Sending queried Repository object with ID(%d)\n", repID)
			}
			return
		}
		log.Printf("[WARNING]: A repository object with ID(%d) not found\n", repID)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func createRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rep Repository
	if err := json.NewDecoder(r.Body).Decode(&rep); err != nil {
		log.Printf("[ERROR]: Failed to decode request body into a Repository object: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		id++
		rep.ID = id

		reps = append(reps, rep)
		if err := json.NewEncoder(w).Encode(rep); err != nil {
			log.Printf("[ERROR]: Failed to encode repository: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Printf("[INFO]: Sending created Repository object with ID(%d)\n", id)
		}
	}

}

func updateRepository(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	repIDstr := vars["repID"]

	repID, err := strconv.ParseInt(repIDstr, 0, 0)
	if err != nil {
		log.Printf("[ERROR]: Failed to parse string '%s' into RepID integer\n", repIDstr)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		rep, index := getRepository(repID)
		if index != -1 {
			reps = append(reps[:index], reps[index+1:]...)
			if err := json.NewDecoder(r.Body).Decode(&rep); err != nil {
				log.Printf("[ERROR]: Failed to decode request body into a Repository object: %s\n", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				rep.ID = repID

				reps = append(reps, *rep)
				if err := json.NewEncoder(w).Encode(rep); err != nil {
					log.Printf("[ERROR]: Failed to encode repository: %s\n", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					log.Printf("[INFO]: Sending updated Repository object with ID(%d)\n", repID)
				}
			}
		}
	}
}

// ======================================
// Main
// ======================================
var server_port = "0.0.0.0:8000"

func main() {
	reps = readJSON("mock_data.json")
	id = findMaxID(reps)

	r := mux.NewRouter()
	r.HandleFunc("/repositories", getRepositories).Methods("GET")
	r.HandleFunc("/repositories/{repID}", getRepositoryByID).Methods("GET")
	r.HandleFunc("/repositories", createRepository).Methods("POST")
	r.HandleFunc("/repositories/{repID}", updateRepository).Methods("PUT")
	r.HandleFunc("/repositories/{repID}", deleteRepository).Methods("DELETE")

	log.Println("[INFO]: Starting server at", server_port)
	if err := http.ListenAndServe(server_port, r); err != nil {
		log.Fatalln(err)
	}
}
