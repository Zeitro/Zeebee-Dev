package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"

	"zeebee-backend/common/errorx"
)

type UpdateProfileLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *UpdateProfileLogic) UpdateProfile(in *lender.OrganizationUpdateProfileRequest) (resp *lender.OrganizationProfileResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()

	// Get current user claim from gRPC metadata (JWT token)
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get user claim from context")
	}

	// Get organization details to verify ownership
	organization, err := l.svcCtx.OrganizationRepo.GetOrganization(l.ctx, claim.OrganizationId, constants.UserTypeLender, false)

	if err != nil {
		return nil, errorx.Wrap(err, "fail to get organization")
	}

	// Verify user is the organization owner
	if organization.Id != claim.OrganizationId {
		return nil, errorx.New("you do not belong to this organization")
	}

	// Copy request fields to existing model (only copies non-zero fields)
	err = copierx.Assign[lender.OrganizationUpdateProfileRequestData, model.Organization](in.Data, organization)
	if err != nil {
		return nil, err
	}

	// Update the model
	_, err = l.svcCtx.OrganizationRepo.UpdateOneById(l.ctx, claim.OrganizationId, organization)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to update organization")
	}

	// Return updated organization by calling Profile again
	profileLogic := NewProfileLogic(l.ctx, l.svcCtx)
	return profileLogic.Profile(&lender.EmptyRequest{})
}
