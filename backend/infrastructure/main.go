package main

import (
	"fmt"
	"os"
	"path"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type BackendStackProps struct {
	awscdk.StackProps
}

func NewBackendStack(scope constructs.Construct, id string, props *BackendStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	NewAuth(stack, jsii.String("cognito"))

	dataTable := NewDB(stack, jsii.String("dynamo"))

	NewAPI(stack, jsii.String("api"), dataTable)
	NewIOT(stack, jsii.String("iot"), dataTable)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewBackendStack(app, "BackendStack", &BackendStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}

func functionDir(functionName string) string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(functionName)

	return path.Join(pwd, "..", "functions", functionName)
}
