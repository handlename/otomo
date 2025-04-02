terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

locals {
  quarifier = "current"
}

resource "terraform_data" "lambda_function" {
  triggers_replace = {
    role_arn      = var.iam_role_arn
    function_name = var.function_name
    function_dir  = var.function_dir
    env_vars      = var.env_vars
  }

  provisioner "local-exec" {
    command     = "make deploy"
    working_dir = self.triggers_replace.function_dir
    when        = create
    environment = merge(
      { ROLE_ARN = self.triggers_replace.role_arn },
      var.env_vars,
    )
  }

  provisioner "local-exec" {
    command     = "make destroy"
    working_dir = self.triggers_replace.function_dir
    when        = destroy
  }
}

# data "aws_lambda_function_url" "url" {
#   function_name = var.function_name
# }
