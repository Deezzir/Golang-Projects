package handlers

import (
	"net/http"
	"serverless/pkg/config"
	"serverless/pkg/user"
	"serverless/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient          dynamodbiface.DynamoDBAPI
	MethodNotAllowedMsg = utils.ErrorBody{Message: "Method Not Allowed"}
)

func init() {
	dynaClient = config.SetAWSSession()
}

func Handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
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

func getUser(req events.APIGatewayProxyRequest, table string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if len(email) > 0 {
		result, err := user.FetchUser(email, table, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, utils.ErrorBody{Message: *aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}

	result, err := user.FetchUsers(table, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, utils.ErrorBody{Message: *aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func createUser(req events.APIGatewayProxyRequest, table string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.CreateUser(req, table, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, utils.ErrorBody{Message: *aws.String(err.Error())})
	}
	return apiResponse(http.StatusCreated, result)
}

func updateUser(req events.APIGatewayProxyRequest, table string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.UpdateUser(req, table, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, utils.ErrorBody{Message: *aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func deleteUser(req events.APIGatewayProxyRequest, table string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(req, table, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, utils.ErrorBody{Message: *aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, nil)
}

func unhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, MethodNotAllowedMsg)
}
