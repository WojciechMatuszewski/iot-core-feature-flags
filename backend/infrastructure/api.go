package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

func NewAPI(stack constructs.Construct, id *string, flagsTable awsdynamodb.Table) {
	scope := awscdk.NewConstruct(stack, id)

	flagsAPI := awsapigatewayv2.NewHttpApi(scope, jsii.String("flagsAPI"), &awsapigatewayv2.HttpApiProps{
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowCredentials: jsii.Bool(false),
			AllowHeaders: &[]*string{
				jsii.String("*"),
			},
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{awsapigatewayv2.CorsHttpMethod_ANY},
			AllowOrigins: &[]*string{jsii.String("*")},
		},
	})

	getFlagsLambda := awslambda.NewFunction(stack, jsii.String("getFlags"), &awslambda.FunctionProps{
		Timeout: awscdk.Duration_Seconds(jsii.Number(20)),
		Tracing: awslambda.Tracing_ACTIVE,
		Code: awslambda.AssetCode_FromAsset(jsii.String(functionDir("get-flags")), &awss3assets.AssetOptions{
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
		Environment: &map[string]*string{
			"FLAGS_TABLE_NAME": flagsTable.TableName(),
		},
	})
	flagsTable.GrantReadData(getFlagsLambda)

	setFlagLambda := awslambda.NewFunction(stack, jsii.String("setFlag"), &awslambda.FunctionProps{
		Timeout: awscdk.Duration_Seconds(jsii.Number(20)),
		Tracing: awslambda.Tracing_ACTIVE,
		Code: awslambda.AssetCode_FromAsset(jsii.String(functionDir("set-flag")), &awss3assets.AssetOptions{
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
		Environment: &map[string]*string{
			"FLAGS_TABLE_NAME": flagsTable.TableName(),
		},
	})
	flagsTable.GrantWriteData(setFlagLambda)

	getFlagsIntegration := awsapigatewayv2integrations.NewLambdaProxyIntegration(&awsapigatewayv2integrations.LambdaProxyIntegrationProps{
		Handler:              getFlagsLambda,
		PayloadFormatVersion: awsapigatewayv2.PayloadFormatVersion_VERSION_2_0(),
	})
	flagsAPI.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: getFlagsIntegration,
		Path:        jsii.String("/flags/{clientID}"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
		},
	})

	setFlagIntegration := awsapigatewayv2integrations.NewLambdaProxyIntegration(&awsapigatewayv2integrations.LambdaProxyIntegrationProps{
		Handler:              setFlagLambda,
		PayloadFormatVersion: awsapigatewayv2.PayloadFormatVersion_VERSION_2_0(),
	})
	flagsAPI.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Integration: setFlagIntegration,
		Path:        jsii.String("/flag/{clientID}"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
		},
	})

	awscdk.NewCfnOutput(scope, jsii.String("apiRootUrl"), &awscdk.CfnOutputProps{
		Value: flagsAPI.ApiEndpoint(),
	})

}
