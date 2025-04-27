package environment

import "os"

func GetTofuExecutable() string {
	if os.Getenv("SPACEBUILDER_TOFU_PATH") != "" {
		return os.Getenv("SPACEBUILDER_TOFU_PATH")
	}

	return "binaries/tofu"
}
