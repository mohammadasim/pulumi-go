package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// create a context
		conf := config.New(ctx, "")
		bucket, err := s3.NewBucketV2(ctx, conf.Require("bucketName"), nil)
		if err != nil {
			return err
		}
		ssmParm, err := ssm.NewParameter(ctx, "lambdaParam", &ssm.ParameterArgs{
			Name:      pulumi.String(conf.Require("ssmName")),
			Type:      pulumi.String(ssm.ParameterTypeString),
			Value:     pulumi.String(conf.Get("ssmValue")),
			Overwrite: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}
		// Get lambda execution policy
		policy, err := iam.LookupPolicy(ctx, &iam.LookupPolicyArgs{
			Name: pulumi.StringRef("AWSLambdaBasicExecutionRole"),
		})
		if err != nil {
			return err
		}
		// Create assume role
		// firs define the policy as a map
		assumeRolePolicyMap := map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": "lambda.amazonaws.com",
					},
					"Action": "sts:AssumeRole",
				},
			},
		}
		// convert the map to json using json.Marshal
		assumeRolePolicy, err := json.Marshal(assumeRolePolicyMap)
		if err != nil {
			return err
		}

		// Create IAM role for lambda
		lambdaRole, err := iam.NewRole(ctx, "goLambdaRole", &iam.RoleArgs{
			Name:              pulumi.String("goLambdaRole"),
			AssumeRolePolicy:  pulumi.String(string(assumeRolePolicy)),
			Description:       pulumi.StringPtr("IAM role for the Golang Lambda Function"),
			ManagedPolicyArns: pulumi.StringArrayArrayInput(policy.Arn),
		})
		// Export the name of the bucket
		ctx.Export("bucketName", bucket.ID())
		ctx.Export("ssmName", ssmParm.Name)
		return nil
	})
}
