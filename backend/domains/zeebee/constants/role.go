package constants

import "zeebee-backend/common/enums"

const AdminRoleName = "Admin"

type DataScope string

var DataScopeName = "data_scope"

const (
	DataScopeAll DataScope = "All"
	DataScopeOwn DataScope = "Own"
)

// Priority
func (d DataScope) Priority() int {
	switch d {
	case DataScopeAll:
		return 1
	case DataScopeOwn:
		return 2
	default:
		return 999
	}
}

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
