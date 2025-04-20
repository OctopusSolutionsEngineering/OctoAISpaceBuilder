package application

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CleanupOldPlans(c *gin.Context) {
	err := handler.RemoveOldPlans()

	if err != nil {
		zap.L().Error("Error removing old plans", zap.Error(err))
		c.IndentedJSON(500, gin.H{"error": "Failed to remove old plans"})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Old plans removed successfully"})
}
