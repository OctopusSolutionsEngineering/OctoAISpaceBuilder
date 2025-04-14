# https://developer.hashicorp.com/terraform/cloud-docs/policy-enforcement/define-policies/opa
# /home/matthew/Code/OctoAISpaceBuilder/binaries/opa_linux_amd64 exec --fail --decision "terraform/analysis/allow" --bundle /home/matthew/Code/OctoAISpaceBuilder/policy /tmp/tempdir1112410754/plan.json; echo $?

package terraform.analysis

import input as tfplan

# The default is to not pass validation
default allow := false

# Don't allow any changes to non-Octopus Deploy resources
affects_non_octopusdeploy_resources if {
    some resource_change in tfplan.resource_changes
    not startswith(resource_change.type, "octopusdeploy_")
}

# This is the combined rule we want to check
allow if {
	not affects_non_octopusdeploy_resources
}