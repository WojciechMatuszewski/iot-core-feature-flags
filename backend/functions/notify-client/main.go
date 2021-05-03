package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, sEvent events.DynamoDBEvent) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	event := sEvent.Records[0]

	clientID := event.Change.NewImage["pk"].String()

	flagName := event.Change.NewImage["sk"].String()
	flagValue := event.Change.NewImage["value"].Boolean()

	flagUpdate := make(map[string]bool)
	flagUpdate[flagName] = flagValue

	buf, err := json.Marshal(flagUpdate)
	if err != nil {
		panic(err)
	}

	iotClient := iotdataplane.NewFromConfig(cfg)
	_, err = iotClient.Publish(ctx, &iotdataplane.PublishInput{
		Topic:   aws.String(fmt.Sprintf("%v/%v", "flags", clientID)),
		Payload: buf,
		Qos:     1,
	})
	if err != nil {
		panic(err)
	}

	return nil
}
