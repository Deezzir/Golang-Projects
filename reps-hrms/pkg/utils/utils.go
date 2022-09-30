package utils

import (
	"log"
	"os"
)

var (
	ErrLogger  *log.Logger
	InfoLogger *log.Logger
	WarnLogger *log.Logger
)

type Response struct {
	Message string `json:"message"`
}

var (
	NotFoundMsg    Response
	InvalidJSONMsg Response
	ServerErrorMsg Response
	BadRequestMsg  Response
	OKMsg          Response
)

func init() {
	logFlags := log.LstdFlags | log.Lshortfile
	ErrLogger = log.New(os.Stderr, "[ERROR] ", logFlags)
	InfoLogger = log.New(os.Stdout, "[INFO] ", logFlags)
	WarnLogger = log.New(os.Stdout, "[WARN] ", logFlags)

	createMsgs()
}

func createMsgs() {
	NotFoundMsg = Response{
		Message: "Not Found",
	}
	InvalidJSONMsg = Response{
		Message: "Invalid JSON",
	}
	ServerErrorMsg = Response{
		Message: "Internal Server Error",
	}
	BadRequestMsg = Response{
		Message: "Bad Request",
	}
	OKMsg = Response{
		Message: "OK",
	}
}
