package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	snsTopicArn = os.Getenv("SnsTopicArn")
)

type RegionalInstances struct {
	Region      string
	InstanceIds []string
}

type Input struct {
	UnmanagedInstances []RegionalInstances
}

type Response struct{}

var (
	snsClient *sns.Client
)

func HandleRequest(ctx context.Context, input Input) (Response, error) {
	// compose message
	msg := "Unmanaged instances are found. You may need to install SSM Agent or check the network connectivity to SSM endpoint.\n"
	for _, ri := range input.UnmanagedInstances {
		msg += fmt.Sprintf("\n---------- %s ----------\n", ri.Region)
		for _, id := range ri.InstanceIds {
			msg += fmt.Sprintf("%s\n", id)
		}
	}
	// load config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return Response{}, fmt.Errorf("failed to load config: %s", err)
	}
	// publish message
	snsClient = sns.NewFromConfig(cfg)
	err = Publish(ctx, snsClient, msg, snsTopicArn)
	if err != nil {
		return Response{}, fmt.Errorf("failed to publish message: %s", err)
	}
	return Response{}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
