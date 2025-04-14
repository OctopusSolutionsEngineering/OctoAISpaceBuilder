package validation

import (
	"errors"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"strings"
)

func ValidateTerraformPlan(plan model.TerraformPlan) error {
	if strings.TrimSpace(plan.ID) == "" {
		return errors.New("id is required")
	}

	if strings.TrimSpace(plan.Plan) == "" {
		return errors.New("plan is required")
	}

	return nil
}
