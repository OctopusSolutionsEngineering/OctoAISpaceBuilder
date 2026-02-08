package execute

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"go.uber.org/zap"
)

func Execute(executable string, args []string, env map[string]string) (string, string, int, error) {
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
		zap.L().Error("Error executing command", zap.String("command", executable), zap.Error(err))
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

func MakeExecutable(executable string) error {
	info, err := os.Stat(executable)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", executable, err)
	}

	// Add execute permission to current permissions
	return os.Chmod(executable, info.Mode()|0111)
}

func MakeAllExecutable(directory string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			zap.L().Info("Scanning " + path + " for files to make executable")
			return nil
		}

		// Check if file has execute bit for any user
		if info.Mode()&0111 == 0 {
			if err := MakeExecutable(path); err != nil {
				return fmt.Errorf("failed to make %s executable: %w", path, err)
			}
			zap.L().Info("Made the file executable: " + path)
		} else {
			zap.L().Info("The file is already executable: " + path)
		}

		return nil
	})
}
