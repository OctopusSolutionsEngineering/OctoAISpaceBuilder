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

	if plan.PlanBinaryBase64 == nil || strings.TrimSpace(*plan.PlanBinaryBase64) == "" {
		return errors.New("plan_binary is required")
	}

	return nil
}
