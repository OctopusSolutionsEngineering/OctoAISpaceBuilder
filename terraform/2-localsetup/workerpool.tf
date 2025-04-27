resource "octopusdeploy_static_worker_pool" "workerpool_docker" {
  name        = "Hosted Ubuntu"
  description = "Emulates the cloud hosted ubuntu worker pool"
  is_default  = false
  sort_order  = 3
}
