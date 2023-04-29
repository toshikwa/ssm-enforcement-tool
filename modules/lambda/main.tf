locals {
  lambda_dir = "assets/functions"
}
variable "app" {
  type = string
}
variable "name" {
  type = string
}
variable "policy_actions" {
  type = list(string)
}
variable "envs" {
  type = map(string)
}

# >>>>>>>>>> IAM >>>>>>>>>> #
resource "aws_iam_role" "default" {
  name = "${var.app}_${var.name}_function_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}
data "aws_iam_policy_document" "default" {
  statement {
    effect    = "Allow"
    actions   = var.policy_actions
    resources = ["*"]
  }
}
resource "aws_iam_policy" "default" {
  name   = "${var.app}_${var.name}_function_policy"
  policy = data.aws_iam_policy_document.default.json
}
resource "aws_iam_role_policy_attachment" "default" {
  role       = aws_iam_role.default.name
  policy_arn = aws_iam_policy.default.arn
}
# <<<<<<<<<< IAM <<<<<<<<<< #

# >>>>>>>>>> Lambda >>>>>>>>>> #
data "archive_file" "lambda" {
  type        = "zip"
  source_file = "${local.lambda_dir}/bin/${var.name}"
  output_path = "${local.lambda_dir}/archive/${var.name}.zip"
}
resource "aws_lambda_function" "default" {
  function_name    = "${var.app}_${var.name}_function"
  filename         = "${local.lambda_dir}/archive/${var.name}.zip"
  role             = aws_iam_role.default.arn
  handler          = var.name
  runtime          = "go1.x"
  timeout          = "60"
  memory_size      = 128
  source_code_hash = data.archive_file.lambda.output_base64sha256
  environment {
    variables = var.envs
  }
}
# <<<<<<<<<< Lambda <<<<<<<<<< #

output "function_name" {
  value = aws_lambda_function.default.function_name
}

output "function_arn" {
  value = aws_lambda_function.default.arn
}
