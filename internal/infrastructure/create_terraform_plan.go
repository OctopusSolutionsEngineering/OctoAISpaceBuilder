package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/validation"
)

// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/data/aztables
func CreateFeedbackAzureStorageTable(plan model.TerraformPlan) error {
	if err := validation.ValidateTerraformPlan(plan); err != nil {
		return err
	}

	service, err := aztables.NewServiceClientFromConnectionString(GetStorageConnectionString(), nil)

	if err != nil {
		return err
	}

	ctx := context.Background()

	if err := CreateTable(service, ctx); err != nil {
		return err
	}

	client := service.NewClient("TerraformPlan")

	myEntity := aztables.EDMEntity{
		Entity: aztables.Entity{
			PartitionKey: plan.Server,
			RowKey:       plan.ID,
		},
		Properties: map[string]any{
			"PlanBinary": plan.PlanBinaryBase64,
		},
	}
	marshalled, err := json.Marshal(myEntity)

	if _, err := client.AddEntity(ctx, marshalled, nil); err != nil {
		return err
	}

	return nil
}
