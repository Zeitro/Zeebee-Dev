package auth

import (
	"fmt"

	"zeebee-backend/domains/zeebee/constants"
)

func BuildPolicy(roleId int64, resourceId int64, resourceType constants.UserType) []string {
	return []string{
		fmt.Sprintf(roleFormat, roleId),
		string(resourceType),
		fmt.Sprintf(resourceFormat, resourceId),
		"allow",
	}
}

func BuildInterfacePolicy(roleId int64, resourceId int64, resourceType constants.UserType) []interface{} {
	return []interface{}{
		fmt.Sprintf(roleFormat, roleId),
		string(resourceType),
		fmt.Sprintf(resourceFormat, resourceId),
		"allow",
	}
}
