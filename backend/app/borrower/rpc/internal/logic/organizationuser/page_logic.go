package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *borrower.OrganizationUserPageRequest) (resp *borrower.OrganizationUserPageResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	pageReq := dto.OrganizationUserPageRequestDto{}
	err = copierx.Assign[borrower.OrganizationUserPageRequestData, dto.OrganizationUserPageRequestDto](in.Data, &pageReq)
	if err != nil {
		return nil, err
	}
	pageReq.Filter.UserType = claim.UserType
	pageReq.Filter.OrganizationId = claim.OrganizationId
	pageReq.Filter.UserId = claim.UserId

	results, err := l.baseLogic.OrganizationUserLogic.Page(l.ctx, pageReq)
	if err != nil {
		return nil, err
	}
	data := &borrower.OrganizationUserPageResponseData{}
	data, err = copierx.Trans[dto.PageResultDto[*dto.OrganizationUserEntityDto], borrower.OrganizationUserPageResponseData](&results, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationUserPageResponse{
		Data: data,
	}, nil
}
