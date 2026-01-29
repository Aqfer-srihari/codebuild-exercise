package main

import (
	"log"
	"question5updation/internal/handler"
	"question5updation/internal/storage"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func lambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return handler.LambdaRouter(req), nil
}

func main() {

	if err := storage.InitDB(); err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	lambda.Start(lambdaHandler)
}
