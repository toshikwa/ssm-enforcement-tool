package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type Ec2DescribeRegionsApi interface {
	DescribeRegions(
		ctx context.Context,
		params *ec2.DescribeRegionsInput,
		optFns ...func(*ec2.Options),
	) (*ec2.DescribeRegionsOutput, error)
}

func ListAllRegions(ctx context.Context, api Ec2DescribeRegionsApi) ([]string, error) {
	res, err := api.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}
	regions := []string{}
	for _, r := range res.Regions {
		regions = append(regions, *r.RegionName)
	}
	return regions, nil
}

type Ec2DescribeInstancesApi interface {
	DescribeInstances(
		ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options),
	) (*ec2.DescribeInstancesOutput, error)
}

func ListRunningInstanceIds(ctx context.Context, api Ec2DescribeInstancesApi) ([]string, error) {
	var token *string
	filters := []ec2Types.Filter{
		{
			Name:   aws.String("instance-state-name"),
			Values: []string{"running"},
		},
	}
	instanceIds := []string{}
	for {
		// describe running instances in the region
		res, err := api.DescribeInstances(ctx, &ec2.DescribeInstancesInput{Filters: filters, NextToken: token})
		if err != nil {
			return nil, err
		}
		// append all instances
		for _, r := range res.Reservations {
			for _, i := range r.Instances {
				instanceIds = append(instanceIds, *i.InstanceId)
			}
		}
		token = res.NextToken
		if token == nil {
			break
		}
	}
	return instanceIds, nil
}

type SsmDescribeInstanceInformationApi interface {
	DescribeInstanceInformation(
		ctx context.Context,
		params *ssm.DescribeInstanceInformationInput,
		optFns ...func(*ssm.Options),
	) (*ssm.DescribeInstanceInformationOutput, error)
}

func FilterManagedInstanceIds(ctx context.Context, api SsmDescribeInstanceInformationApi, instanceIds []string) ([]string, error) {
	if len(instanceIds) == 0 {
		return []string{}, nil
	}
	var token *string
	filters := []ssmTypes.InstanceInformationStringFilter{
		{
			Key:    aws.String("InstanceIds"),
			Values: instanceIds,
		},
	}
	managedInstanceIds := []string{}
	for {
		res, err := api.DescribeInstanceInformation(ctx, &ssm.DescribeInstanceInformationInput{Filters: filters, NextToken: token})
		if err != nil {
			return nil, err
		}
		for _, info := range res.InstanceInformationList {
			managedInstanceIds = append(managedInstanceIds, *info.InstanceId)
		}
		token = res.NextToken
		if token == nil {
			break
		}
	}
	return managedInstanceIds, nil
}
