package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

func NewDB(stack constructs.Construct, id *string) awsdynamodb.Table {
	scope := awscdk.NewConstruct(stack, id)

	flagsTable := awsdynamodb.NewTable(scope, jsii.String("flags"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("pk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("sk"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		Stream:        awsdynamodb.StreamViewType_NEW_IMAGE,
	})

	awscdk.NewCfnOutput(scope, jsii.String("flagsTableName"), &awscdk.CfnOutputProps{
		Value: flagsTable.TableName(),
	})

	return flagsTable
}
