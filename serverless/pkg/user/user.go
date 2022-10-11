package user

import (
	"encoding/json"
	"errors"
	"serverless/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	ErrorFailedToFetch   = "failed to fetch record"
	ErrorFailedToDecode  = "failed to decode record"
	ErrorFailedToEncode  = "failed to encode record"
	ErrorFailedToDelete  = "failed to delete record"
	ErrorFailedToSave    = "failed to save record"
	ErrorInvalidUserData = "invalid User Data"
	ErrorInvalidEmail    = "invalid Email"
	ErrorUserExists      = "user already exists"
	ErrorUserNotExists   = "user does not exist"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email, table string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(table),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetch)
	}

	user := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, errors.New(ErrorFailedToDecode)
	}

	return user, nil
}

func FetchUsers(table string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(table),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetch)
	}

	users := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, users)
	if err != nil {
		return nil, errors.New(ErrorFailedToDecode)
	}

	return users, nil
}

func CreateUser(req events.APIGatewayProxyRequest, table string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var user User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !utils.IsValidEmail(user.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	u, _ := FetchUser(user.Email, table, dynaClient)
	if u != nil && len(u.Email) != 0 {
		return nil, errors.New(ErrorUserExists)
	}

	obj, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorFailedToEncode)
	}

	input := &dynamodb.PutItemInput{
		Item:      obj,
		TableName: aws.String(table),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToSave)
	}

	return &user, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, table string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var user User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	u, _ := FetchUser(user.Email, table, dynaClient)
	if u == nil || len(u.Email) == 0 {
		return nil, errors.New(ErrorUserNotExists)
	}

	obj, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorFailedToEncode)
	}

	input := &dynamodb.PutItemInput{
		Item:      obj,
		TableName: aws.String(table),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToSave)
	}

	return &user, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, table string, dynaClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(table),
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorFailedToDelete)
	}

	return nil
}
