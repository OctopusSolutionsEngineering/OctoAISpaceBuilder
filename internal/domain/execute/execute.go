package execute

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Execute(executable string, args []string, env map[string]string) (string, string, int, error) {
	if err := makeExecutable(executable); err != nil {
		return "", "", 0, err
	}

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

	// Default exit code
	exitCode := 0

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Get the exit code from the wait status
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
		return stdout.String(), stderr.String(), exitCode, err
	}

	return stdout.String(), stderr.String(), exitCode, nil
}

func makeExecutable(executable string) error {
	info, err := os.Stat(executable)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", executable, err)
	}

	// Add execute permission to current permissions
	return os.Chmod(executable, info.Mode()|0111)
}
