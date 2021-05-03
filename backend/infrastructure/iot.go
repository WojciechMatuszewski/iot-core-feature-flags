package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/pkg/errors"
)

func NewIOT(stack constructs.Construct, id *string, dataTable awsdynamodb.Table) {
	scope := awscdk.NewConstruct(stack, id)

	iotEndpoint, err := getIOTEndpoint()
	if err != nil {
		panic(err)
	}

	notifyClientLambda := awslambda.NewFunction(stack, jsii.String("notifyClientLambda"), &awslambda.FunctionProps{
		Timeout: awscdk.Duration_Seconds(jsii.Number(20)),
		Tracing: awslambda.Tracing_ACTIVE,
		Code: awslambda.AssetCode_FromAsset(jsii.String(functionDir("notify-client")), &awss3assets.AssetOptions{
			Bundling: &awscdk.BundlingOptions{
				Image: awslambda.Runtime_GO_1_X().BundlingDockerImage(),
				Command: &[]*string{
					jsii.String("bash"),
					jsii.String("-c"),
					jsii.String("go build -o /asset-output/main"),
				},
				User: jsii.String("root"),
			},
		}),
		Runtime: awslambda.Runtime_GO_1_X(),
		Handler: jsii.String("main"),
	})
	notifyClientLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("iot:Publish"),
		},
		Effect: awsiam.Effect_ALLOW,
		Resources: &[]*string{
			jsii.String(fmt.Sprintf("arn:%v:iot:%v:%v:topic/flags/*", *awscdk.Aws_PARTITION(), *awscdk.Aws_REGION(), *awscdk.Aws_ACCOUNT_ID())),
		},
	}))

	notifyClientLambda.AddEventSource(awslambdaeventsources.NewDynamoEventSource(dataTable, &awslambdaeventsources.DynamoEventSourceProps{
		StartingPosition:   awslambda.StartingPosition_LATEST,
		BatchSize:          jsii.Number(1),
		BisectBatchOnError: jsii.Bool(false),
		RetryAttempts:      jsii.Number(0),
	}))

	// awsiot.NewCfnTopicRule(scope, jsii.String("lambdarule"), &awsiot.CfnTopicRuleProps{
	// 	TopicRulePayload: awsiot.CfnTopicRule_TopicRulePayloadProperty{
	// 		Actions: []awsiot.CfnTopicRule_ActionProperty{
	// 			{
	// 				Lambda: awsiot.CfnTopicRule_LambdaActionProperty{
	// 					FunctionArn: getFlagsLambda.FunctionArn(),
	// 				},
	// 			},
	// 		},
	// 		Sql:          jsii.String("SELECT * FROM 'myTopic/+' WHERE action = 'getFlags'"),
	// 		RuleDisabled: jsii.Bool(true),
	// 	},
	// 	RuleName: jsii.String("lambdarule"),
	// })

	// awslambda.NewCfnPermission(scope, jsii.String("IotInvokeLambda"), &awslambda.CfnPermissionProps{
	// 	Action:       jsii.String("lambda:InvokeFunction"),
	// 	FunctionName: getFlagsLambda.FunctionName(),
	// 	Principal:    jsii.String("iot.amazonaws.com"),
	// })

	awscdk.NewCfnOutput(scope, jsii.String("iotEndpoint"), &awscdk.CfnOutputProps{
		Value: jsii.String(iotEndpoint),
	})
}

type CmdOutput struct {
	EndpointAddress string `json:"endpointAddress"`
}

// An alternative would be to write custom resource.
func getIOTEndpoint() (string, error) {
	var out CmdOutput
	cmd := exec.Command("aws", "iot", "describe-endpoint", "--endpoint-type", "iot:Data-ATS")

	buff, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to execute the command")
	}

	err = json.Unmarshal(buff, &out)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal the data")
	}

	return out.EndpointAddress, nil
}
