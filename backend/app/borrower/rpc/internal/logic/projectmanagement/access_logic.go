package projectmanagementlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/domains/zeebee/data/dto"
)

func requireProjectManagementActor(ctx context.Context, baseLogic *logic.BaseLogic) (*dto.ProjectManagementActor, error) {
	claim, err := baseLogic.GetUserClaimFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return baseLogic.ProjectManagementLogic.RequireActor(ctx, claim.OrganizationId, claim.UserId)
}
