package application

import (
	"github.com/DataDog/jsonapi"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/application/responses"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func Health(c *gin.Context) {
	planMetricsSnapshot := terraformPlanMetrics.Snapshot(time.Now())
	metricsSnapshot := terraformApplyMetrics.Snapshot(time.Now())

	responseJSON, err := jsonapi.Marshal(model.Health{
		ID:                                uuid.New().String(),
		Status:                            "OK",
		TerraformPlanSuccessLast1Minute:   planMetricsSnapshot.SuccessLast1Minute,
		TerraformPlanFailedLast1Minute:    planMetricsSnapshot.FailedLast1Minute,
		TerraformPlanSuccessLast15Minute:  planMetricsSnapshot.SuccessLast15Minutes,
		TerraformPlanFailedLast15Minute:   planMetricsSnapshot.FailedLast15Minutes,
		TerraformPlanSuccessLast1Hour:     planMetricsSnapshot.SuccessLast1Hour,
		TerraformPlanFailedLast1Hour:      planMetricsSnapshot.FailedLast1Hour,
		TerraformPlanSuccessAllTime:       planMetricsSnapshot.SuccessAllTime,
		TerraformPlanFailedAllTime:        planMetricsSnapshot.FailedAllTime,
		TerraformApplySuccessLast1Minute:  metricsSnapshot.SuccessLast1Minute,
		TerraformApplyFailedLast1Minute:   metricsSnapshot.FailedLast1Minute,
		TerraformApplySuccessLast15Minute: metricsSnapshot.SuccessLast15Minutes,
		TerraformApplyFailedLast15Minute:  metricsSnapshot.FailedLast15Minutes,
		TerraformApplySuccessLast1Hour:    metricsSnapshot.SuccessLast1Hour,
		TerraformApplyFailedLast1Hour:     metricsSnapshot.FailedLast1Hour,
		TerraformApplySuccessAllTime:      metricsSnapshot.SuccessAllTime,
		TerraformApplyFailedAllTime:       metricsSnapshot.FailedAllTime,
	})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.GenerateError("Failed to process request", err))
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(200, string(responseJSON))
}
