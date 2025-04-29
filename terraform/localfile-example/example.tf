// We don't support providers other than the Octopus one
resource "local_file" "foo" {
  content  = var.octopus_space_id
  filename = "${path.module}/foo.bar"
}

variable "octopus_space_id" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The ID of the Octopus space to populate."
}