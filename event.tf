# >>>>>>>>>> IAM >>>>>>>>>> #
resource "aws_iam_role_policy" "invoke_sfn_policy" {
  name = "${local.app}_eventbridge_policy"
  role = aws_iam_role.eventbridge_role.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["states:StartExecution"]
        Resource = [aws_sfn_state_machine.default.arn]
      }
    ]
  })
}
data "aws_iam_policy_document" "assume_role_eventbridge" {
  statement {
    actions = [
      "sts:AssumeRole"
    ]
    principals {
      type = "Service"
      identifiers = [
        "events.amazonaws.com"
      ]
    }
  }
}
resource "aws_iam_role" "eventbridge_role" {
  name               = "${local.app}_eventbridge_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_eventbridge.json
}
# <<<<<<<<<< IAM <<<<<<<<<< #

resource "aws_cloudwatch_event_rule" "trigger_sfn" {
  name                = "${local.app}_rule"
  schedule_expression = "cron(0 1 * * ? *)"
}
resource "aws_cloudwatch_event_target" "target_sfn" {
  rule      = aws_cloudwatch_event_rule.trigger_sfn.name
  target_id = "${local.app}_event_target"
  arn       = aws_sfn_state_machine.default.arn
  role_arn  = aws_iam_role.eventbridge_role.arn
}
