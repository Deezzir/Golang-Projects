package main

import (
	"serverless-yt/pkg/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.Handler)
}
