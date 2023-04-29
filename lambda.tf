module "list_unmanaged_instances_lambda" {
  source = "./modules/lambda"
  name   = "list_unmanaged_instances"
  app    = local.app
  policy_actions = [
    "ec2:DescribeRegions",
    "ec2:DescribeInstances",
    "ssm:DescribeInstanceInformation",
  ]
  envs = {}
}

module "attach_ssm_permission_lambda" {
  source = "./modules/lambda"
  name   = "attach_ssm_permissions"
  app    = local.app
  policy_actions = [
    "ec2:DescribeInstances",
    "ec2:AssociateIamInstanceProfile",
    "iam:GetInstanceProfile",
    "iam:ListAttachedRolePolicies",
    "iam:AttachRolePolicy",
    "iam:PassRole",
  ]
  envs = {
    SsmPolicyArn        = data.aws_iam_policy.ssm_policy.arn
    SsmPolicyName       = data.aws_iam_policy.ssm_policy.name
    SsmAgentProfileName = aws_iam_instance_profile.ssm_agent_profile.name
    SsmAgentProfileArn  = aws_iam_instance_profile.ssm_agent_profile.arn
  }
}

module "send_alert_to_sns_lambda" {
  source = "./modules/lambda"
  name   = "send_alert_to_sns"
  app    = local.app
  policy_actions = [
    "sns:Publish",
  ]
  envs = {
    SnsTopicArn = aws_sns_topic.default.arn
  }
}
