package organizationlogic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/domains/zeebee/constants"

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
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *UpdateProfileLogic) UpdateProfile(in *lender.OrganizationUpdateProfileRequest) (resp *lender.OrganizationDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Get current user claim from gRPC metadata (JWT token)
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get User Claim From Context")
	}

	// Get organization details to verify ownership
	organization, err := l.svcCtx.OrganizationRepo.GetOrganization(l.ctx, claim.OrganizationId, constants.UserTypeLender, false)

	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get Organization")
	}

	// Verify user is the organization owner
	if organization.OwnerId != claim.UserId {
		return nil, errorx.New("Only Owner Can Update Organization Information")
	}

	// Copy request fields to existing model (only copies non-zero fields)
	err = copier.CopyWithOption(&organization, in.Data, copier.Option{
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Copy Request Data")
	}

	// Update the model
	_, err = l.svcCtx.OrganizationRepo.UpdateOneById(l.ctx, claim.OrganizationId, organization)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Update Organization")
	}

	// Return updated organization by calling Profile again
	profileLogic := NewProfileLogic(l.ctx, l.svcCtx)
	return profileLogic.Profile(&lender.EmptyRequest{})
}
