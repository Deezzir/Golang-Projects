package handlers

import (
	"serverless-yt/pkg/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func init() {
	dynaClient = config.SetAWSSession()
}

func Handler(req events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	switch req.HTTPMethod {
	case "GET":
		return getUser(req, config.TABLE_NAME, dynaClient)
	case "POST":
		return createUser(req, config.TABLE_NAME, dynaClient)
	case "PUT":
		return updateUser(req, config.TABLE_NAME, dynaClient)
	case "DELETE":
		return deleteUser(req, config.TABLE_NAME, dynaClient)
	default:
		return unhandledMethod()
	}
}

func getUser() {

}

func createUser() {

}

func updateUser() {

}

func deleteUser() {

}

func unhandledMethod() {

}
