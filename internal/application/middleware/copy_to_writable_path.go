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

	// These directories make up the installation directory
	dirs := []string{"binaries", "provider", "policy"}

	for _, dir := range dirs {

		destPath := filepath.Join(tempDir, dir)

		exists, err := pathExists(destPath)

		if err != nil {
			zap.L().Error("Failed to test the presence of the "+dir+" directory", zap.Error(err))
			c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to test the presence of the "+dir+" directory", err))
			c.Abort()
			return
		}

		if !exists {
			zap.L().Info("Copying "+dir+" to "+destPath, zap.Error(err))

			// Create the directory in temp
			if err := os.MkdirAll(destPath, 0755); err != nil {
				zap.L().Error("Failed to create the "+dir+" directory", zap.Error(err))
				c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to create the "+dir+" directory", err))
				c.Abort()
				return
			}

			// Recursively copy the directory
			if err := copyDir(dir, destPath); err != nil {
				zap.L().Error("Failed to copy the "+dir+" directory", zap.Error(err))
				c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to copy the "+dir+" directory", err))
				c.Abort()
				return
			}
		}
	}

	c.Next()
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
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
	// Skip existing files
	if exists, err := pathExists(dst); err != nil {
		return err
	} else if exists {
		return nil
	}

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
