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

type SummaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SummaryLogic {
	return &SummaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SummaryLogic) Summary(in *borrower.LoanSummaryRequest) (resp *borrower.LoanSummaryResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	input, err := copierx.Trans[borrower.LoanSummaryRequestData, dto.LoanSummaryRequestDto](in.Data, copierx.StrToSqlTimePtr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	input.OrganizationId = claim.OrganizationId
	input.OrganizationUserId = claim.OrganizationUserId

	output, err := l.baseLogic.LoanLogic.BorrowerLoanSummary(l.ctx, *input)
	if err != nil {
		return nil, err
	}

	data, err := copierx.Trans[dto.LoanSummaryResponseDto, borrower.LoanSummaryResponseData](output,
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}

	return &borrower.LoanSummaryResponse{
		Data: data,
	}, nil
}
