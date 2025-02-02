resource "aws_dynamodb_table" "users" {
  name           = "${var.app_name}-users"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "user_id"
  
  attribute {
    name = "user_id"
    type = "S"
  }
  attribute {
    name = "email"
    type = "S"
  }
  global_secondary_index {
    name               = "email-index"
    hash_key          = "email"
    projection_type    = "ALL"
  }
  tags = {
    Environment = var.environment
    Region      = var.region_name
  }
}

resource "aws_dynamodb_table" "canvases" {
  name           = "${var.app_name}-canvases"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "user_id"
  range_key      = "canvas_id"

  attribute {
    name = "user_id"
    type = "S"
  }
  attribute {
    name = "canvas_id"
    type = "S"
  }
  ttl {
    attribute_name = "ttl"
    enabled       = true
  }
  tags = {
    Environment = var.environment
    Region      = var.region_name
  }
}

resource "aws_dynamodb_table" "svg_data" {
  name           = "${var.app_name}-svg-data"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "canvas_id"
  range_key      = "svg_id"
  
  attribute {
    name = "canvas_id"
    type = "S"
  }
  attribute {
    name = "svg_id"
    type = "S"
  }
  attribute {
    name = "user_id"
    type = "S"
  }
  attribute {
    name = "created_at"
    type = "S"
  }
  global_secondary_index {
    name               = "user-id-index"
    hash_key          = "user_id"
    range_key         = "created_at"
    projection_type    = "ALL"
  }
  tags = {
    Environment = var.environment
    Region      = var.region_name
  }
}

resource "aws_iam_role" "canvas_app_role" {
  name = "${var.app_name}-app-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "dynamodb_access" {
  name = "dynamodb-access"
  role = aws_iam_role.canvas_app_role.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:DeleteItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = [
          aws_dynamodb_table.users.arn,
          aws_dynamodb_table.canvases.arn,
          aws_dynamodb_table.svg_data.arn,
          "${aws_dynamodb_table.users.arn}/index/*",
          "${aws_dynamodb_table.canvases.arn}/index/*",
          "${aws_dynamodb_table.svg_data.arn}/index/*"
        ]
      }
    ]
  })
}
