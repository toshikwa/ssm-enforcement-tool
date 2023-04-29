# SSM enforcement tool in Terraform

This project contains a set of infrastructure implemented in Terraform to monitor your "not-managed-by-SSM" instances accross all regions.

It automatically attaches the proper permissions to your instances, then sends the alert to your email if you need to install SSM Agent or check the network connectivity to SSM endpoint.

## Usage

You neet to have [AWS CLI](https://github.com/aws/aws-cli) and [Terraform](https://github.com/hashicorp/terraform) installed on your machine.

Then, execute the command below to deploy SSM enforcement tool to your AWS account.

```
# you need to set your own email address to receive the alert
export EMAIL="YOUR_EMAIL_ADDRESS"
export AWS_DEFAULT_REGION=ap-northeast-1
export ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)

# make S3 bucket to store tfstate
aws s3 mb s3://ssm-enforcement-tool-${ACCOUNT_ID} --region ${AWS_DEFAULT_REGION}

# initialize (you may need to add `-reconfigure` option)
terraform init \
  -backend-config="bucket=ssm-enforcement-tool-${ACCOUNT_ID}"\
  -backend-config="region=${AWS_DEFAULT_REGION}"

# deploy
terraform apply -auto-approve -var="email=${EMAIL}"
```

You will recieve an email to confirm the subscription to the alert.
