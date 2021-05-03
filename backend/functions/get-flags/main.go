package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	lambda.Start(handler)
}

type FlagItem struct {
	PK    string `dynamodbav:"pk"`
	SK    string `dynamodbav:"sk"`
	Value bool   `dynamodbav:"value"`
}

type Output map[string]string

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	clientID, found := event.PathParameters["clientID"]
	if !found {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Missing parameter",
		}, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	db := dynamodb.NewFromConfig(cfg)

	exp, err := expression.NewBuilder().WithKeyCondition(expression.KeyEqual(expression.Key("pk"), expression.Value(clientID))).Build()
	if err != nil {
		panic(err)
	}

	out, err := db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(os.Getenv("FLAGS_TABLE_NAME")),
		KeyConditionExpression:    exp.KeyCondition(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
	})
	if err != nil {
		panic(err)
	}

	flagsList := make([]FlagItem, len(out.Items))
	err = attributevalue.UnmarshalListOfMaps(out.Items, &flagsList)
	if err != nil {
		panic(err)
	}

	output := make(map[string]bool)
	for _, v := range flagsList {
		output[v.SK] = v.Value
	}

	outBuf, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(outBuf),
	}, nil

}
