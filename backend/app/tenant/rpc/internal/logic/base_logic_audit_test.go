package logic_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/logic/testutil"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// RecordAudit is async; eventually-await the row then assert it carried the
// input fields through copierx into the Operation table.
func TestRecordAudit_WritesOperationRow(t *testing.T) {
	env := testutil.NewEnv(t)
	logic.ResetBaseLogicForTest()
	bl := logic.GetBaseLogic(env.SvcCtx)

	bl.RecordAudit(context.Background(), dto.OperationCreateDto{
		EntityType:   constants.OperationEntityTypeUser,
		EntityId:     42,
		OperatorId:   42,
		OperatorName: "alice@example.com",
		OperatorType: constants.UserTypeLender,
		Operation:    constants.OperationActionUserLogin,
	})

	require.Eventually(t, func() bool {
		var count int64
		env.DB.Model(&model.Operation{}).Count(&count)
		return count == 1
	}, time.Second, 5*time.Millisecond)

	var op model.Operation
	require.NoError(t, env.DB.First(&op).Error)
	assert.Equal(t, int64(42), op.OperatorId)
	assert.Equal(t, "alice@example.com", op.OperatorName)
	assert.Equal(t, constants.UserTypeLender, op.OperatorType)
	assert.Equal(t, constants.OperationActionUserLogin, op.Operation)
	assert.Equal(t, constants.OperationEntityTypeUser, op.EntityType)
}

// RecordAudit must not block or return an error to the caller; it swallows store
// failures silently (logged). Verify it returns immediately with a nil-shaped row.
func TestRecordAudit_DoesNotBlockOnStore(t *testing.T) {
	env := testutil.NewEnv(t)
	logic.ResetBaseLogicForTest()
	bl := logic.GetBaseLogic(env.SvcCtx)

	done := make(chan struct{})
	go func() {
		bl.RecordAudit(context.Background(), dto.OperationCreateDto{
			EntityType: constants.OperationEntityTypeUser,
			Operation:  constants.OperationActionUserPasswordEdit,
		})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("RecordAudit blocked the caller")
	}
}
