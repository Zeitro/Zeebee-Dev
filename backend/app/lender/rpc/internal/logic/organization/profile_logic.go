package organizationlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type ProfileLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
	logx.Logger
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
		Logger:    logx.WithContext(ctx),
	}
}

func (l *ProfileLogic) Profile(in *lender.EmptyRequest) (resp *lender.OrganizationDetailResponse, err error) {
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

	// Validate user belongs to an organization
	if claim.OrganizationId == 0 {
		return nil, errorx.New("User Not Belong To Any Organization")
	}

	// Get organization details with preloaded associations (emails, etc.)
	organization, err := l.baseLogic.OrganizationLogic.DetailWithType(
		l.ctx,
		claim.OrganizationId,
		constants.UserTypeLender,
	)
	if err != nil {
		return nil, errorx.Wrap(err, "Failed To Get Organization")
	}

	// Convert business object to proto entity with time field conversion
	data, err := copierx.Trans[bo.OrganizationBo, lender.OrganizationEntity](
		organization,
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	if claim.UserId == organization.OwnerId {
		data.IsEditable = true
	}

	return &lender.OrganizationDetailResponse{
		Data: data,
	}, nil
}
