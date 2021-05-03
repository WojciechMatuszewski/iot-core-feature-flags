package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

func main() {
	lambda.Start(handler)
}

type Input struct {
	FlagName string `json:"flagName"`
	Value    bool   `json:"value"`
}

type FlagItem struct {
	PK    string `dynamodbav:"pk"`
	SK    string `dynamodbav:"sk"`
	Value bool   `dynamodbav:"value"`
}

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (
	events.APIGatewayV2HTTPResponse, error) {

	var input Input
	err := json.Unmarshal([]byte(event.Body), &input)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to parse the JSON",
		}, nil
	}

	clientID, found := event.PathParameters["clientID"]
	if !found {
		panic(errors.New("clientID parameter not found"))
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)

	item := FlagItem{
		PK:    clientID,
		SK:    input.FlagName,
		Value: input.Value,
	}
	avs, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}

	_, err = dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      avs,
		TableName: aws.String(os.Getenv("FLAGS_TABLE_NAME")),
	})
	if err != nil {
		panic(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "Flag set",
	}, nil

}
