package projectmanagementprojectmemberlogic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/domains/zeebee/data/dto"
)

func TestMemberCandidateEntity(t *testing.T) {
	entity, err := memberCandidateEntity(&dto.ProjectManagementMemberCandidateEntityDto{
		OrganizationUserId: 95,
		FullName:           "Awei Chen",
		Email:              "544278032@qq.com",
	})
	require.NoError(t, err)
	require.Equal(t, int64(95), entity.OrganizationUserId)
	require.Equal(t, "Awei Chen", entity.FullName)
	require.Equal(t, "544278032@qq.com", entity.Email)
}
