package projectmanagementlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	borrowerlogic "zeebee-backend/app/borrower/rpc/internal/logic"
)

func TestRequireProjectManagementActorRejectsMissingTenantIdentity(t *testing.T) {
	_, err := requireProjectManagementActor(context.Background(), &borrowerlogic.BaseLogic{})

	require.Equal(t, codes.Unauthenticated, status.Code(err))
}
