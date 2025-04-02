variable "function_name" {
  type = string
}

variable "function_dir" {
  type = string
}

variable "iam_role_arn" {
  type = string
}

variable "env_vars" {
  description = "Environment varidales for Lambda function"
  type        = map(string)
  default     = {}
}
