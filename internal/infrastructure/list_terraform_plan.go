package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"time"
)

func ListTerraformPlan(olderThanMinutes int) ([]TableEntityId, error) {
	service, err := aztables.NewServiceClientFromConnectionString(GetStorageConnectionString(), nil)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	if err := CreateTable(service, ctx); err != nil {
		return nil, err
	}

	client := service.NewClient("TerraformPlan")

	// Calculate the cutoff time
	cutoffTime := time.Now().Add(-time.Duration(olderThanMinutes) * time.Minute)

	// Query for old records based on timestamp
	// Azure Table Storage uses Timestamp field by default
	filter := fmt.Sprintf("Timestamp lt datetime'%s'", cutoffTime.UTC().Format("2006-01-02T15:04:05Z"))

	pager := client.NewListEntitiesPager(&aztables.ListEntitiesOptions{
		Filter: &filter,
	})

	results := []TableEntityId{}
	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		// Process each entity in the page
		for _, entity := range response.Entities {
			var tableEntity aztables.EDMEntity
			err = json.Unmarshal(entity, &tableEntity)

			results = append(results, TableEntityId{RowKey: tableEntity.RowKey, PartitionKey: tableEntity.PartitionKey})

		}
	}

	return results, nil
}
