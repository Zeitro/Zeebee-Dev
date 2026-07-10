package organizationlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type ProfileLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ProfileLogic) Profile(in *proto.EmptyRequest) (resp *lender.OrganizationProfileResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	// Get current user claim from gRPC metadata (JWT token)
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get user claim from context")
	}

	// Validate user belongs to an organization
	if claim.OrganizationId == 0 {
		return nil, errorx.New("User does not belong to any organization")
	}

	// Get organization details with preloaded associations (emails, etc.)
	organization, err := l.baseLogic.OrganizationLogic.DetailWithType(
		l.ctx,
		claim.OrganizationId,
		constants.UserTypeLender,
	)
	if err != nil {
		return nil, errorx.Wrap(err, "fail to get organization")
	}

	// Convert business object to proto entity with time field conversion
	data, err := copierx.Trans[dto.OrganizationEntityDto, proto.OrganizationProfileEntity](
		organization,
		copierx.TimeToRFC3339Str(),
		copierx.StrMapToBytes(),
	)
	if err != nil {
		return nil, err
	}
	if claim.UserId == organization.OwnerId {
		data.IsEditable = true
	}
	for _, setting := range data.OrganizationSettings {
		if setting.Purpose == constants.OrganizationPurposePreapprovalLetter {
			preapprovalLetterTemplate, err := l.svcCtx.OrganizationTemplateRepo.DetailByPurpose(l.ctx, claim.OrganizationId, constants.OrganizationPurposePreapprovalLetter, false)
			if err != nil {
				return nil, err
			}
			setting.HasConfigured = preapprovalLetterTemplate.HasConfigured
		} else {
			setting.HasConfigured = true
		}
	}

	return &lender.OrganizationProfileResponse{
		Data: data,
	}, nil
}
