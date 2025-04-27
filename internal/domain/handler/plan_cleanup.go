package handler

import (
	"errors"
	"fmt"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/infrastructure"
	"go.uber.org/zap"
)

const OldPlanThresholdMinutes = 5

func RemoveOldPlans() error {
	oldPlans, err := infrastructure.ListPlanAzureStorageBlob(OldPlanThresholdMinutes)

	if err != nil {
		return err
	}

	zap.L().Info("Cleaning up " + fmt.Sprint(len(oldPlans)) + " old plans")

	// Try to delete all the old plans, collecting any errors rather than returning immediately
	deleteErrors := []error{}
	for _, oldPlan := range oldPlans {
		if err := infrastructure.DeletePlanAzureStorageBlob(oldPlan); err != nil {
			deleteErrors = append(deleteErrors, err)
		}
	}

	if len(deleteErrors) > 0 {
		return errors.Join(deleteErrors...)
	}

	return nil
}
