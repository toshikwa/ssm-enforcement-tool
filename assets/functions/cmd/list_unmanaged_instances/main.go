package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type RegionalInstances struct {
	Region      string
	InstanceIds []string
}

type Response struct {
	UnmanagedInstances []RegionalInstances
	Compliant          bool
}

func Diff(sub, all []string) ([]string, error) {
	if len(sub) == len(all) {
		return []string{}, nil
	}
	// sort slices to compare
	sort.Strings(sub)
	sort.Strings(all)
	// list diff
	res := []string{}
	i, j := 0, 0
	for j != len(sub) && i < len(all) {
		if all[i] == sub[j] {
			j += 1
		} else {
			res = append(res, all[i])
		}
		i += 1
	}
	return res, nil
}

var (
	ec2Client *ec2.Client
	ssmClient *ssm.Client
)

func HandleRequest(ctx context.Context) (Response, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return Response{}, fmt.Errorf("Failed to load config: %s", err)
	}

	// list all regions
	ec2Client = ec2.NewFromConfig(cfg)
	regions, err := ListAllRegions(ctx, ec2Client)
	if err != nil {
		return Response{}, fmt.Errorf("failed to list all regions: %s", err)
	}

	instances := []RegionalInstances{}
	for _, r := range regions {
		// api clients for each region
		ec2Client = ec2.NewFromConfig(cfg, func(o *ec2.Options) { o.Region = r })
		ssmClient = ssm.NewFromConfig(cfg, func(o *ssm.Options) { o.Region = r })
		// list instance ids
		ids, err := ListRunningInstanceIds(ctx, ec2Client)
		if err != nil {
			return Response{}, fmt.Errorf("failed to list running instance ids in %s: %s", r, err)
		}
		// filter managed instance ids
		managedIds, err := FilterManagedInstanceIds(ctx, ssmClient, ids)
		if err != nil {
			return Response{}, fmt.Errorf("failed to filter managed instance ids in %s: %s", r, err)
		}
		// list unmanaged instance ids
		unmanagedIds, err := Diff(managedIds, ids)
		if err != nil {
			return Response{}, fmt.Errorf("failed to diff unmanaged instance ids: %s", err)
		}
		if len(unmanagedIds) != 0 {
			instances = append(instances, RegionalInstances{Region: r, InstanceIds: unmanagedIds})
		}
	}
	return Response{UnmanagedInstances: instances, Compliant: len(instances) == 0}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
