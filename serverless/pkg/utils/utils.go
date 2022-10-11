package utils

import (
	"log"
	"os"
	"regexp"
)

var (
	InfoLogger *log.Logger
	ErrLogger  *log.Logger
	WarnLogger *log.Logger
)

type ErrorBody struct {
	Message string `json:"error,omitempty"`
}

var (
	emailRegexp = regexp.MustCompile("(?i)" + // case insensitive
		"^[a-z0-9!#$%&'*+/=?^_`{|}~.-]+" + // local part
		"@" +
		"[a-z0-9-]+(\\.[a-z0-9-]+)+\\.?$") // domain part
)

func IsValidEmail(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegexp.MatchString(email)
}

func init() {
	logFlags := log.LstdFlags | log.Lshortfile
	InfoLogger = log.New(os.Stdout, "[INFO] ", logFlags)
	ErrLogger = log.New(os.Stderr, "[ERROR] ", logFlags)
	WarnLogger = log.New(os.Stdout, "[WARN] ", logFlags)
}
