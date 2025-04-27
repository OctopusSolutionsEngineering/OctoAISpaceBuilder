package validation

import (
	"testing"

	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateTerraformPlanRequest(t *testing.T) {
	tests := []struct {
		name    string
		plan    model.TerraformPlan
		wantErr bool
	}{
		{
			name:    "valid plan with space_id",
			plan:    model.TerraformPlan{SpaceId: "Spaces-123"},
			wantErr: false,
		},
		{
			name:    "empty space_id",
			plan:    model.TerraformPlan{SpaceId: ""},
			wantErr: true,
		},
		{
			name:    "whitespace space_id",
			plan:    model.TerraformPlan{SpaceId: "   "},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTerraformPlanRequest(tt.plan)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "space_id is required")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
