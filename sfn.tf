# >>>>>>>>>> IAM >>>>>>>>>> #
data "aws_iam_policy_document" "assume_role_sfn" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["states.${data.aws_region.current.name}.amazonaws.com"]
    }
  }
}
resource "aws_iam_role" "sfn_role" {
  name               = "${local.app}_sfn_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_sfn.json
}
resource "aws_iam_policy_attachment" "sfn" {
  name       = "${local.app}_policy_attachment"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaRole"
  roles      = ["${aws_iam_role.sfn_role.name}"]
}
# <<<<<<<<<< IAM <<<<<<<<<< #

resource "aws_sfn_state_machine" "default" {
  name     = local.app
  type     = "STANDARD"
  role_arn = aws_iam_role.sfn_role.arn
  definition = jsonencode(
    {
      Comment = "State machine which monitores unmanaged instances and triggeres the alert"
      StartAt = "ListUnmanagedInstances"
      States = {
        ListUnmanagedInstances = {
          Type     = "Task"
          Resource = module.list_unmanaged_instances_lambda.function_arn
          Next     = "CheckCompliance"
        }
        CheckCompliance = {
          Type = "Choice"
          Choices = [
            {
              Comment       = "All instances are compliant"
              Variable      = "$.Compliant"
              BooleanEquals = true
              Next          = "Succeed"
            },
            {
              Comment       = "There are non-compliant instances"
              Variable      = "$.Compliant"
              BooleanEquals = false
              Next          = "AttachSsmPermissions"
            },
          ]
        }
        AttachSsmPermissions = {
          Type      = "Task"
          Resource  = module.attach_ssm_permission_lambda.function_arn
          InputPath = "$"
          Next      = "WaitFor10Mins"
        }
        WaitFor10Mins = {
          Type    = "Wait"
          Seconds = 600
          Next    = "ListUnmanagedInstancesAgain"
        }
        ListUnmanagedInstancesAgain = {
          Type     = "Task"
          Resource = module.list_unmanaged_instances_lambda.function_arn
          Next     = "CheckComplianceAgain"
        }
        CheckComplianceAgain = {
          Type = "Choice"
          Choices = [
            {
              Comment       = "All instances are compliant"
              Variable      = "$.Compliant"
              BooleanEquals = true
              Next          = "Succeed"
            },
            {
              Comment       = "There are non-compliant instances"
              Variable      = "$.Compliant"
              BooleanEquals = false
              Next          = "SendAlertToSns"
            },
          ]
        }
        SendAlertToSns = {
          Type     = "Task"
          Resource = module.send_alert_to_sns_lambda.function_arn
          Next     = "Succeed"
        }
        Succeed = {
          Type = "Succeed"
        }
      }
    }
  )
}
