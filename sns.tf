variable "email" {
  type = string
}

resource "aws_sns_topic" "default" {
  name = "${local.app}_topic"
}

resource "aws_sns_topic_subscription" "default" {
  topic_arn = aws_sns_topic.default.arn
  protocol  = "email"
  endpoint  = var.email
}
