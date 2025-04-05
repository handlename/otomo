terraform {
  required_version = "1.11.3"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "local" {}
}

provider "aws" {
  region = "us-east-1"

  default_tags {
    tags = {
      Service   = local.service
      ManagedBy = "terraform"
    }
  }
}

data "aws_caller_identity" "current" {}
data "aws_canonical_user_id" "current" {}
data "aws_region" "current" {}
