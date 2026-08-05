package loanlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type BorrowerLoanPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewBorrowerLoanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerLoanPageLogic {
	return &BorrowerLoanPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BorrowerLoanPageLogic) BorrowerLoanPage(in *borrower.LoanPageRequest) (resp *borrower.LoanPageResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	pageReq := &dto.LoanBorrowerPageRequestDto{}
	pageReq, err = copierx.Trans[borrower.LoanPageRequestData, dto.LoanBorrowerPageRequestDto](in.Data, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}

	pageReq.Filter.OrganizationId = claim.OrganizationId
	pageReq.Filter.UserId = claim.OrganizationUserId

	pageResult, err := l.baseLogic.LoanLogic.BorrowerLoanPage(l.ctx, *pageReq)
	if err != nil {
		return nil, err
	}

	resp = &borrower.LoanPageResponse{
		Data: &borrower.LoanPageResponseData{},
	}
	resp.Data, err = copierx.Trans[dto.PageResultDto[*dto.LoanEntityDto], borrower.LoanPageResponseData](pageResult,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
