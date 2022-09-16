package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/joho/godotenv"

	"log"
	"os"
)

var (
	db    *gorm.DB
	table string = "repositories"
)

func getCredentials() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("[ERROR]: failed to load the .env file\n")
	}

	user := os.Getenv("DBUSERNAME")
	pass := os.Getenv("DBPASSWORD")

	if user == "" || pass == "" {
		log.Fatalf("[ERROR]: DBUSERNAME or DBPASSWORD is not present in the .env file")
	}

	return user + ":" + pass
}

func SetDB() {
	if db == nil {
		creds := getCredentials()

		d, err := gorm.Open("mysql", creds+"/"+table+"?charset=utf8&parseTime=true&loc=local")
		if err != nil {
			log.Fatalf("[ERROR]: failed to connect to the Data Base\n%s\n", err.Error())
		}

		db = d
	}
}

func GetDB() *gorm.DB {
	return db
}
