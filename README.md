# SSM enforcement tool in Terraform

This project contains a set of infrastructure implemented in Terraform to monitor your "not-managed-by-SSM" instances accross all regions.

It automatically attaches the proper permissions to your instances, then sends the alert to your email if you need to install SSM Agent or check the network connectivity to SSM endpoint.

## Usage

You neet to have [AWS CLI](https://github.com/aws/aws-cli) and [Terraform](https://github.com/hashicorp/terraform) installed on your machine.

You can simply deploy SSM enforcement tool to your AWS account as follows.

```
./deploy.sh -e your@email.com -r ap-northeast-1
```

You will recieve an email to confirm the subscription to the alert.
