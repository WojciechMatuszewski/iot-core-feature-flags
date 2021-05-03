package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func main() {
	lambda.Start(cfn.LambdaWrap(handler))
}

func handler(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {
	if event.RequestType != cfn.RequestCreate {
		return
	}

	userPoolID, ok := event.ResourceProperties["userPoolId"].(string)
	if !ok {
		panic(errors.New("userPoolId not found"))
	}

	userPoolClientID, ok := event.ResourceProperties["userPoolClientId"].(string)
	if !ok {
		panic(errors.New("userPoolClientId ot found"))
	}

	adminGroupName, ok := event.ResourceProperties["adminGroupName"].(string)
	if !ok {
		panic(errors.New("adminGroupName not found"))
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	cognitoClient := cognitoidentityprovider.NewFromConfig(cfg)

	adminUsername := "admin@admin.com"
	adminPassword := "test12345"

	_, err = cognitoClient.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(userPoolClientID),
		Password: aws.String(adminPassword),
		Username: aws.String(adminUsername),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(adminUsername),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = cognitoClient.AdminConfirmSignUp(ctx, &cognitoidentityprovider.AdminConfirmSignUpInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(adminUsername),
	})
	if err != nil {
		panic(err)
	}

	_, err = cognitoClient.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String(adminGroupName),
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(adminUsername),
	})
	if err != nil {
		panic(err)
	}

	data = map[string]interface{}{
		"adminUsername": adminUsername,
		"adminPassword": adminPassword,
	}

	return
}
