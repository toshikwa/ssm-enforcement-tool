package main

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type Ec2DescribeInstancesApi interface {
	DescribeInstances(
		ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options),
	) (*ec2.DescribeInstancesOutput, error)
}

func DescribeInstances(ctx context.Context, api Ec2DescribeInstancesApi, ids []string) ([]ec2Types.Instance, error) {
	var token *string
	instances := []ec2Types.Instance{}
	for {
		// describe instances
		res, err := api.DescribeInstances(ctx, &ec2.DescribeInstancesInput{InstanceIds: ids, NextToken: token})
		if err != nil {
			return nil, err
		}
		// append all instances
		for _, r := range res.Reservations {
			instances = append(instances, r.Instances...)
		}
		token = res.NextToken
		if token == nil {
			break
		}
	}
	return instances, nil
}

type Ec2AssociateIamInstanceProfileApi interface {
	AssociateIamInstanceProfile(
		ctx context.Context,
		params *ec2.AssociateIamInstanceProfileInput,
		optFns ...func(*ec2.Options),
	) (*ec2.AssociateIamInstanceProfileOutput, error)
}

func AssociateIamInstanceProfile(
	ctx context.Context,
	api Ec2AssociateIamInstanceProfileApi,
	profileArn, profileName, instanceId string,
) error {
	_, err := api.AssociateIamInstanceProfile(
		ctx, &ec2.AssociateIamInstanceProfileInput{
			IamInstanceProfile: &ec2Types.IamInstanceProfileSpecification{
				Arn:  aws.String(profileArn),
				Name: aws.String(profileName),
			},
			InstanceId: &instanceId,
		},
	)
	return err
}

type IamGetInstanceProfileApi interface {
	GetInstanceProfile(
		ctx context.Context,
		params *iam.GetInstanceProfileInput,
		optFns ...func(*iam.Options),
	) (*iam.GetInstanceProfileOutput, error)
}

func GetInstanceProfileRoleName(ctx context.Context, api IamGetInstanceProfileApi, pfArn string) (string, error) {
	pfName := pfArn[strings.LastIndex(pfArn, "/")+1:]
	res, err := api.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{InstanceProfileName: aws.String(pfName)})
	if err != nil {
		return "", err
	}
	return *res.InstanceProfile.Roles[0].RoleName, nil
}

type IamListAttachedRolePoliciesApi interface {
	ListAttachedRolePolicies(
		ctx context.Context,
		params *iam.ListAttachedRolePoliciesInput,
		optFns ...func(*iam.Options),
	) (*iam.ListAttachedRolePoliciesOutput, error)
}

func ListAttachedRolePolicies(ctx context.Context, api IamListAttachedRolePoliciesApi, roleName string) ([]iamTypes.AttachedPolicy, error) {
	res, err := api.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{RoleName: aws.String(roleName)})
	if err != nil {
		return nil, err
	}
	return res.AttachedPolicies, nil
}

type IamAttachRolePolicyApi interface {
	AttachRolePolicy(
		ctx context.Context,
		params *iam.AttachRolePolicyInput,
		optFns ...func(*iam.Options),
	) (*iam.AttachRolePolicyOutput, error)
}

func AttachRolePolicies(ctx context.Context, api IamAttachRolePolicyApi, roleName, policyArn string) error {
	_, err := api.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{RoleName: aws.String(roleName), PolicyArn: aws.String(policyArn)})
	return err
}
