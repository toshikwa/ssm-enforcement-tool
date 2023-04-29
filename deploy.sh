#/bin/bash

USAGE_MSG="Usage: \`./deploy.sh -e your@email.com -r ap-northeast-1\`"
while getopts e:r: OPT
do
  case $OPT in
    "e" ) EMAIL="$OPTARG" ;;
    "r" ) AWS_REGION="$OPTARG" ;;
    * ) echo ${USAGE_MSG} ; exit 1 ;;
  esac
done

if [ -z "${EMAIL}" ]; then
  echo "You need to specify your own email address."
  echo ${USAGE_MSG}
  exit 1
fi

if [ -z "${AWS_REGION}" ]; then
  echo "You need to specify the region."
  echo ${USAGE_MSG}
  exit 1
fi

# account id
ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)

# make S3 bucket to store tfstate
aws s3 mb s3://ssm-enforcement-tool-${ACCOUNT_ID} --region ${AWS_REGION}

# initialize
terraform init \
  -backend-config="bucket=ssm-enforcement-tool-${ACCOUNT_ID}" \
  -backend-config="region=${AWS_REGION}" \
  -reconfigure

# deploy
terraform apply -auto-approve -var="email=${EMAIL}"