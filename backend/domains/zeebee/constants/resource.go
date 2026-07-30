package constants

import "zeebee-backend/common/enums"

type ResourceType string

const ResourceTypeName = "resource_type"

const (
	ResourceTypeMenu   ResourceType = "Menu"
	ResourceTypeButton ResourceType = "Button"
)

type ResourceTreeScope string

const ResourceTreeScopeName = "resource_tree_scope"
const (
	ResourceTreeScopeAll          ResourceTreeScope = "All"
	ResourceTreeScopeOrganization ResourceTreeScope = "Organization"
)

func init() {
	// register resource type
	enums.Register(ResourceTypeName, []enums.EnumOpt{
		{
			Label: string(ResourceTypeMenu),
			Value: string(ResourceTypeMenu),
		},
		{
			Label: string(ResourceTypeButton),
			Value: string(ResourceTypeButton),
		},
	})

	// register resource tree scope
	enums.Register(ResourceTreeScopeName, []enums.EnumOpt{
		{
			Label: string(ResourceTreeScopeAll),
			Value: string(ResourceTreeScopeAll),
		},
		{
			Label: string(ResourceTreeScopeOrganization),
			Value: string(ResourceTreeScopeOrganization),
		},
	})
}
