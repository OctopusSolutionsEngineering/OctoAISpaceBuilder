# https://developer.hashicorp.com/terraform/cloud-docs/policy-enforcement/define-policies/opa
# /home/matthew/Code/OctoAISpaceBuilder/binaries/opa_linux_amd64 exec --fail --decision "terraform/analysis/allow" --bundle /home/matthew/Code/OctoAISpaceBuilder/policy /tmp/tempdir1112410754/plan.json; echo $?

package terraform.analysis

# The default is to not pass validation
default allow := false

# Don't allow any changes to non-Octopus Deploy resources
affects_non_octopusdeploy_resources[msg] if {
    resource_change := input.resource_changes[_]
    not startswith(resource_change.type, "octopusdeploy_")
    # Generate a failure message
    msg := sprintf("Attempted to create type of ': %v", [resource_change.type])
}

# Make sure all sensitive values defulat to "Change Me!"
custom_sensitive_vars[msg] if {
    # Get resources from planned_values
    resource := input.planned_values.root_module.resources[_]

    # Check if sensitive_values exists in the resource
    is_object(resource.sensitive_values)

    # Find all true values in sensitive_values which mark sensitive fields
    [path, value] := walk(resource.sensitive_values)
    value == true

    # Get the corresponding value from the actual resource values
    is_object(resource.values)

    # Find the corrosponding properties under values
    [actual_path, actual_value] := walk(resource.values)
    actual_path == path

    # Find those properties that are not set to a default value
    actual_value != "Change Me!"
    actual_value != "CHANGE_ME"
    actual_value != "CHANGE ME"
    actual_value != "AWS_SECRET_KEY"

    # Feed usernames have been flagged as sensitive, so we include a couple of exceptions here
    actual_value != "x-access-token"
    actual_value != "username"
    actual_value != "mcasperson"

    # Variable references are ok
    not regex.match(`^#\{[^}]+\}$`, actual_value)

    actual_value != null

    # Generate a failure message
    msg := sprintf("Resource %s has a sensitive value at path %v that is not a generic placholder: %v",
                  [resource.address, concat(".", path), actual_value])
}

multiple_providers[msg] if {
    # provider_config must only have 'octopusdeploy' as key
    keys := object.keys(input.configuration.provider_config)
    count(keys) != 1
    msg := "provider_config must only contain the 'octopusdeploy' property"
}

not_octopus_provider[msg] if {
    keys := object.keys(input.configuration.provider_config)
    keys[0] != "octopusdeploy"
    msg := "provider_config must only contain the 'octopusdeploy' property"
}

provider_full_name_not_octopus[msg] if {
    # full_name must match exactly
    input.configuration.provider_config.octopusdeploy.full_name != "registry.opentofu.org/octopusdeploy/octopusdeploy"
    msg := "provider_config.octopusdeploy.full_name must be 'registry.opentofu.org/octopusdeploy/octopusdeploy'"
}

# This is the combined rule we want to check
allow if {
	count(affects_non_octopusdeploy_resources) == 0
	count(custom_sensitive_vars) == 0
	count(multiple_providers) == 0
	count(not_octopus_provider) == 0
	count(provider_full_name_not_octopus) == 0
}