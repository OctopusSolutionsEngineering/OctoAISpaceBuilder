package middleware

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/environment"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CopyToWritablePath copies the binaries to a writeable, temporary location.
// This service relies on a number of binary files like OPA and terraform providers.
// When running in Azure functions, the file system is read-only.
// Copying files to the temp dir allows us to set the executable flag on the files.
func CopyToWritablePath(c *gin.Context) {
	if !environment.IsInAzureFunctions() {
		// Nothing to do when running locally
		c.Next()
		return
	}

	// Get the temporary directory
	tempDir := os.TempDir()
	binariesDestPath := filepath.Join(tempDir, "binaries")
	providerDestPath := filepath.Join(tempDir, "provider")

	binariesExists, err := dirExists(binariesDestPath)

	if err != nil {
		zap.L().Error("Failed to test the presence of the binaries directory", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to test the presence of the binaries directory", err))
		c.Abort()
		return
	}

	if !binariesExists {
		// Create the binaries directory in temp
		if err := os.MkdirAll(binariesDestPath, 0755); err != nil {
			zap.L().Error("Failed to create the binaries directory", zap.Error(err))
			c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to create the binaries directory", err))
			c.Abort()
			return
		}

		// Recursively copy the binaries directory
		if err := copyDir("binaries", binariesDestPath); err != nil {
			zap.L().Error("Failed to copy the binaries directory", zap.Error(err))
			c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to copy the binaries directory", err))
			c.Abort()
			return
		}
	}

	providersExists, err := dirExists(providerDestPath)

	if err != nil {
		zap.L().Error("Failed to test the presence of the providers directory", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to test the presence of the providers directory", err))
		c.Abort()
		return
	}

	if !providersExists {
		// Create the provider directory in temp
		if err := os.MkdirAll(providerDestPath, 0755); err != nil {
			zap.L().Error("Failed to create the provider directory", zap.Error(err))
			c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to create the provider directory", err))
			c.Abort()
			return
		}

		// Recursively copy the provider directory
		if err := copyDir("provider", providerDestPath); err != nil {
			zap.L().Error("Failed to copy the provider directory", zap.Error(err))
			c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to copy the provider directory", err))
			c.Abort()
			return
		}
	}

	c.Next()
}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		// Create directories
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy files
		return copyFile(path, dstPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			zap.L().Error("Failed to close the file ", zap.Error(err))
		}
	}()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			zap.L().Error("Failed to close the file ", zap.Error(err))
		}
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Preserve file permissions
	return os.Chmod(dst, mode)
}
