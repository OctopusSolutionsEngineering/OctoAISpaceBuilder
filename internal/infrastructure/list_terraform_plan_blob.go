package infrastructure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/samber/lo"
	"strings"
	"time"
)

func ListPlanAzureStorageBlob(olderThanMinutes int) ([]string, error) {
	ctx := context.Background()

	connectionString := GetStorageConnectionString()
	if connectionString == "" {
		return nil, fmt.Errorf("AzureWebJobsStorage environment variable is not set")
	}

	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}

	if _, err := client.CreateContainer(ctx, "terraformplan", nil); err != nil {
		if !strings.Contains(err.Error(), "ContainerAlreadyExists") {
			return nil, err
		}
	}

	// Calculate the cutoff time
	cutoffTime := time.Now().Add(-time.Duration(olderThanMinutes) * time.Minute)

	pager := client.NewListBlobsFlatPager("terraformplan", nil)

	// Loop through the blobs
	ids := []string{}
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		for _, blob := range resp.Segment.BlobItems {
			if blob.Properties.CreationTime.Before(cutoffTime) {
				split := strings.Split(*blob.Name, ".")
				if len(split) == 3 && split[0] == "plan" && !lo.Contains(ids, split[1]) {
					ids = append(ids, split[1])
				}
			}
		}
	}

	return ids, nil
}
