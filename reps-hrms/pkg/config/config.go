package config

import (
	"context"
	"os"
	"reps-hrms/pkg/utils"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var (
	DB_NAME   string
	MONGO_URI string
)

func init() {
	loadEnv()
	setEnv()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		utils.ErrLogger.Fatalf("Failed to load the .env file - %s\n", err.Error())
	}
}

func setEnv() {
	DB_NAME = os.Getenv("DB_NAME")
	MONGO_URI = os.Getenv("MONGO_URI")

	if DB_NAME == "" || MONGO_URI == "" {
		utils.ErrLogger.Fatalln("DB_NAME or MONGO_URI are not set in .env file")
	}
}

func InitMongoDB() (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db := client.Database(DB_NAME)
	mongoDB := &MongoDB{
		Client: client,
		DB:     db,
	}
	utils.InfoLogger.Println("MongoDB is initializing")
	return mongoDB, nil
}
