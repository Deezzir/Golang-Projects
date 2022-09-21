package utils

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

var ErrLogger *log.Logger
var WarnLogger *log.Logger
var InfoLogger *log.Logger

var NotFoundMsg []byte
var InvalidJSONMsg []byte
var ServerErrorMsg []byte
var BadRequestMsg []byte

func ParseBody(r *http.Request, obj interface{}) bool {
	if body, err := io.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), obj); err != nil {
			ErrLogger.Printf("Failed to decode JSON from request Body - %s\n", err.Error())
			return false
		}
	}
	return true
}

// func readJSON(filename string) []interface{} {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatalf("[ERROR]: failed to open the provided file - '%s'\n", filename)
// 	}

// 	var objects []interface{}
// 	if err := json.NewDecoder(file).Decode(&objects); err != nil {
// 		log.Fatalf("[ERROR]: Failed to decode JSON in the provided file - '%s'\n%s\n", filename, err.Error())
// 	}
// 	return objects
// }

func createMsgs() {
	if result, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Not Found",
	}); err != nil {
		ErrLogger.Fatalf("Failed to encode 'NotFoundMsg' into JSON - %s\n", err.Error())
	} else {
		NotFoundMsg = result
	}

	if result, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Invalid JSON",
	}); err != nil {
		ErrLogger.Fatalf("Failed to encode 'InvalidJSONMsg' into JSON - %s\n", err.Error())
	} else {
		InvalidJSONMsg = result
	}

	if result, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Internal Server Error Occured",
	}); err != nil {
		ErrLogger.Fatalf("Failed to encode 'ServerErrorMsg' into JSON - %s\n", err.Error())
	} else {
		ServerErrorMsg = result
	}

	if result, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: "Bad Request",
	}); err != nil {
		ErrLogger.Fatalf("Failed to encode 'BadRequestMsg' into JSON - %s\n", err.Error())
	} else {
		BadRequestMsg = result
	}
}

func init() {
	logFlags := log.LstdFlags | log.Lshortfile
	InfoLogger = log.New(os.Stdout, "[INFO] ", logFlags)
	ErrLogger = log.New(os.Stderr, "[ERROR] ", logFlags)
	WarnLogger = log.New(os.Stdout, "[WARNING] ", logFlags)

	createMsgs()
}
