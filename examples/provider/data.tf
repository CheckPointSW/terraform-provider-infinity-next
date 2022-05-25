data "local_file" "schema_validation_file" {
  filename = format("${path.module}/${var.schema_validation_file_name}")
}