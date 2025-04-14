package execute

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Execute(executable string, args []string, env map[string]string) (string, string, error) {
	cmd := exec.Command(executable, args...)

	// Set environment variables if provided
	if env != nil {
		for key, value := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("execution failed: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), stderr.String(), nil
}
