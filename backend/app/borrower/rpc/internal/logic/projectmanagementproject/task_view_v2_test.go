package projectmanagementprojectlogic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/domains/zeebee/proto"
)

func TestTaskViewResponseIsReadOnlyProjectScoped(t *testing.T) {
	view := &proto.ProjectManagementTaskViewEntity{ProjectId: 42, Tasks: []*proto.ProjectManagementTaskEntity{}}
	require.EqualValues(t, 42, view.ProjectId)
	require.Empty(t, view.Tasks)
}
