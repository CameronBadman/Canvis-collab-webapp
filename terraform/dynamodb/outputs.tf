output "users_table_name" {
  value = aws_dynamodb_table.users.name
}

output "users_table_arn" {
  value = aws_dynamodb_table.users.arn
}

output "canvases_table_name" {
  value = aws_dynamodb_table.canvases.name
}

output "canvases_table_arn" {
  value = aws_dynamodb_table.canvases.arn
}

output "svg_data_table_name" {
  value = aws_dynamodb_table.svg_data.name
}

output "svg_data_table_arn" {
  value = aws_dynamodb_table.svg_data.arn
}

output "iam_role_arn" {
  value = aws_iam_role.canvas_app_role.arn
}
