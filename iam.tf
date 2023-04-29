data "aws_iam_policy" "ssm_policy" {
  arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}
data "aws_iam_policy_document" "assume_role_ec2" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}
resource "aws_iam_role" "ssm_agent_role" {
  name               = "${local.app}_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_ec2.json
}
resource "aws_iam_instance_profile" "ssm_agent_profile" {
  name = "${local.app}_profile"
  role = aws_iam_role.ssm_agent_role.name
}
