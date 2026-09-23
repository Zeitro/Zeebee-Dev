package common

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/domains/zeebee/data/dto"
)

func TestDropLegacyProjectWorkflowAPIs(t *testing.T) {
	resources := []dto.ResourceSyncItemDto{{
		Apis: []dto.ApiSyncItemDto{
			{Path: "/api/v1/project_management/project/:id/workflow"},
			{Path: "/api/v1/project_management/project/:id"},
		},
		Children: []dto.ResourceSyncItemDto{{
			Apis: []dto.ApiSyncItemDto{{Path: "/api/v1/project_management/workflow/:id/swap_stage_order"}},
		}},
	}}

	dropLegacyProjectWorkflowAPIs(resources)

	require.Len(t, resources[0].Apis, 1)
	require.Equal(t, "/api/v1/project_management/project/:id", resources[0].Apis[0].Path)
	require.Empty(t, resources[0].Children[0].Apis)
}
