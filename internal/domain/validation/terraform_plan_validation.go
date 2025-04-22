package validation

import (
	"errors"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"strings"
)

func ValidateTerraformPlanRequest(plan model.TerraformPlan) error {
	if strings.TrimSpace(plan.SpaceId) == "" {
		return errors.New("space_id is required")
	}

	return nil
}
