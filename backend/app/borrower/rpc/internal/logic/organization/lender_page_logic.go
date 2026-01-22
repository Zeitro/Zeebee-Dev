package organizationlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type LenderPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLenderPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LenderPageLogic {
	return &LenderPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LenderPageLogic) LenderPage(in *borrower.OrganizationLenderPageRequest) (resp *borrower.OrganizationLenderPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	pageReq, err := copierx.Trans[borrower.OrganizationLenderPageRequestData, dto.OrganizationPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}
	pageReq.Filter.OrganizationType = constants.UserTypeLender
	pageReq.Filter.Status = constants.OrganizationStatusActive

	results, err := l.svcCtx.OrganizationRepo.Page(l.ctx, *pageReq)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[gormx.PageResponse[*model.Organization], borrower.OrganizationLenderPageResponseData](&results, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return &borrower.OrganizationLenderPageResponse{
		Data: data,
	}, nil
}
