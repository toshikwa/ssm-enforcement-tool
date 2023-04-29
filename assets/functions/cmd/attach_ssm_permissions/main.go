package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

var (
	ssmPolicyName       = os.Getenv("SsmPolicyName")
	ssmPolicyArn        = os.Getenv("SsmPolicyArn")
	ssmAgentProfileName = os.Getenv("SsmAgentProfileName")
	ssmAgentProfileArn  = os.Getenv("SsmAgentProfileArn")
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
	ec2Client *ec2.Client
	iamClient *iam.Client
)

func HandleRequest(ctx context.Context, input Input) (Response, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return Response{}, fmt.Errorf("Failed to load config: %s", err)
	}

	// check if instances are attached with profiles
	iamClient = iam.NewFromConfig(cfg)
	pfArns := map[string]struct{}{}
	for _, ri := range input.UnmanagedInstances {
		ec2Client = ec2.NewFromConfig(cfg, func(o *ec2.Options) { o.Region = ri.Region })
		instances, err := DescribeInstances(ctx, ec2Client, ri.InstanceIds)
		if err != nil {
			return Response{}, fmt.Errorf("failed to describe instances in %s: %s", ri.Region, err)
		}
		for _, i := range instances {
			if i.IamInstanceProfile == nil {
				// attach ssm profile to the instance
				err := AssociateIamInstanceProfile(ctx, ec2Client, ssmAgentProfileArn, ssmAgentProfileName, *i.InstanceId)
				if err != nil {
					return Response{}, fmt.Errorf("failed to associate profile (name: %s, arn: %s) to %s: %s", ssmAgentProfileName, ssmAgentProfileArn, *i.InstanceId, err)
				}
			} else {
				// log profile arn for later use
				pfArns[*i.IamInstanceProfile.Arn] = struct{}{}
			}
		}
	}

	// attach SSM agent role
	for pfArn := range pfArns {
		roleName, err := GetInstanceProfileRoleName(ctx, iamClient, pfArn)
		if err != nil {
			return Response{}, fmt.Errorf("failed to get role name of profile %s: %s", pfArn, err)
		}
		policies, err := ListAttachedRolePolicies(ctx, iamClient, roleName)
		if err != nil {
			return Response{}, fmt.Errorf("failed to list policies (role %s): %s", roleName, err)
		}

		for _, p := range policies {
			if *p.PolicyName == ssmPolicyName {
				goto SKIP
			}
			err := AttachRolePolicies(ctx, iamClient, roleName, ssmPolicyArn)
			if err != nil {
				return Response{}, fmt.Errorf("failed to attach policy %s to role %s: %s", ssmPolicyName, roleName, err)
			}
		}
	SKIP:
	}
	return Response{}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
