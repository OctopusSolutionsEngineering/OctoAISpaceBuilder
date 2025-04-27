package environment

import "os"

func GetOpaExecutable() string {
	if os.Getenv("SPACEBUILDER_OPA_PATH") != "" {
		return os.Getenv("SPACEBUILDER_OPA_PATH")
	}

	return "binaries/opa_linux_amd64"
}

func GetOpaPolicyPath() string {
	if os.Getenv("SPACEBUILDER_OPA_POLICY_PATH") != "" {
		return os.Getenv("SPACEBUILDER_OPA_POLICY_PATH")
	}

	return "policy/"
}
