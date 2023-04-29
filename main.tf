data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

terraform {
  required_version = ">= 1.4.0"
  backend "s3" {
    key     = "terraform.tfstate"
    encrypt = true
  }
}

provider "aws" {}

locals {
  app = "ssm_enforcement_tool"
}
