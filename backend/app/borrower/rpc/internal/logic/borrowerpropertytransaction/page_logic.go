package borrowerpropertytransactionlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/repo"
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

func (l *PageLogic) Page(in *borrower.BorrowerPropertyTransactionPageRequest) (resp *borrower.BorrowerPropertyTransactionPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageReq := dto.BorrowerPropertyTransactionPageRequestDto{
		BasePageRequestDto: dto.BasePageRequestDto{
			Page:  in.Data.Page,
			Limit: in.Data.Limit,
		},
		Filter: dto.BorrowerPropertyTransactionPageRequestFilter{
			OrganizationId: claim.OrganizationId,
		},
	}

	if in.Data.Filter != nil {
		pageReq.Filter.Keyword = in.Data.Filter.Keyword
		pageReq.Filter.Type = in.Data.Filter.Type
		pageReq.Filter.Status = in.Data.Filter.Status
	}

	pageResult, err := l.svcCtx.BorrowerPropertyTransactionRepo.Page(l.ctx, pageReq)
	if err != nil {
		return nil, err
	}
	var data *borrower.BorrowerPropertyTransactionPageResponseData
	data, err = copierx.Trans[gormx.PageResponse[repo.BorrowerPropertyTransactionPageRow], borrower.BorrowerPropertyTransactionPageResponseData](&pageResult, copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
	if err != nil {
		return
	}

	// list := make([]*borrower.BorrowerPropertyTransactionPageEntity, 0, len(pageResult.List))
	// for _, row := range pageResult.List {
	// 	entity, transErr := copierx.Trans[model.BorrowerPropertyTransaction, borrower.BorrowerPropertyTransactionPageEntity](
	// 		&model.BorrowerPropertyTransaction{
	// 			Base:             model.Base{Id: row.Id, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt},
	// 			BorrowerPropertyId: row.BorrowerPropertyId,
	// 			OrganizationId:   row.OrganizationId,
	// 			Code:             row.Code,
	// 			Type:             row.Type,
	// 			Status:           row.Status,
	// 			CoeDate:          row.CoeDate,
	// 			PurchaseAmount:   row.PurchaseAmount,
	// 			EmdAmount:        row.EmdAmount,
	// 			PropertyType:     row.PropertyType,
	// 			Pm:               row.Pm,
	// 			Strategy:         row.Strategy,
	// 			Note:             row.Note,
	// 		},
	// 		copierx.TimeToRFC3339Str(),
	// 		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	// 	)
	// 	if transErr != nil {
	// 		return nil, errorx.GetFailedTransRespErr(transErr)
	// 	}
	//
	// 	addrEntity, transErr := copierx.Trans[data.Address, borrower.AddressEntity](&row.Address)
	// 	if transErr == nil {
	// 		entity.Address = addrEntity
	// 	}
	//
	// 	list = append(list, entity)
	// }

	return &borrower.BorrowerPropertyTransactionPageResponse{
		Data: data,
	}, nil
}
