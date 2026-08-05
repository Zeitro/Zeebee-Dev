package organizationtemplatelogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type DetailByPurposeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
	// repo      *repo.AllRepos
}

func NewDetailByPurposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailByPurposeLogic {
	return &DetailByPurposeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailByPurposeLogic) DetailByPurpose(in *lender.OrganizationTemplateDetailByPurposeRequest) (resp *lender.OrganizationTemplateDetailResponse, err error) {

	claims, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeLender)
	if err != nil {
		return nil, err
	}
	detail, err := l.baseLogic.OrganizationTemplateLogic.DetailByPurpose(l.ctx, claims.OrganizationId, in.Data.Purpose)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationTemplateEntityDto, proto.OrganizationTemplateEntity](
		detail,
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationTemplateDetailResponse{
		Data: data,
	}, nil
}
