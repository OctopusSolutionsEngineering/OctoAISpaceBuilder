// We don't support providers other than the Octopus one
resource "local_file" "foo" {
  content  = "foo!"
  filename = "${path.module}/foo.bar"
}