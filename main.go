package main

import (
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
		// Export the name of the bucket
		ctx.Export("bucketName", bucket.ID())
		ctx.Export("ssmName", ssmParm.Name)
		return nil
	})
}
