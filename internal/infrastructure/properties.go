package infrastructure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"time"
)

func GetTimeProperty(name string, entity aztables.EDMEntity, defaultValue time.Time) time.Time {
	if value, ok := entity.Properties[name]; ok {
		if strValue, ok := value.(string); ok {
			parsedTime, err := time.Parse(time.RFC3339, strValue)
			if err != nil {
				return defaultValue
			}

			return parsedTime
		}
	}

	return defaultValue
}

func GetStringProperty(name string, entity aztables.EDMEntity, defaultValue string) string {
	if value, ok := entity.Properties[name]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}

	return defaultValue
}

func GetBoolProperty(name string, entity aztables.EDMEntity, defaultValue bool) bool {
	if value, ok := entity.Properties[name]; ok {
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}
	}

	return defaultValue
}
