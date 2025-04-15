This microservice implements a web app, designed to be deployed as an Azure function, that generates a Terraform plan for an Octopus space with a call to `/api/terraformplan`, stores the plan, and then applies it with a call to `/api/terraformapply`. 

The plans are validated using Open Policy Agent to restrict what they can create. 

The Terraform plugins are pre-downloaded and downloading random plugins is disabled.

This services forces the use of a local backend, which is not persisted between calls. The Terraform configuration is expected to be stateless. See [Octoterra](https://github.com/OctopusSolutionsEngineering/OctopusTerraformExport) for more details on stateless Terraform configurations.