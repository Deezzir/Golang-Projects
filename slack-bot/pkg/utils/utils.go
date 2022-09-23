package utils

import (
	"log"
	"os"
	//"path/filepath"
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

func ListDir(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		ErrorLogger.Printf("Failed to read /%s directory\n", path)
	}

	var fileNames []string
	for _, file := range files {
		// fileName := filepath.Join(path, file.Name())
		fileNames = append(fileNames, file.Name())
	}

	return fileNames
}
