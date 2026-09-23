package projectmanagementtemplatelogic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/admin/rpc/pb/admin"
)

func TestPublicTemplateStatusContractIsClosed(t *testing.T) {
	require.Equal(t, "Active", (&admin.ProjectManagementTemplateStatusRequestData{Status: "Active"}).Status)
	require.Equal(t, "Disabled", (&admin.ProjectManagementTemplateStatusRequestData{Status: "Disabled"}).Status)
	for _, status := range []string{"Draft", "Archived"} {
		require.NotEqual(t, "Active", status)
		require.NotEqual(t, "Disabled", status)
	}
}
