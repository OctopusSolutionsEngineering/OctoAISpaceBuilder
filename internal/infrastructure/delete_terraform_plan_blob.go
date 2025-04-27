package infrastructure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"strings"
)

func DeletePlanAzureStorageBlob(id string) error {

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

	if _, err := client.DeleteBlob(ctx, "terraformplan", "plan."+id+".binary", nil); err != nil {
		if !strings.Contains(err.Error(), "BlobNotFound") {
			return err
		}
	}

	if _, err := client.DeleteBlob(ctx, "terraformplan", "plan."+id+".lock", nil); err != nil {
		if !strings.Contains(err.Error(), "BlobNotFound") {
			return err
		}
	}

	if _, err := client.DeleteBlob(ctx, "terraformplan", "plan."+id+".configuration", nil); err != nil {
		if !strings.Contains(err.Error(), "BlobNotFound") {
			return err
		}
	}

	return nil
}
