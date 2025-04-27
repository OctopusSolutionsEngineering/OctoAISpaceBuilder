package infrastructure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.uber.org/zap"
	"io"
	"strings"
)

func ReadPlanAzureStorageBlob(id string) ([]byte, []byte, []byte, error) {

	ctx := context.Background()

	connectionString := GetStorageConnectionString()
	if connectionString == "" {
		return nil, nil, nil, fmt.Errorf("AzureWebJobsStorage environment variable is not set")
	}

	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	if _, err := client.CreateContainer(ctx, "terraformplan", nil); err != nil {
		if !strings.Contains(err.Error(), "ContainerAlreadyExists") {
			return nil, nil, nil, err
		}
	}

	// Download the binary plan
	downloadResponse, err := client.DownloadStream(ctx, "terraformplan", "plan."+id+".binary", nil)
	if err != nil {
		return nil, nil, nil, err
	}
	planBinary, err := io.ReadAll(downloadResponse.Body)
	if err != nil {
		return nil, nil, nil, err
	}

	defer func() {
		if err := downloadResponse.Body.Close(); err != nil {
			zap.L().Error("failed to close download response body", zap.Error(err))
		}
	}()

	// Download the lock file
	lockResponse, err := client.DownloadStream(ctx, "terraformplan", "plan."+id+".lock", nil)
	if err != nil {
		return nil, nil, nil, err
	}
	planLock, err := io.ReadAll(lockResponse.Body)
	if err != nil {
		return nil, nil, nil, err
	}
	defer func() {
		if err := lockResponse.Body.Close(); err != nil {
			zap.L().Error("failed to close download response body", zap.Error(err))
		}
	}()

	// Download the configuration file
	configResponse, err := client.DownloadStream(ctx, "terraformplan", "plan."+id+".configuration", nil)
	if err != nil {
		return nil, nil, nil, err
	}
	planConfig, err := io.ReadAll(configResponse.Body)
	if err != nil {
		return nil, nil, nil, err
	}
	defer func() {
		if err := configResponse.Body.Close(); err != nil {
			zap.L().Error("failed to close download response body", zap.Error(err))
		}
	}()

	return planBinary, planLock, planConfig, nil
}
