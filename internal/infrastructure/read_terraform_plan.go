package infrastructure

import (
	"encoding/json"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"golang.org/x/net/context"
)

func ReadFeedbackAzureStorageTable(terraformApply model.TerraformApply) (string, string, string, string, error) {
	service, err := aztables.NewServiceClientFromConnectionString(GetStorageConnectionString(), nil)

	if err != nil {
		return "", "", "", "", err
	}

	ctx := context.Background()

	if err := CreateTable(service, ctx); err != nil {
		return "", "", "", "", err
	}

	client := service.NewClient("TerraformPlan")

	resp, err := client.GetEntity(ctx, terraformApply.Server, terraformApply.PlanId, nil)

	if err != nil {
		return "", "", "", "", err
	}

	var entity aztables.EDMEntity
	err = json.Unmarshal(resp.Value, &entity)
	if err != nil {
		return "", "", "", "", err
	}

	planBinary := ""
	if value, ok := entity.Properties["PlanBinary"]; ok {
		if value, ok := value.(string); ok {
			planBinary = value
		}
	}

	spaceId := ""
	if value, ok := entity.Properties["SpaceId"]; ok {
		if value, ok := value.(string); ok {
			spaceId = value
		}
	}

	lockFile := ""
	if value, ok := entity.Properties["LockFile"]; ok {
		if value, ok := value.(string); ok {
			lockFile = value
		}
	}

	configuration := ""
	if value, ok := entity.Properties["Configuration"]; ok {
		if value, ok := value.(string); ok {
			configuration = value
		}
	}

	if planBinary != "" && spaceId != "" {
		return planBinary, spaceId, lockFile, configuration, nil
	}

	return "", "", "", "", errors.New("could not find PlanJson, SpaceId, LockFile, or Configuration in entity properties")
}
