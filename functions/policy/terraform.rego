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
    actual_value != null

    # Generate a failure message
    msg := sprintf("Resource %s has a sensitive value at path %v that is not 'Change Me!': %v",
                  [resource.address, concat(".", path), actual_value])
}

# This is the combined rule we want to check
allow if {
	count(affects_non_octopusdeploy_resources) == 0
	count(custom_sensitive_vars) == 0
}