package organizationuserlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
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

func (l *PageLogic) Page(in *lender.OrganizationUserPageRequest) (resp *lender.OrganizationUserPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	pageReq := &dto.OrganizationUserPageRequestDto{}
	pageReq, err = copierx.Trans[lender.OrganizationUserPageRequestData, dto.OrganizationUserPageRequestDto](in.Data)
	// err = copier.Copy(&pageReq, &in.Data)
	if err != nil {
		return nil, err
	}

	pageReq.Filter.UserType = claim.UserType
	pageReq.Filter.OrganizationId = claim.OrganizationId
	pageReq.Filter.UserId = claim.UserId

	results, err := l.baseLogic.OrganizationUserLogic.Page(l.ctx, *pageReq)
	if err != nil {
		return nil, err
	}
	var list []*lender.OrganizationUserPageEntity
	for _, v := range results.List {
		item, err := copierx.Trans[dto.OrganizationUserEntityDto, lender.OrganizationUserPageEntity](v, copierx.TimeToRFC3339Str())
		if err != nil {
			return nil, errorx.GetFailedTransRespErr(err)
		}
		list = append(list, item)
	}

	return &lender.OrganizationUserPageResponse{
		Data: &lender.OrganizationUserPageResponseData{
			Page:      results.Page,
			Limit:     results.Limit,
			Total:     results.Total,
			TotalPage: results.TotalPage,
			List:      list,
		},
	}, nil
}
