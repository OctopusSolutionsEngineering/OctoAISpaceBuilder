package infrastructure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

func DeleteTerraformPlan(tableEntity TableEntityId) error {
	service, err := aztables.NewServiceClientFromConnectionString(GetStorageConnectionString(), nil)

	if err != nil {
		return err
	}

	ctx := context.Background()

	if err := CreateTable(service, ctx); err != nil {
		return err
	}

	client := service.NewClient("TerraformPlan")

	_, err = client.DeleteEntity(ctx, tableEntity.PartitionKey, tableEntity.RowKey, nil)

	return err
}
