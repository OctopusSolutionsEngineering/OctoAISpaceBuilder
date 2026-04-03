package application

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DataDog/jsonapi"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTerraformApplyCallMetricsSnapshotWindows(t *testing.T) {
	metrics := terraformApplyCallMetrics{}
	currentMinute := time.Now().Unix() / 60

	metrics.buckets[currentMinute%60] = metricsBucket{minuteEpoch: currentMinute, success: 2, failed: 1}
	metrics.buckets[(currentMinute-10)%60] = metricsBucket{minuteEpoch: currentMinute - 10, success: 3, failed: 4}
	metrics.buckets[(currentMinute-30)%60] = metricsBucket{minuteEpoch: currentMinute - 30, success: 5, failed: 6}
	metrics.buckets[(currentMinute-59)%60] = metricsBucket{minuteEpoch: currentMinute - 59, success: 7, failed: 8}

	metrics.allTimeSuccess = 50
	metrics.allTimeFailed = 60

	snapshot := metrics.Snapshot(time.Unix(currentMinute*60, 0))

	assert.Equal(t, int64(2), snapshot.SuccessLast1Minute)
	assert.Equal(t, int64(1), snapshot.FailedLast1Minute)
	assert.Equal(t, int64(5), snapshot.SuccessLast15Minutes)
	assert.Equal(t, int64(5), snapshot.FailedLast15Minutes)
	assert.Equal(t, int64(17), snapshot.SuccessLast1Hour)
	assert.Equal(t, int64(19), snapshot.FailedLast1Hour)
	assert.Equal(t, int64(50), snapshot.SuccessAllTime)
	assert.Equal(t, int64(60), snapshot.FailedAllTime)
}

func TestCreateTerraformApplyRecordsFailedCall(t *testing.T) {
	terraformApplyMetrics = &terraformApplyCallMetrics{}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/terraformapply", CreateTerraformApply)

	req, err := http.NewRequest(http.MethodPost, "/api/terraformapply", bytes.NewBufferString("this is not jsonapi"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	snapshot := terraformApplyMetrics.Snapshot(time.Now())
	assert.Equal(t, int64(0), snapshot.SuccessAllTime)
	assert.Equal(t, int64(1), snapshot.FailedAllTime)
}

func TestCreateTerraformPlanRecordsFailedCall(t *testing.T) {
	terraformPlanMetrics = &terraformApplyCallMetrics{}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/terraformplan", CreateTerraformPlan)

	req, err := http.NewRequest(http.MethodPost, "/api/terraformplan", bytes.NewBufferString("this is not jsonapi"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	snapshot := terraformPlanMetrics.Snapshot(time.Now())
	assert.Equal(t, int64(0), snapshot.SuccessAllTime)
	assert.Equal(t, int64(1), snapshot.FailedAllTime)
}

func TestHealthEndpointIncludesTerraformApplyMetrics(t *testing.T) {
	terraformPlanMetrics = &terraformApplyCallMetrics{}
	terraformApplyMetrics = &terraformApplyCallMetrics{}
	terraformPlanMetrics.Record(true)
	terraformPlanMetrics.Record(false)
	terraformApplyMetrics.Record(true)
	terraformApplyMetrics.Record(false)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/api/health", Health)

	req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Health
	err = jsonapi.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, int64(1), response.TerraformPlanSuccessAllTime)
	assert.Equal(t, int64(1), response.TerraformPlanFailedAllTime)
	assert.Equal(t, int64(1), response.TerraformApplySuccessAllTime)
	assert.Equal(t, int64(1), response.TerraformApplyFailedAllTime)
}
