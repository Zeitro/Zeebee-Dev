package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type LenderPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewLenderPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderPageLogic {
	return &LenderPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *LenderPageLogic) LenderPage(in *proto.CommonPageRequest) (resp *borrower.OrganizationLenderPageResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	pageReq, err := copierx.Trans[proto.CommonPageRequestData, dto.OrganizationPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	pageReq.Filter.OrganizationType = constants.UserTypeLender
	pageReq.Filter.Status = constants.OrganizationStatusActive

	pageResult, err := l.baseLogic.OrganizationLogic.Page(l.ctx, *pageReq)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.PageResultDto[dto.OrganizationPageEntityDto], borrower.OrganizationLenderPageResponseData](&pageResult, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationLenderPageResponse{
		Data: data,
	}, nil
}
