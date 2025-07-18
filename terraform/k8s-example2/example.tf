provider "octopusdeploy" {
  space_id = "${trimspace(var.octopus_space_id)}"
}

terraform {

  required_providers {
    octopusdeploy = { source = "OctopusDeploy/octopusdeploy", version = "1.0.1" }
  }
  required_version = ">= 1.6.0"
}

variable "octopus_space_id" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The ID of the Octopus space to populate."
}

output "octopus_space_id" {
  value = "${var.octopus_space_id}"
}
data "octopusdeploy_spaces" "octopus_space_name" {
  ids  = ["${var.octopus_space_id}"]
  skip = 0
  take = 1
}
output "octopus_space_name" {
  value = "${data.octopusdeploy_spaces.octopus_space_name.spaces[0].name}"
}

variable "project_group_default_project_group_name" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The name of the project group to lookup"
  default     = "Default Project Group"
}
data "octopusdeploy_project_groups" "project_group_default_project_group" {
  ids          = null
  partial_name = "${var.project_group_default_project_group_name}"
  skip         = 0
  take         = 1
  lifecycle {
    postcondition {
      error_message = "Failed to resolve a project group called $${var.project_group_default_project_group_name}. This resource must exist in the space before this Terraform configuration is applied."
      condition     = length(self.project_groups) != 0
    }
  }
}

data "octopusdeploy_environments" "environment_development" {
  ids          = null
  partial_name = "Development"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_environment" "environment_development" {
  count                        = "${length(data.octopusdeploy_environments.environment_development.environments) != 0 ? 0 : 1}"
  name                         = "Development"
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

data "octopusdeploy_environments" "environment_security" {
  ids          = null
  partial_name = "Security"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_environment" "environment_security" {
  count                        = "${length(data.octopusdeploy_environments.environment_security.environments) != 0 ? 0 : 1}"
  name                         = "Security"
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

data "octopusdeploy_lifecycles" "lifecycle_security" {
  ids          = null
  partial_name = "Security"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_lifecycle" "lifecycle_security" {
  count       = "${length(data.octopusdeploy_lifecycles.lifecycle_security.lifecycles) != 0 ? 0 : 1}"
  name        = "Security"
  description = "This lifecycle automatically deploys to the Development environment, progresses through the Test and Production environments, and then automatically deploys to the Security environment. The Security environment is used to scan SBOMs for any vulnerabilities and deployments to the Security environment are initiated by triggers on a daily basis."

  phase {
    automatic_deployment_targets          = ["${length(data.octopusdeploy_environments.environment_development.environments) != 0 ? data.octopusdeploy_environments.environment_development.environments[0].id : octopusdeploy_environment.environment_development[0].id}"]
    optional_deployment_targets           = []
    name                                  = "Development"
    is_optional_phase                     = false
    minimum_environments_before_promotion = 0
  }
  phase {
    automatic_deployment_targets          = []
    optional_deployment_targets           = ["${length(data.octopusdeploy_environments.environment_test.environments) != 0 ? data.octopusdeploy_environments.environment_test.environments[0].id : octopusdeploy_environment.environment_test[0].id}"]
    name                                  = "Test"
    is_optional_phase                     = false
    minimum_environments_before_promotion = 0
  }
  phase {
    automatic_deployment_targets          = []
    optional_deployment_targets           = ["${length(data.octopusdeploy_environments.environment_production.environments) != 0 ? data.octopusdeploy_environments.environment_production.environments[0].id : octopusdeploy_environment.environment_production[0].id}"]
    name                                  = "Production"
    is_optional_phase                     = false
    minimum_environments_before_promotion = 0
  }
  phase {
    automatic_deployment_targets          = ["${length(data.octopusdeploy_environments.environment_security.environments) != 0 ? data.octopusdeploy_environments.environment_security.environments[0].id : octopusdeploy_environment.environment_security[0].id}"]
    optional_deployment_targets           = []
    name                                  = "Security"
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

data "octopusdeploy_channels" "channel_kubernetes_web_app_default" {
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

data "octopusdeploy_feeds" "feed_ghcr_anonymous" {
  feed_type    = "Docker"
  ids          = null
  partial_name = "GHCR Anonymous"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_docker_container_registry" "feed_ghcr_anonymous" {
  count                                = "${length(data.octopusdeploy_feeds.feed_ghcr_anonymous.feeds) != 0 ? 0 : 1}"
  name                                 = "GHCR Anonymous"
  registry_path                        = ""
  api_version                          = "v2"
  feed_uri                             = "https://ghcrfacade-a6awccayfpcpg4cg.eastus-01.azurewebsites.net"
  package_acquisition_location_options = ["ExecutionTarget", "NotAcquired"]
  lifecycle {
    ignore_changes  = [password]
    prevent_destroy = true
  }
}

variable "project_kubernetes_web_app_step_deploy_a_kubernetes_web_app_via_yaml_package_octopub_selfcontained_packageid" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The package ID for the package named octopub-selfcontained from step Deploy a Kubernetes Web App via YAML in project Kubernetes Web App"
  default     = "octopussolutionsengineering/octopub-selfcontained"
}
resource "octopusdeploy_deployment_process" "deployment_process_kubernetes_web_app" {
  project_id = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) != 0 ? data.octopusdeploy_projects.project_kubernetes_web_app.projects[0].id : octopusdeploy_project.project_kubernetes_web_app[0].id}"

  step {
    condition           = "Success"
    name                = "Approve Production Deployment"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Manual"
      name                               = "Approve Production Deployment"
      notes                              = "This manual intervention is used to provide a prompt before a production deployment."
      condition                          = "Success"
      run_on_server                      = false
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = false
      worker_pool_id                     = ""
      properties                         = {
        "Octopus.Action.Manual.BlockConcurrentDeployments" = "False"
        "Octopus.Action.Manual.Instructions" = "Do you approve the production deployment?"
        "Octopus.Action.RunOnServer" = "false"
      }
      environments                       = ["${length(data.octopusdeploy_environments.environment_production.environments) != 0 ? data.octopusdeploy_environments.environment_production.environments[0].id : octopusdeploy_environment.environment_production[0].id}"]
      excluded_environments              = []
      channels                           = []
      tenant_tags                        = []
      features                           = []
    }

    properties   = {}
    target_roles = []
  }
  step {
    condition           = "Success"
    name                = "Deploy a Kubernetes Web App via YAML"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.KubernetesDeployRawYaml"
      name                               = "Deploy a Kubernetes Web App via YAML"
      notes                              = "This step deploys a Kubernetes Deployment resource defined as raw YAML."
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = true
      is_required                        = false
      worker_pool_variable               = "Octopus.Project.WorkerPool"
      properties                         = {
        "Octopus.Action.Kubernetes.ServerSideApply.Enabled" = "True"
        "Octopus.Action.RunOnServer" = "true"
        "Octopus.Action.Kubernetes.ResourceStatusCheck" = "True"
        "Octopus.Action.Kubernetes.DeploymentTimeout" = "180"
        "OctopusUseBundledTooling" = "False"
        "Octopus.Action.Kubernetes.ServerSideApply.ForceConflicts" = "True"
        "Octopus.Action.KubernetesContainers.CustomResourceYaml" = "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: #{Kubernetes.Deployment.Name}\n  labels:\n    app: #{Kubernetes.Deployment.Name}\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      app: #{Kubernetes.Deployment.Name}\n  template:\n    metadata:\n      labels:\n        app: #{Kubernetes.Deployment.Name}\n    spec:\n      containers:\n      - name: octopub\n        # The image is sourced from the package reference. This allows the package version to be selected\n        # at release creation time.\n        image: \"ghcr.io/#{Octopus.Action.Package[octopub-selfcontained].PackageId}:#{Octopus.Action.Package[octopub-selfcontained].PackageVersion}\"\n        ports:\n        - containerPort: 8080\n        resources:\n          limits:\n            cpu: \"1\"\n            memory: \"512Mi\"\n          requests:\n            cpu: \"0.5\"\n            memory: \"256Mi\"\n        livenessProbe:\n          httpGet:\n            path: /health/products\n            port: 8080\n          initialDelaySeconds: 30\n          periodSeconds: 10\n        readinessProbe:\n          httpGet:\n            path: /health/products\n            port: 8080\n          initialDelaySeconds: 5\n          periodSeconds: 5"
        "Octopus.Action.KubernetesContainers.Namespace" = "#{Octopus.Environment.Name | ToLower}"
        "Octopus.Action.Script.ScriptSource" = "Inline"
      }

      container {
        feed_id = "${length(data.octopusdeploy_feeds.feed_github_container_registry.feeds) != 0 ? data.octopusdeploy_feeds.feed_github_container_registry.feeds[0].id : octopusdeploy_docker_container_registry.feed_github_container_registry[0].id}"
        image   = "ghcr.io/octopusdeploylabs/k8s-workertools"
      }

      environments          = []
      excluded_environments = ["${length(data.octopusdeploy_environments.environment_security.environments) != 0 ? data.octopusdeploy_environments.environment_security.environments[0].id : octopusdeploy_environment.environment_security[0].id}"]
      channels              = []
      tenant_tags           = []

      package {
        name                      = "octopub-selfcontained"
        package_id                = "${var.project_kubernetes_web_app_step_deploy_a_kubernetes_web_app_via_yaml_package_octopub_selfcontained_packageid}"
        acquisition_location      = "NotAcquired"
        extract_during_deployment = false
        feed_id                   = "${length(data.octopusdeploy_feeds.feed_ghcr_anonymous.feeds) != 0 ? data.octopusdeploy_feeds.feed_ghcr_anonymous.feeds[0].id : octopusdeploy_docker_container_registry.feed_ghcr_anonymous[0].id}"
        properties                = { Extract = "False", Purpose = "DockerImageReference", SelectionMode = "immediate" }
      }
      features = []
    }

    properties   = {}
    target_roles = ["Kubernetes"]
  }
  step {
    condition           = "Success"
    name                = "Print Message When no Targets"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Script"
      name                               = "Print Message When no Targets"
      notes                              = "This step detects when the previous steps were skipped due to no targets being defined and prints a message with a link to documentation for creating targets."
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = false
      worker_pool_variable               = "Octopus.Project.WorkerPool"
      properties                         = {
        "Octopus.Action.Script.ScriptBody" = "# The variable to check must be in the format\n# #{Octopus.Step[\u003cname of the step that deploys the web app\u003e].Status.Code}\nif (\"#{Octopus.Step[Deploy a Kubernetes Web App via YAML].Status.Code}\" -eq \"Abandoned\") {\n  Write-Highlight \"To complete the deployment you must have a Kubernetes target with the tag 'Kubernetes'.\"\n  Write-Highlight \"We recoomend the [Kubernetes Agent](https://octopus.com/docs/kubernetes/targets/kubernetes-agent).\"\n}"
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.Syntax" = "PowerShell"
        "OctopusUseBundledTooling" = "False"
        "Octopus.Action.RunOnServer" = "true"
      }
      environments                       = []
      excluded_environments              = ["${length(data.octopusdeploy_environments.environment_security.environments) != 0 ? data.octopusdeploy_environments.environment_security.environments[0].id : octopusdeploy_environment.environment_security[0].id}"]
      channels                           = []
      tenant_tags                        = []
      features                           = []
    }

    properties   = {}
    target_roles = []
  }
  step {
    condition           = "Success"
    name                = "Scan for Vulnerabilities"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Script"
      name                               = "Scan for Vulnerabilities"
      notes                              = "This step extracts the Docker image, finds any bom.json files, and scans them for vulnerabilities using Trivy. \n\nThis step is expected to be run with each deployment to ensure vulnerabilities are discovered as early as possible. \n\nIt is also run daily via a project trigger that reruns the deployment in the Security environment. This allows unknown vulnerabilities to be discovered after a production deployment."
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = false
      worker_pool_variable               = "Octopus.Project.WorkerPool"
      properties                         = {
        "Octopus.Action.Script.ScriptBody" = "Write-Host \"Pulling Trivy Docker Image\"\nWrite-Host \"##octopus[stdout-verbose]\"\ndocker pull ghcr.io/aquasecurity/trivy\nWrite-Host \"##octopus[stdout-default]\"\n\n# Extract files from the Docker image\nWrite-Host \"Extracting files from Docker image...\"\n# The IMAGE_NAME variable must be in the format\n# ghcr.io/#{Octopus.Action[\u003cName of the step applying the Kubernetes deployment resource\u003e].Package[\u003cName of the referenced package\u003e].PackageId}:#{Octopus.Action[\u003cName of the step applying the Kubernetes deployment resource\u003e].Package[Name of the referenced package\u003e].PackageVersion}\n$IMAGE_NAME = \"ghcr.io/#{Octopus.Action[Deploy a Kubernetes Web App via YAML].Package[octopub-selfcontained].PackageId}:#{Octopus.Action[Deploy a Kubernetes Web App via YAML].Package[octopub-selfcontained].PackageVersion}\"\n$CONTAINER_NAME = \"temporary\"\n\n# Create directory if it doesn't exist\nNew-Item -Path \"extracted_files\" -ItemType Directory -Force | Out-Null\n\nWrite-Host \"##octopus[stdout-verbose]\"\n# Create a temporary container\ndocker create --name $CONTAINER_NAME $IMAGE_NAME 2\u003e\u00261\n\n# Export the container's root filesystem\n# PowerShell equivalent for tar extraction\ndocker export $CONTAINER_NAME | tar \"--wildcards-match-slash\" \"--wildcards\" \"*/bom.json\" -C \"extracted_files\" -xvf - 2\u003e\u00261\n\n# Remove the temporary container\ndocker rm $CONTAINER_NAME 2\u003e\u00261\nWrite-Host \"##octopus[stdout-default]\"\n\n$TIMESTAMP = [DateTimeOffset]::Now.ToUnixTimeMilliseconds()\n$SUCCESS = 0\n\n# Find all bom.json files\n$bomFiles = Get-ChildItem -Path \".\" -Filter \"bom.json\" -Recurse -File\n\nforeach ($file in $bomFiles) {\n    Write-Host \"Scanning $($file.FullName)\"\n\n    # Delete any existing report file\n    if (Test-Path \"$PWD/depscan-bom.json\") {\n        Remove-Item \"$PWD/depscan-bom.json\" -Force\n    }\n\n    # Generate the report, capturing the output\n    try {\n        $OUTPUT = docker run --rm -v \"$($file.FullName):/input/$($file.Name)\" ghcr.io/aquasecurity/trivy sbom -q \"/input/$($file.Name)\"\n        $exitCode = $LASTEXITCODE\n    }\n    catch {\n        $OUTPUT = $_.Exception.Message\n        $exitCode = 1\n    }\n\n    # Run again to generate the JSON output\n    docker run --rm -v \"$${PWD}:/output\" -v \"$($file.FullName):/input/$($file.Name)\" ghcr.io/aquasecurity/trivy sbom -q -f json -o /output/depscan-bom.json \"/input/$($file.Name)\"\n\n    # Octopus Deploy artifact\n    New-OctopusArtifact \"$PWD/depscan-bom.json\"\n\n    # Parse JSON output to count vulnerabilities\n    $jsonContent = Get-Content -Path \"depscan-bom.json\" | ConvertFrom-Json\n    $CRITICAL = ($jsonContent.Results | ForEach-Object { $_.Vulnerabilities } | Where-Object { $_.Severity -eq \"CRITICAL\" }).Count\n    $HIGH = ($jsonContent.Results | ForEach-Object { $_.Vulnerabilities } | Where-Object { $_.Severity -eq \"HIGH\" }).Count\n\n    if (\"#{Octopus.Environment.Name}\" -eq \"Security\") {\n        Write-Highlight \"🟥 $CRITICAL critical vulnerabilities\"\n        Write-Highlight \"🟧 $HIGH high vulnerabilities\"\n    }\n\n    # Set success to 1 if exit code is not zero\n    if ($exitCode -ne 0) {\n        $SUCCESS = 1\n    }\n\n    # Print the output\n    $OUTPUT | ForEach-Object {\n        if ($_.Length -gt 0) {\n            Write-Host $_\n        }\n    }\n}\n\n# Cleanup\nfor ($i = 1; $i -le 10; $i++) {\n    try {\n        if (Test-Path \"bundle\") {\n            Set-ItemProperty -Path \"bundle\" -Name IsReadOnly -Value $false -Recurse -ErrorAction SilentlyContinue\n            Remove-Item -Path \"bundle\" -Recurse -Force -ErrorAction Stop\n            break\n        }\n    }\n    catch {\n        Write-Host \"Attempting to clean up files\"\n        Start-Sleep -Seconds 1\n    }\n}\n\n# Set Octopus variable\nSet-OctopusVariable -Name \"VerificationResult\" -Value $SUCCESS\n\nexit 0"
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.Syntax" = "PowerShell"
        "OctopusUseBundledTooling" = "False"
        "Octopus.Action.RunOnServer" = "true"
      }
      environments                       = []
      excluded_environments              = []
      channels                           = []
      tenant_tags                        = []
      features                           = []
    }

    properties   = {}
    target_roles = []
  }
  depends_on = []
}

variable "variable_06e1298a6b448c42c52466647770012776889b9fd03777616984b6737d0d4a9c_value" {
  type        = string
  nullable    = true
  sensitive   = false
  description = "The value associated with the variable Kubernetes.Deployment.Name"
  default     = "octopub"
}
resource "octopusdeploy_variable" "kubernetes_web_app_kubernetes_deployment_name_1" {
  count        = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) != 0 ? 0 : 1}"
  owner_id     = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) == 0 ?octopusdeploy_project.project_kubernetes_web_app[0].id : data.octopusdeploy_projects.project_kubernetes_web_app.projects[0].id}"
  value        = "${var.variable_06e1298a6b448c42c52466647770012776889b9fd03777616984b6737d0d4a9c_value}"
  name         = "Kubernetes.Deployment.Name"
  type         = "String"
  description  = "This is the name of the Kubernetes deployment. It is exposed as a variable to allow the name to be shared between the deployment process and the runbooks."
  is_sensitive = false
  lifecycle {
    ignore_changes  = [sensitive_value]
    prevent_destroy = true
  }
  depends_on = []
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

resource "octopusdeploy_variable" "kubernetes_web_app_octopus_project_workerpool_1" {
  count        = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) != 0 ? 0 : 1}"
  owner_id     = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) == 0 ?octopusdeploy_project.project_kubernetes_web_app[0].id : data.octopusdeploy_projects.project_kubernetes_web_app.projects[0].id}"
  value        = "${data.octopusdeploy_worker_pools.workerpool_hosted_ubuntu.worker_pools[0].id}"
  name         = "Octopus.Project.WorkerPool"
  type         = "WorkerPool"
  description  = "Using a variable to define the worker pool used by steps allows the worker pool to be easily changed in a central location."
  is_sensitive = false
  lifecycle {
    ignore_changes  = [sensitive_value]
    prevent_destroy = true
  }
  depends_on = []
}

variable "variable_0c5050f146b8c3c60013cbd83e1f4c5146e0cbd70d63e18d2fed539a3717c129_value" {
  type        = string
  nullable    = true
  sensitive   = false
  description = "The value associated with the variable OctopusPrintEvaluatedVariables"
  default     = "False"
}
resource "octopusdeploy_variable" "kubernetes_web_app_octopusprintevaluatedvariables_1" {
  count        = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) != 0 ? 0 : 1}"
  owner_id     = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) == 0 ?octopusdeploy_project.project_kubernetes_web_app[0].id : data.octopusdeploy_projects.project_kubernetes_web_app.projects[0].id}"
  value        = "${var.variable_0c5050f146b8c3c60013cbd83e1f4c5146e0cbd70d63e18d2fed539a3717c129_value}"
  name         = "OctopusPrintEvaluatedVariables"
  type         = "String"
  description  = "Set this variable to true to log the evaluated value of variables available at the beginning of each step in the deployment as Verbose messages. See https://octopus.com/docs/support/debug-problems-with-octopus-variables for more details."
  is_sensitive = false
  lifecycle {
    ignore_changes  = [sensitive_value]
    prevent_destroy = true
  }
  depends_on = []
}

variable "variable_661cf577e84dfa02eac136deec93000f9cee60a23514aaf7d40f3b5de3994542_value" {
  type        = string
  nullable    = true
  sensitive   = false
  description = "The value associated with the variable OctopusPrintVariables"
  default     = "False"
}
resource "octopusdeploy_variable" "kubernetes_web_app_octopusprintvariables_1" {
  count        = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) != 0 ? 0 : 1}"
  owner_id     = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) == 0 ?octopusdeploy_project.project_kubernetes_web_app[0].id : data.octopusdeploy_projects.project_kubernetes_web_app.projects[0].id}"
  value        = "${var.variable_661cf577e84dfa02eac136deec93000f9cee60a23514aaf7d40f3b5de3994542_value}"
  name         = "OctopusPrintVariables"
  type         = "String"
  description  = "Set this variable to true to log the variables available at the beginning of each step in the deployment as Verbose messages. See https://octopus.com/docs/support/debug-problems-with-octopus-variables for more details."
  is_sensitive = false
  lifecycle {
    ignore_changes  = [sensitive_value]
    prevent_destroy = true
  }
  depends_on = []
}

data "octopusdeploy_projects" "projecttrigger_kubernetes_web_app_daily_security_scan" {
  ids          = null
  partial_name = "Kubernetes Web App"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_project_scheduled_trigger" "projecttrigger_kubernetes_web_app_daily_security_scan" {
  count       = "${length(data.octopusdeploy_projects.projecttrigger_kubernetes_web_app_daily_security_scan.projects) != 0 ? 0 : 1}"
  space_id    = "${trimspace(var.octopus_space_id)}"
  name        = "Daily Security Scan"
  description = "This trigger reruns the deployment in the Security environment. This means any new vulnerabilities detected after a production deployment will be identified."
  timezone    = "UTC"
  is_disabled = false
  project_id  = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) != 0 ? data.octopusdeploy_projects.project_kubernetes_web_app.projects[0].id : octopusdeploy_project.project_kubernetes_web_app[0].id}"
  tenant_ids  = []

  once_daily_schedule {
    start_time   = "2025-04-26T09:00:00"
    days_of_week = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]
  }

  deploy_latest_release_action {
    source_environment_id      = "${length(data.octopusdeploy_environments.environment_security.environments) != 0 ? data.octopusdeploy_environments.environment_security.environments[0].id : octopusdeploy_environment.environment_security[0].id}"
    destination_environment_id = "${length(data.octopusdeploy_environments.environment_security.environments) != 0 ? data.octopusdeploy_environments.environment_security.environments[0].id : octopusdeploy_environment.environment_security[0].id}"
    should_redeploy            = true
  }
  lifecycle {
    prevent_destroy = true
  }
}

variable "project_kubernetes_web_app_name" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The name of the project exported from Kubernetes Web App"
  default     = "Kubernetes Web App 2"
}
variable "project_kubernetes_web_app_description_prefix" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "An optional prefix to add to the project description for the project Kubernetes Web App"
  default     = ""
}
variable "project_kubernetes_web_app_description_suffix" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "An optional suffix to add to the project description for the project Kubernetes Web App"
  default     = ""
}
variable "project_kubernetes_web_app_description" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The description of the project exported from Kubernetes Web App"
  default     = "This project provides an example Kubernetes deployment using YAML, Cloud Target Discovery, and SBOM scanning to an AWS EKS Kubernetes cluster."
}
variable "project_kubernetes_web_app_tenanted" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The tenanted setting for the project Untenanted"
  default     = "Untenanted"
}
data "octopusdeploy_projects" "project_kubernetes_web_app" {
  ids          = null
  partial_name = "${var.project_kubernetes_web_app_name}"
  skip         = 0
  take         = 1
}
resource "octopusdeploy_project" "project_kubernetes_web_app" {
  count                                = "${length(data.octopusdeploy_projects.project_kubernetes_web_app.projects) != 0 ? 0 : 1}"
  name                                 = "${var.project_kubernetes_web_app_name}"
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  discrete_channel_release             = false
  is_disabled                          = false
  is_version_controlled                = false
  lifecycle_id                         = "${length(data.octopusdeploy_lifecycles.lifecycle_security.lifecycles) != 0 ? data.octopusdeploy_lifecycles.lifecycle_security.lifecycles[0].id : octopusdeploy_lifecycle.lifecycle_security[0].id}"
  project_group_id                     = "${data.octopusdeploy_project_groups.project_group_default_project_group.project_groups[0].id}"
  included_library_variable_sets       = []
  tenanted_deployment_participation    = "${var.project_kubernetes_web_app_tenanted}"

  connectivity_policy {
    allow_deployments_to_no_targets = true
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "None"
  }

  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.NextPatch}"
  }
  description = "${var.project_kubernetes_web_app_description_prefix}${var.project_kubernetes_web_app_description}${var.project_kubernetes_web_app_description_suffix}"
  lifecycle {
    prevent_destroy = true
  }
}


