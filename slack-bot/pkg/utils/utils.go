package utils

import (
	"log"
	"os"
)

var (
	ErrorLogger   *log.Logger
	InfoLogger    *log.Logger
	CommandLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
	CommandLogger = log.New(os.Stdout, "[COMNMAND EVENT]: ", 0)
}
