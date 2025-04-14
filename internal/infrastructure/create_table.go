package infrastructure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"golang.org/x/net/context"
	"strings"
)

func CreateTable(service *aztables.ServiceClient, ctx context.Context) error {
	if _, err := service.CreateTable(ctx, "TerraformPlan", nil); err != nil {
		if !strings.Contains(err.Error(), "TableAlreadyExists") {
			return err
		}
	}

	return nil
}
