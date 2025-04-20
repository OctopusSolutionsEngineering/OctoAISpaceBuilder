This microservice implements a web app, designed to be deployed as an Azure function, that generates a Terraform plan for an Octopus space with a call to `/api/terraformplan`, stores the plan, and then applies it with a call to `/api/terraformapply`. 

This services forces the use of a local backend, which is not persisted between calls. The Terraform configuration is expected to be stateless. See [Octoterra](https://github.com/OctopusSolutionsEngineering/OctopusTerraformExport) for more details on stateless Terraform configurations.

## Security Considerations

* The use of local state is enforced with overrides.
* Plugins must use a local mirror with the CLI configuration, so no providers are downloaded.
* OPA enforces rules that only allow Octopus resources to be created.
* OPA enforces rules that prohibit sensitive values being defined in the configuration.
* The client can only send Terraform configuration. The service defines all the file names (i.e. the service controls overrides).
* The client can not send variable values as separate files.
* Old plan records are deleted after 5 minutes.