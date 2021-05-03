package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

func NewAuth(stack constructs.Construct, id *string) {
	scope := awscdk.NewConstruct(stack, id)

	userPool := awscognito.NewUserPool(scope, jsii.String("userPool"), &awscognito.UserPoolProps{
		LambdaTriggers: &awscognito.UserPoolTriggers{},
		PasswordPolicy: &awscognito.PasswordPolicy{
			MinLength:        jsii.Number(6),
			RequireDigits:    jsii.Bool(false),
			RequireLowercase: jsii.Bool(false),
			RequireSymbols:   jsii.Bool(false),
			RequireUppercase: jsii.Bool(false),
		},
		SelfSignUpEnabled: jsii.Bool(true),
		SignInAliases: &awscognito.SignInAliases{
			Email: jsii.Bool(true),
		},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	userPoolClient := awscognito.NewUserPoolClient(scope, jsii.String("userPoolClient"), &awscognito.UserPoolClientProps{
		AuthFlows:      &awscognito.AuthFlow{UserPassword: jsii.Bool(true), UserSrp: jsii.Bool(true)},
		GenerateSecret: jsii.Bool(false),
		SupportedIdentityProviders: &[]awscognito.UserPoolClientIdentityProvider{
			awscognito.UserPoolClientIdentityProvider_COGNITO(),
		},

		UserPool: userPool,
	})

	identityPool := awscognito.NewCfnIdentityPool(scope, jsii.String("identityPool"), &awscognito.CfnIdentityPoolProps{
		AllowUnauthenticatedIdentities: jsii.Bool(true),
		CognitoIdentityProviders: []awscognito.CfnIdentityPool_CognitoIdentityProviderProperty{
			{
				ClientId:     userPoolClient.UserPoolClientId(),
				ProviderName: userPool.UserPoolProviderName(),
			},
		},
	})

	identityRole := awsiam.NewRole(scope, jsii.String("identityRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewFederatedPrincipal(jsii.String("cognito-identity.amazonaws.com"), &map[string]interface{}{}, jsii.String("sts:AssumeRoleWithWebIdentity")),
		InlinePolicies: &map[string]awsiam.PolicyDocument{
			"AllowAccess": awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
				Statements: &[]awsiam.PolicyStatement{
					awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
						Actions: &[]*string{
							// jsii.String("iot:Receive"),
							// jsii.String("iot:Publish"),
							// jsii.String("iot:Subscribe"),
							jsii.String("iot:Connect"),
						},
						Effect: awsiam.Effect_ALLOW,
						Resources: &[]*string{
							// jsii.String("*"),
							jsii.String(fmt.Sprintf("arn:%v:iot:%v:%v:client/*", *awscdk.Aws_PARTITION(), *awscdk.Aws_REGION(), *awscdk.Aws_ACCOUNT_ID())),
						},
					}),
					awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
						Actions: &[]*string{
							jsii.String("iot:Subscribe"),
						},
						Effect: awsiam.Effect_ALLOW,
						Resources: &[]*string{
							jsii.String(fmt.Sprintf("arn:%v:iot:%v:%v:topicfilter/flags/*", *awscdk.Aws_PARTITION(), *awscdk.Aws_REGION(), *awscdk.Aws_ACCOUNT_ID())),
						},
					}),
					awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
						Actions: &[]*string{
							jsii.String("iot:Receive"),
						},
						Effect: awsiam.Effect_ALLOW,
						Resources: &[]*string{
							jsii.String(fmt.Sprintf("arn:%v:iot:%v:%v:topic/flags/*", *awscdk.Aws_PARTITION(), *awscdk.Aws_REGION(), *awscdk.Aws_ACCOUNT_ID())),
						},
					}),
				},
			}),
		},
	})

	adminGroup := awscognito.NewCfnUserPoolGroup(scope, jsii.String("adminGroup"), &awscognito.CfnUserPoolGroupProps{
		UserPoolId:  userPool.UserPoolId(),
		Description: jsii.String("Admin group"),
		GroupName:   jsii.String("admin"),
	})

	createAdminLambda := awslambda.NewFunction(stack, jsii.String("createAdmin"), &awslambda.FunctionProps{
		Timeout: awscdk.Duration_Seconds(jsii.Number(20)),
		Tracing: awslambda.Tracing_ACTIVE,
		Code: awslambda.AssetCode_FromAsset(jsii.String(functionDir("create-admin")), &awss3assets.AssetOptions{
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

	createAdminLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("cognito-idp:AdminConfirmSignUp"),
			jsii.String("cognito-idp:AdminAddUserToGroup"),
		},
		Effect: awsiam.Effect_ALLOW,
		Resources: &[]*string{
			jsii.String(fmt.Sprintf("arn:%v:cognito-idp:%v:%v:userpool/%v", *awscdk.Aws_PARTITION(), *awscdk.Aws_REGION(), *awscdk.Aws_ACCOUNT_ID(), *userPool.UserPoolId())),
		},
	}))

	admin := awscdk.NewCustomResource(scope, jsii.String("createAdminResource"), &awscdk.CustomResourceProps{
		ServiceToken: createAdminLambda.FunctionArn(),
		Properties: &map[string]interface{}{
			"userPoolId":       *userPool.UserPoolId(),
			"userPoolClientId": *userPoolClient.UserPoolClientId(),
			"adminGroupName":   *adminGroup.GroupName(),
		},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	awscognito.NewCfnIdentityPoolRoleAttachment(scope, jsii.String("identityPoolroles"), &awscognito.CfnIdentityPoolRoleAttachmentProps{
		IdentityPoolId: awscdk.Fn_Ref(identityPool.LogicalId()),
		Roles: map[string]string{
			"authenticated":   *identityRole.RoleArn(),
			"unauthenticated": *identityRole.RoleArn(),
		},
	})

	awscdk.NewCfnOutput(scope, jsii.String("userPoolId"), &awscdk.CfnOutputProps{
		Value: userPool.UserPoolId(),
	})

	awscdk.NewCfnOutput(scope, jsii.String("userPoolClientId"), &awscdk.CfnOutputProps{
		Value: userPoolClient.UserPoolClientId(),
	})

	awscdk.NewCfnOutput(scope, jsii.String("identityPoolId"), &awscdk.CfnOutputProps{
		Value: awscdk.Fn_Ref(identityPool.LogicalId()),
	})

	awscdk.NewCfnOutput(scope, jsii.String("adminUsername"), &awscdk.CfnOutputProps{
		Value: admin.GetAttString(jsii.String("adminUsername")),
	})

	awscdk.NewCfnOutput(scope, jsii.String("adminPassword"), &awscdk.CfnOutputProps{
		Value: admin.GetAttString(jsii.String("adminPassword")),
	})
}
