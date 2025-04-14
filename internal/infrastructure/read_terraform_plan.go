package infrastructure

import (
	"encoding/json"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"golang.org/x/net/context"
)

func ReadFeedbackAzureStorageTable(terraformApply model.TerraformApply) (string, error) {
	service, err := aztables.NewServiceClientFromConnectionString(GetStorageConnectionString(), nil)

	if err != nil {
		return "", err
	}

	ctx := context.Background()

	if err := CreateTable(service, ctx); err != nil {
		return "", err
	}

	client := service.NewClient("TerraformPlan")

	resp, err := client.GetEntity(ctx, terraformApply.Server, terraformApply.PlanId, nil)

	if err != nil {
		return "", err
	}

	var entity aztables.EDMEntity
	err = json.Unmarshal(resp.Value, &entity)
	if err != nil {
		return "", err
	}

	if value, ok := entity.Properties["PlanBinary"]; ok {
		if value, ok := value.(string); ok {
			return value, nil
		}
	}

	return "", errors.New("could not find PlanJson in entity properties")
}
