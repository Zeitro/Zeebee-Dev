package apilogic

import (
	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

func targetType(callerType constants.UserType, target string) (constants.UserType, error) {
	targetType, err := logic.ParseUserType(target)
	if err != nil {
		return "", err
	}
	if callerType != constants.UserTypeAdmin && callerType != targetType {
		return "", errorx.New("cannot manage another platform api")
	}
	return targetType, nil
}
