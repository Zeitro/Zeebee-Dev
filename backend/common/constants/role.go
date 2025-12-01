package constants

import "zeebee-backend/common/enums"

const AdminRoleName = "admin"

type DataScope string

var DataScopeName = "data_scope"

const (
	DataScopeAll DataScope = "All"
	DataScopeOwn DataScope = "Own"
)

func init() {
	// register data scope
	enums.Register(DataScopeName, []enums.EnumOpt{
		{
			Label: string(DataScopeAll),
			Value: string(DataScopeAll),
		},
		{
			Label: string(DataScopeOwn),
			Value: string(DataScopeOwn),
		},
	})
}
