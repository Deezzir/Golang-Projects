package config

import (
	"reps-store/pkg/utils"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"log"
	"os"
	"time"
)

var (
	db     *gorm.DB
	dbname string = "test"
	tcp    string = "tcp(127.0.0.1:3306)"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		utils.ErrLogger.Fatalf("Failed to load the .env file - %s\n", err.Error())
	}
}

func getCredentials() string {
	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")

	if user == "" || pass == "" {
		utils.ErrLogger.Fatalln("DBUSERNAME or DBPASSWORD is not present in the .env file")
	}

	return user + ":" + pass
}

func SetDataBase() {
	if db == nil {
		credentials := getCredentials()
		dsn := credentials + "@" + tcp + "/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

		logger := logger.New(
			log.New(os.Stderr, "[DATABASE] ", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)

		d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger,
		})
		if err != nil {
			utils.ErrLogger.Fatalf("Failed to connect to the Data Base - %s\n", err.Error())
		} else {
			utils.InfoLogger.Println("Successfully connected to the DB")
		}

		db = d
	}
}

func GetDataBase() *gorm.DB {
	return db
}
