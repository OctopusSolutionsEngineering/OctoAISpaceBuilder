package infrastructure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"strings"
)

func CreatePlanAzureStorageBlob(id string, planBinary []byte, lockFile []byte, configuration []byte) error {

	ctx := context.Background()

	connectionString := GetStorageConnectionString()
	if connectionString == "" {
		return fmt.Errorf("AzureWebJobsStorage environment variable is not set")
	}

	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return err
	}

	if _, err := client.CreateContainer(ctx, "terraformplan", nil); err != nil {
		if !strings.Contains(err.Error(), "ContainerAlreadyExists") {
			return err
		}
	}

	// Upload the binary plan
	if _, err := client.UploadBuffer(ctx, "terraformplan", "plan."+id+".binary", planBinary, nil); err != nil {
		return err
	}

	// Upload the lock file
	if _, err := client.UploadBuffer(ctx, "terraformplan", "plan."+id+".lock", lockFile, nil); err != nil {
		return err
	}

	// Upload the lock file
	if _, err := client.UploadBuffer(ctx, "terraformplan", "plan."+id+".configuration", configuration, nil); err != nil {
		return err
	}

	return nil
}
