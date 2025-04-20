package application

import (
	"github.com/OctopusSolutionsEngineering/OctoAISpaceBuilder/internal/domain/handler"
	"github.com/gin-gonic/gin"
)

func CleanupOldPlans(c *gin.Context) {
	err := handler.RemoveOldPlans()

	if err != nil {
		c.IndentedJSON(500, gin.H{"error": "Failed to remove old plans"})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Old plans removed successfully"})
}
