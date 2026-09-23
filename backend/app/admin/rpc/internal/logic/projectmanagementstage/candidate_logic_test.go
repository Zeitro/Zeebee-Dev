package projectmanagementstagelogic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCandidateWorkflowType(t *testing.T) {
	bound := map[int64]struct{}{1: {}}
	require.True(t, includeCandidateStage(bound, 2, workflowTypeUnbound))
	require.False(t, includeCandidateStage(bound, 1, workflowTypeUnbound))
	require.True(t, includeCandidateStage(bound, 1, workflowTypeBound))
	require.False(t, includeCandidateStage(bound, 2, workflowTypeBound))
}
