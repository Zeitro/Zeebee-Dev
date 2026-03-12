package auth

import (
	"fmt"

	"zeebee-backend/domains/zeebee/constants"
)

func BuildPolicy(roleId int64, resourceId int64, resourceType constants.UserType) []string {
	return []string{
		fmt.Sprintf("Role:%d", roleId),
		string(resourceType),
		fmt.Sprintf("Resource:%d", resourceId),
		"allow",
	}
}

func BuildInterfacePolicy(roleId int64, resourceId int64, resourceType constants.UserType) []interface{} {
	return []interface{}{
		fmt.Sprintf("Role:%d", roleId),
		string(resourceType),
		fmt.Sprintf("Resource:%d", resourceId),
		"allow",
	}
}
