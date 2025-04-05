data "aws_iam_policy" "aws_lambda_basic_execution_role" {
  name = "AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "assume_role" {
  for_each = toset([
    "lambda",
  ])

  version = "2012-10-17"

  statement {
    actions = [
      "sts:AssumeRole",
    ]

    effect = "Allow"

    principals {
      type = "Service"

      identifiers = ["${each.value}.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "${local.service}-lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role["lambda"].json
}

resource "aws_iam_role_policy_attachment" "lambda_aws_lambda_basic_execution_role" {
  role       = aws_iam_role.lambda.name
  policy_arn = data.aws_iam_policy.aws_lambda_basic_execution_role.arn
}

resource "aws_iam_role_policy" "lambda" {
  name   = "lambda"
  role   = aws_iam_role.lambda.id
  policy = data.aws_iam_policy_document.lambda.json
}

data "aws_iam_policy_document" "lambda" {
  version = "2012-10-17"

  statement {
    actions = [
      "bedrock:Invoke*",
    ]

    effect = "Allow"

    resources = ["*"]
  }

  statement {
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameter*",
    ]

    effect = "Allow"

    resources = [
      "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.id}:parameter/${local.service}",
      "arn:aws:ssm:${data.aws_region.current.name}:${data.aws_caller_identity.current.id}:parameter/${local.service}/*",
    ]
  }
}
