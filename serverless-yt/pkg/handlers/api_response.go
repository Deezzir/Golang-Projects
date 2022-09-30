package handlers

import (
	"encoding/json"
	"serverless-yt/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
)

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	strBody, err := json.Marshal(body)
	if err != nil {
		utils.ErrLogger.Printf("Failed to encode the Response Body - %s\n", err.Error())
		return nil, err
	}

	res := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: status,
		Body:       string(strBody),
	}

	return &res, nil
}
