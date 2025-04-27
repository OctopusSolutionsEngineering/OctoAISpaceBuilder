package application

import (
	"bytes"
	"github.com/DataDog/jsonapi"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/files"
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestTerraformPlanAndApplyEndpoint plans and then applies a terraform configuration via the application
// level Gin interface.
func TestTerraformPlanAndApplyEndpoint(t *testing.T) {
	// Use the system installation of OPA
	if err := os.Setenv("SPACEBUILDER_OPA_PATH", "opa"); err != nil {
		t.Fatalf("Failed to set SPACEBUILDER_OPA_PATH: %v", err)
	}

	// Use the system installation of Tofu
	if err := os.Setenv("SPACEBUILDER_TOFU_PATH", "tofu"); err != nil {
		t.Fatalf("Failed to set SPACEBUILDER_TOFU_PATH: %v", err)
	}

	// Disable the Terraform CLI config to allow provider downloads
	if err := os.Setenv("SPACEBUILDER_DISABLE_TERRAFORM_CLI_CONFIG", "true"); err != nil {
		t.Fatalf("Failed to set SPACEBUILDER_DISABLE_TERRAFORM_CLI_CONFIG: %v", err)
	}

	// Set the OPA policy path to the local policy directory
	if err := os.Setenv("SPACEBUILDER_OPA_POLICY_PATH", "../../functions/policy/"); err != nil {
		t.Fatalf("Failed to set SPACEBUILDER_DISABLE_TERRAFORM_CLI_CONFIG: %v", err)
	}

	base, err := files.CopyDir("../../terraform")

	if err != nil {
		t.Fatalf("Failed to copy Terraform files: %v", err)
	}

	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		spaceId, err := testFramework.Act(t, container, base, "2-localsetup", []string{})

		if err != nil {
			t.Fatalf("Failed to create space: %v", err)
		}

		// Set up Gin in test mode
		gin.SetMode(gin.TestMode)

		// Create a test router with the endpoint
		router := gin.Default()

		// Register the endpoint with mocked middleware
		router.POST("/api/terraformplan", CreateTerraformPlan)
		router.POST("/api/terraformapply", CreateTerraformApply)

		reponse := func() model.TerraformPlan {
			configuration, err := os.ReadFile("../../terraform/k8s-example/example.tf")

			if err != nil {
				t.Fatalf("Failed to read configuration file: %v", err)
			}

			// Apply the changes
			body := model.TerraformPlan{
				ID:            "unused",
				SpaceId:       spaceId,
				Configuration: string(configuration),
			}

			jsonBody, err := jsonapi.Marshal(body)
			require.NoError(t, err)

			// Create test request
			req, err := http.NewRequest(http.MethodPost, "/api/terraformplan", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Octopus-Url", container.URI)
			req.Header.Set("X-Octopus-ApiKey", test.ApiKey)

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, 201, w.Code)

			var response model.TerraformPlan
			err = jsonapi.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.NotEmpty(t, response.ID)
			assert.Equal(t, spaceId, response.SpaceId)

			return response
		}()

		// Assert response
		func(response model.TerraformPlan) {
			// Plan the changes
			applyBody := model.TerraformApply{
				ID:        "unused",
				PlanId:    response.ID,
				Server:    "",
				ApplyText: nil,
			}

			applyJsonBody, err := jsonapi.Marshal(applyBody)
			require.NoError(t, err)

			// Create test request
			applyReq, err := http.NewRequest(http.MethodPost, "/api/terraformapply", bytes.NewBuffer(applyJsonBody))
			require.NoError(t, err)
			applyReq.Header.Set("Content-Type", "application/json")
			applyReq.Header.Set("X-Octopus-Url", container.URI)
			applyReq.Header.Set("X-Octopus-ApiKey", test.ApiKey)

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, applyReq)

			// Assert status code
			assert.Equal(t, 201, w.Code)

			// Assert response
			var applyResponse model.TerraformApply
			err = jsonapi.Unmarshal(w.Body.Bytes(), &applyResponse)
			assert.NoError(t, err)
			assert.NotEmpty(t, applyResponse.ID)
		}(reponse)

		return nil
	})
}
