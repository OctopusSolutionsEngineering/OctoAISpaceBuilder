package handler

import (
	"errors"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/infrastructure"
)

const OldPlanThresholdMinutes = 5

func RemoveOldPlans() error {
	oldPlans, err := infrastructure.ListTerraformPlan(OldPlanThresholdMinutes)

	if err != nil {
		return err
	}

	// Try to delete all the old plans, collecting any errors rather than returning immediately
	deleteErrors := []error{}
	for _, oldPlan := range oldPlans {
		if err := infrastructure.DeleteTerraformPlan(oldPlan); err != nil {
			deleteErrors = append(deleteErrors, err)
		}
	}

	if len(deleteErrors) > 0 {
		return errors.Join(deleteErrors...)
	}

	return nil
}
