package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/aztables
func CreatePlanAzureStorageTable(id string, planBinary string, spaceId string, server string, lockFile string, configuration string) error {
	service, err := aztables.NewServiceClientFromConnectionString(GetStorageConnectionString(), nil)

	if err != nil {
		return err
	}

	ctx := context.Background()

	if err := CreateTable(service, ctx); err != nil {
		return err
	}

	client := service.NewClient("TerraformPlan")

	// We need to save enough of the terraform configuration to be able to recreate the plan
	// on a new function instance.
	// https://developer.hashicorp.com/terraform/tutorials/automation/automate-terraform#plan-and-apply-on-different-machines
	myEntity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: server,
			RowKey:       id,
		},
		Properties: map[string]any{
			"PlanBinary":    planBinary,
			"SpaceId":       spaceId,
			"LockFile":      lockFile,
			"Configuration": configuration,
		},
	}
	marshalled, err := json.Marshal(myEntity)

	if _, err := client.AddEntity(ctx, marshalled, nil); err != nil {
		return err
	}

	return nil
}
