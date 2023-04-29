package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SnsPublishApi interface {
	Publish(
		ctx context.Context,
		params *sns.PublishInput,
		optFns ...func(*sns.Options),
	) (*sns.PublishOutput, error)
}

func Publish(ctx context.Context, api SnsPublishApi, msg, topicArn string) error {
	_, err := api.Publish(
		ctx,
		&sns.PublishInput{
			Message:  aws.String(msg),
			TopicArn: aws.String(topicArn),
		},
	)
	return err
}
