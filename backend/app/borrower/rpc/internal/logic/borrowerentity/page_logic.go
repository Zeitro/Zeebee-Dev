package borrowerentitylogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
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

func (l *PageLogic) Page(in *proto.CommonPageRequest) (resp *borrower.BorrowerEntityPageResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageReq := &dto.BorrowerEntityPageRequestDto{}
	pageReq, err = copierx.Trans[proto.CommonPageRequestData, dto.BorrowerEntityPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	pageReq.Filter.OrganizationId = claim.OrganizationId
	pageReq.Filter.UserId = claim.OrganizationUserId
	pageResult, err := l.svcCtx.BorrowerEntityRepo.PageForBorrower(l.ctx, *pageReq)
	if err != nil {
		return nil, err
	}

	respData := make([]*borrower.BorrowerEntityPageEntity, 0, len(pageResult.List))
	respDataPtr, err := copierx.Trans[[]*model.BorrowerEntity, []*borrower.BorrowerEntityPageEntity](&pageResult.List,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}
	respData = *respDataPtr

	return &borrower.BorrowerEntityPageResponse{
		Data: &borrower.BorrowerEntityPageResponseData{
			Page:      pageResult.Page,
			Limit:     pageResult.Limit,
			Total:     pageResult.Total,
			TotalPage: pageResult.TotalPage,
			List:      respData,
		},
	}, nil
}
