package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type PageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *PageLogic) Page(in *borrower.OrganizationUserPageRequest) (resp *borrower.OrganizationUserPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
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
