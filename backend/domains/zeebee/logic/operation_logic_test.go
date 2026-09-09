package logic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

func TestValidateOperationPageFilter(t *testing.T) {
	require.NoError(t, validateOperationPageFilter(dto.OperationPageFilterDto{
		TaskId:     1,
		EntityType: constants.OperationEntityTypeProjectManagementProject,
	}))
	require.EqualError(t, validateOperationPageFilter(dto.OperationPageFilterDto{
		TaskId:     1,
		EntityType: constants.OperationEntityTypeBorrowerPropertyTransaction,
	}), "task id filter is only supported for project management project operations")
}
