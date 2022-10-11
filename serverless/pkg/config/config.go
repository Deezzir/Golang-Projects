package config

import (
	"os"
	"serverless/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	AWS_REGION string
	TABLE_NAME string
)

func init() {
	setEnv()
}

func setEnv() {
	AWS_REGION = os.Getenv("AWS_REGION")
	if AWS_REGION == "" {
		utils.ErrLogger.Fatalln("AWS_REGION is not set")
	}

	TABLE_NAME = "go-serverless"
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
