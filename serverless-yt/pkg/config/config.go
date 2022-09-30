package config

import (
	"os"
	"serverless-yt/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/joho/godotenv"
)

var (
	AWS_REGION string
	TABLE_NAME string
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
	AWS_REGION = os.Getenv("AWS_REGION")
	if AWS_REGION == "" {
		utils.ErrLogger.Fatalln("AWS_REGION is not set in .env file")
	}

	TABLE_NAME = os.Getenv("TABLE_NAME")
	if TABLE_NAME == "" {
		utils.ErrLogger.Fatalln("TABLE_NAME is not set in .env file")
	}
}

func SetAWSSession() dynamodbiface.DynamoDBAPI {
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION),
	})
	if err != nil {
		utils.ErrLogger.Fatalf("Failed to create a new AWS session - %s\n", err.Error())
	}

	dynaClient := dynamodb.New(awsSession)
	return dynaClient
}
