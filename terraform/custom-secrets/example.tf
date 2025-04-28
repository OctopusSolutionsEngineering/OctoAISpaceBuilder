provider "octopusdeploy" {
  space_id = "${trimspace(var.octopus_space_id)}"
}
provider "shell" {
  interpreter        = ["pwsh", "-Command"]
  enable_parallelism = false
}
provider "external" {
}

terraform {

  required_providers {
    octopusdeploy = { source = "OctopusDeployLabs/octopusdeploy", version = "0.43.0" }
    shell         = { source = "scottwinkler/shell", version = "1.7.10" }
    external      = { source = "hashicorp/external", version = "2.3.4" }
  }
  required_version = ">= 1.6.0"
}

variable "octopus_space_id" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The ID of the Octopus space to populate."
}

data "octopusdeploy_project_groups" "project_group_kubernetes" {
  ids          = null
  partial_name = "${var.project_group_kubernetes_name}"
  skip         = 0
  take         = 1
}
variable "project_group_kubernetes_name" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The name of the project group to lookup"
  default     = "Kubernetes"
}
resource "octopusdeploy_project_group" "project_group_kubernetes" {
  count = "${length(data.octopusdeploy_project_groups.project_group_kubernetes.project_groups) != 0 ? 0 : 1}"
  name  = "${var.project_group_kubernetes_name}"
  lifecycle {
    prevent_destroy = true
  }
}

data "octopusdeploy_environments" "environment_dev" {
  ids          = null
  partial_name = "Dev"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_environment" "environment_dev" {
  count                        = "${length(data.octopusdeploy_environments.environment_dev.environments) != 0 ? 0 : 1}"
  name                         = "Dev"
  description                  = ""
  allow_dynamic_infrastructure = true
  use_guided_failure           = false

  jira_extension_settings {
    environment_type = "development"
  }

  jira_service_management_extension_settings {
    is_enabled = true
  }

  servicenow_extension_settings {
    is_enabled = true
  }
  lifecycle {
    prevent_destroy = true
  }
}

data "octopusdeploy_environments" "environment_development__security_" {
  ids          = null
  partial_name = "Development (Security)"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_environment" "environment_development__security_" {
  count                        = "${length(data.octopusdeploy_environments.environment_development__security_.environments) != 0 ? 0 : 1}"
  name                         = "Development (Security)"
  description                  = "Used to scan the development releases for security issues. This resource is created and managed by the [Octopus Terraform provider](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs). The Terraform files can be found in the [GitHub repo](https://github.com/mcasperson/AppBuilder-EKS)."
  allow_dynamic_infrastructure = true
  use_guided_failure           = false

  jira_extension_settings {
    environment_type = "unmapped"
  }

  jira_service_management_extension_settings {
    is_enabled = false
  }

  servicenow_extension_settings {
    is_enabled = false
  }
  lifecycle {
    prevent_destroy = true
  }
}

data "octopusdeploy_environments" "environment_production__app_" {
  ids          = null
  partial_name = "Production (App)"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_environment" "environment_production__app_" {
  count                        = "${length(data.octopusdeploy_environments.environment_production__app_.environments) != 0 ? 0 : 1}"
  name                         = "Production (App)"
  description                  = "The production environment."
  allow_dynamic_infrastructure = true
  use_guided_failure           = false

  jira_extension_settings {
    environment_type = "unmapped"
  }

  jira_service_management_extension_settings {
    is_enabled = false
  }

  servicenow_extension_settings {
    is_enabled = false
  }
  lifecycle {
    prevent_destroy = true
  }
}

data "octopusdeploy_environments" "environment_production__security_" {
  ids          = null
  partial_name = "Production (Security)"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_environment" "environment_production__security_" {
  count                        = "${length(data.octopusdeploy_environments.environment_production__security_.environments) != 0 ? 0 : 1}"
  name                         = "Production (Security)"
  description                  = "Used to scan the productions releases for security issues. This resource is created and managed by the [Octopus Terraform provider](https://registry.terraform.io/providers/OctopusDeployLabs/octopusdeploy/latest/docs). The Terraform files can be found in the [GitHub repo](https://github.com/mcasperson/AppBuilder-EKS)."
  allow_dynamic_infrastructure = true
  use_guided_failure           = false

  jira_extension_settings {
    environment_type = "unmapped"
  }

  jira_service_management_extension_settings {
    is_enabled = false
  }

  servicenow_extension_settings {
    is_enabled = false
  }
  lifecycle {
    prevent_destroy = true
  }
}

data "octopusdeploy_lifecycles" "lifecycle_application_and_security" {
  ids          = null
  partial_name = "Application and Security"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_lifecycle" "lifecycle_application_and_security" {
  count       = "${length(data.octopusdeploy_lifecycles.lifecycle_application_and_security.lifecycles) != 0 ? 0 : 1}"
  name        = "Application and Security"
  description = ""

  phase {
    automatic_deployment_targets          = []
    optional_deployment_targets           = ["${length(data.octopusdeploy_environments.environment_dev.environments) != 0 ? data.octopusdeploy_environments.environment_dev.environments[0].id : octopusdeploy_environment.environment_dev[0].id}"]
    name                                  = "Development"
    is_optional_phase                     = false
    minimum_environments_before_promotion = 0
  }
  phase {
    automatic_deployment_targets          = ["${length(data.octopusdeploy_environments.environment_development__security_.environments) != 0 ? data.octopusdeploy_environments.environment_development__security_.environments[0].id : octopusdeploy_environment.environment_development__security_[0].id}"]
    optional_deployment_targets           = []
    name                                  = "Dveelopment Security"
    is_optional_phase                     = false
    minimum_environments_before_promotion = 0
  }
  phase {
    automatic_deployment_targets          = []
    optional_deployment_targets           = ["${length(data.octopusdeploy_environments.environment_production__app_.environments) != 0 ? data.octopusdeploy_environments.environment_production__app_.environments[0].id : octopusdeploy_environment.environment_production__app_[0].id}"]
    name                                  = "Production"
    is_optional_phase                     = false
    minimum_environments_before_promotion = 0
  }
  phase {
    automatic_deployment_targets          = ["${length(data.octopusdeploy_environments.environment_production__security_.environments) != 0 ? data.octopusdeploy_environments.environment_production__security_.environments[0].id : octopusdeploy_environment.environment_production__security_[0].id}"]
    optional_deployment_targets           = []
    name                                  = "Production Security"
    is_optional_phase                     = false
    minimum_environments_before_promotion = 0
  }

  release_retention_policy {
    quantity_to_keep = 30
    unit             = "Days"
  }

  tentacle_retention_policy {
    quantity_to_keep = 30
    unit             = "Days"
  }
  lifecycle {
    prevent_destroy = true
  }
}

data "octopusdeploy_environments" "environment_test" {
  ids          = null
  partial_name = "Test"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_environment" "environment_test" {
  count                        = "${length(data.octopusdeploy_environments.environment_test.environments) != 0 ? 0 : 1}"
  name                         = "Test"
  description                  = ""
  allow_dynamic_infrastructure = true
  use_guided_failure           = false

  jira_extension_settings {
    environment_type = "unmapped"
  }

  jira_service_management_extension_settings {
    is_enabled = false
  }

  servicenow_extension_settings {
    is_enabled = false
  }
  lifecycle {
    prevent_destroy = true
  }
}

data "octopusdeploy_environments" "environment_production" {
  ids          = null
  partial_name = "Production"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_environment" "environment_production" {
  count                        = "${length(data.octopusdeploy_environments.environment_production.environments) != 0 ? 0 : 1}"
  name                         = "Production"
  description                  = ""
  allow_dynamic_infrastructure = true
  use_guided_failure           = false

  jira_extension_settings {
    environment_type = "unmapped"
  }

  jira_service_management_extension_settings {
    is_enabled = false
  }

  servicenow_extension_settings {
    is_enabled = false
  }
  lifecycle {
    prevent_destroy = true
  }
}

data "octopusdeploy_lifecycles" "lifecycle_default_lifecycle" {
  ids          = null
  partial_name = "Default Lifecycle"
  skip         = 0
  take         = 1
  lifecycle {
    postcondition {
      error_message = "Failed to resolve a lifecycle called \"Default Lifecycle\". This resource must exist in the space before this Terraform configuration is applied."
      condition     = length(self.lifecycles) != 0
    }
  }
}

data "octopusdeploy_channels" "channel_my_k8s_project_2_default" {
  ids          = []
  partial_name = "Default"
  skip         = 0
  take         = 1
}

data "octopusdeploy_feeds" "feed_octopus_server__built_in_" {
  feed_type    = "BuiltIn"
  ids          = null
  partial_name = ""
  skip         = 0
  take         = 1
  lifecycle {
    postcondition {
      error_message = "Failed to resolve a feed called \"BuiltIn\". This resource must exist in the space before this Terraform configuration is applied."
      condition     = length(self.feeds) != 0
    }
  }
}

data "octopusdeploy_feeds" "feed_github_container_registry" {
  feed_type    = "Docker"
  ids          = null
  partial_name = "GitHub Container Registry"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_docker_container_registry" "feed_github_container_registry" {
  count                                = "${length(data.octopusdeploy_feeds.feed_github_container_registry.feeds) != 0 ? 0 : 1}"
  name                                 = "GitHub Container Registry"
  registry_path                        = ""
  api_version                          = "v2"
  feed_uri                             = "https://ghcr.io"
  package_acquisition_location_options = ["ExecutionTarget", "NotAcquired"]
  lifecycle {
    ignore_changes  = [password]
    prevent_destroy = true
  }
}

data "octopusdeploy_worker_pools" "workerpool_default_worker_pool" {
  ids          = null
  partial_name = "Default Worker Pool"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_dynamic_worker_pool" "workerpool_default_worker_pool" {
  count       = "${length(data.octopusdeploy_worker_pools.workerpool_default_worker_pool.worker_pools) != 0 ? 0 : 1}"
  name        = "Default Worker Pool"
  description = "Default pool of workers from the Dynamic Worker service"
  is_default  = false
  worker_type = "WindowsDefault"
  lifecycle {
    prevent_destroy = true
  }
}

data "octopusdeploy_worker_pools" "workerpool_hosted_ubuntu" {
  ids          = null
  partial_name = "Hosted Ubuntu"
  skip         = 0
  take         = 1
  lifecycle {
    postcondition {
      error_message = "Failed to resolve a worker pool called \"Hosted Ubuntu\". This resource must exist in the space before this Terraform configuration is applied."
      condition     = length(self.worker_pools) != 0
    }
  }
}

resource "octopusdeploy_deployment_process" "deployment_process_my_k8s_project_2" {
  project_id = "${length(data.octopusdeploy_projects.project_my_k8s_project_2.projects) != 0 ? data.octopusdeploy_projects.project_my_k8s_project_2.projects[0].id : octopusdeploy_project.project_my_k8s_project_2[0].id}"

  step {
    condition           = "Success"
    name                = "Deploy a Kubernetes Web App via YAML"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.KubernetesDeployRawYaml"
      name                               = "Deploy a Kubernetes Web App via YAML"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = true
      is_required                        = false
      worker_pool_id                     = "${data.octopusdeploy_worker_pools.workerpool_hosted_ubuntu.worker_pools[0].id}"
      properties                         = {
        "Octopus.Action.Kubernetes.ResourceStatusCheck" = "True"
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.KubernetesContainers.Namespace" = "#{Octopus.Environment.Name | ToLower}"
        "Octopus.Action.Kubernetes.DeploymentTimeout" = "180"
        "Octopus.Action.Kubernetes.ServerSideApply.ForceConflicts" = "True"
        "Octopus.Action.RunOnServer" = "true"
        "Octopus.Action.Kubernetes.ServerSideApply.Enabled" = "True"
        "Octopus.Action.KubernetesContainers.CustomResourceYaml" = "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: octopub\n  labels:\n    app: octopub\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      app: octopub\n  template:\n    metadata:\n      labels:\n        app: octopub\n    spec:\n      containers:\n      - name: octopub\n        image: octopussamples/octopub-selfcontained\n        ports:\n        - containerPort: 8080\n        resources:\n          limits:\n            cpu: \"1\"\n            memory: \"512Mi\"\n          requests:\n            cpu: \"0.5\"\n            memory: \"256Mi\"\n        livenessProbe:\n          httpGet:\n            path: /health/products\n            port: 8080\n          initialDelaySeconds: 30\n          periodSeconds: 10\n        readinessProbe:\n          httpGet:\n            path: /health/products\n            port: 8080\n          initialDelaySeconds: 5\n          periodSeconds: 5\n"
        "OctopusUseBundledTooling" = "False"
      }

      container {
        feed_id = "${length(data.octopusdeploy_feeds.feed_github_container_registry.feeds) != 0 ? data.octopusdeploy_feeds.feed_github_container_registry.feeds[0].id : octopusdeploy_docker_container_registry.feed_github_container_registry[0].id}"
        image   = "ghcr.io/octopusdeploylabs/k8s-workertools"
      }

      environments          = []
      excluded_environments = ["${length(data.octopusdeploy_environments.environment_development__security_.environments) != 0 ? data.octopusdeploy_environments.environment_development__security_.environments[0].id : octopusdeploy_environment.environment_development__security_[0].id}"]
      channels              = []
      tenant_tags           = []
      features              = []
    }

    properties   = {}
    target_roles = ["Kubernetes"]
  }
  depends_on = []
}

variable "variable_42a2de0c9187a3d773ca5ee10a26490b4e720a0edd470e2907fbe2a77f633531_sensitive_value" {
  type        = string
  nullable    = true
  sensitive   = true
  description = "The secret variable value associated with the variable SecretVariable"
  default     = "A custom secret value"
}
resource "octopusdeploy_variable" "my_k8s_project_2_secretvariable_1" {
  count           = "${length(data.octopusdeploy_projects.project_my_k8s_project_2.projects) != 0 ? 0 : 1}"
  owner_id        = "${length(data.octopusdeploy_projects.project_my_k8s_project_2.projects) == 0 ?octopusdeploy_project.project_my_k8s_project_2[0].id : data.octopusdeploy_projects.project_my_k8s_project_2.projects[0].id}"
  name            = "SecretVariable"
  type            = "Sensitive"
  is_sensitive    = true
  sensitive_value = var.variable_42a2de0c9187a3d773ca5ee10a26490b4e720a0edd470e2907fbe2a77f633531_sensitive_value
  lifecycle {
    ignore_changes  = [sensitive_value]
    prevent_destroy = true
  }
  depends_on = []
}

variable "project_my_k8s_project_2_name" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The name of the project exported from My K8s Project 2"
  default     = "My K8s Project 2"
}
variable "project_my_k8s_project_2_description_prefix" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "An optional prefix to add to the project description for the project My K8s Project 2"
  default     = ""
}
variable "project_my_k8s_project_2_description_suffix" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "An optional suffix to add to the project description for the project My K8s Project 2"
  default     = ""
}
variable "project_my_k8s_project_2_description" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The description of the project exported from My K8s Project 2"
  default     = "This project provides an example Kubernetes deployment using YAML, Cloud Target Discovery, and SBOM scanning to an AWS EKS Kubernetes cluster."
}
variable "project_my_k8s_project_2_tenanted" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The tenanted setting for the project Untenanted"
  default     = "Untenanted"
}
data "octopusdeploy_projects" "project_my_k8s_project_2" {
  ids          = null
  partial_name = "${var.project_my_k8s_project_2_name}"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_project" "project_my_k8s_project_2" {
  count                                = "${length(data.octopusdeploy_projects.project_my_k8s_project_2.projects) != 0 ? 0 : 1}"
  name                                 = "${var.project_my_k8s_project_2_name}"
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  discrete_channel_release             = false
  is_disabled                          = false
  is_version_controlled                = false
  lifecycle_id                         = "${length(data.octopusdeploy_lifecycles.lifecycle_application_and_security.lifecycles) != 0 ? data.octopusdeploy_lifecycles.lifecycle_application_and_security.lifecycles[0].id : octopusdeploy_lifecycle.lifecycle_application_and_security[0].id}"
  project_group_id                     = "${length(data.octopusdeploy_project_groups.project_group_kubernetes.project_groups) != 0 ? data.octopusdeploy_project_groups.project_group_kubernetes.project_groups[0].id : octopusdeploy_project_group.project_group_kubernetes[0].id}"
  included_library_variable_sets       = []
  tenanted_deployment_participation    = "${var.project_my_k8s_project_2_tenanted}"

  connectivity_policy {
    allow_deployments_to_no_targets = true
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "None"
  }

  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.NextPatch}"
  }
  description = "${var.project_my_k8s_project_2_description_prefix}${var.project_my_k8s_project_2_description}${var.project_my_k8s_project_2_description_suffix}"
  lifecycle {
    prevent_destroy = true
  }
}

