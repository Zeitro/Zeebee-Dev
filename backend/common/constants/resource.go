package constants

import "zeebee-backend/common/enums"

type ResourceType string

const ResourceTypeName = "resource_type"

const (
	ResourceTypeMenu ResourceType = "Menu"
	ResourceTypeApi  ResourceType = "Api"
)

func init() {
	// register resource type
	enums.Register(ResourceTypeName, []enums.EnumOpt{
		{
			Label: string(ResourceTypeMenu),
			Value: string(ResourceTypeMenu),
		},
		{
			Label: string(ResourceTypeApi),
			Value: string(ResourceTypeApi),
		},
	})
}
