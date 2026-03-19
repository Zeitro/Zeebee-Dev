package loanlogic

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
	"zeebee-backend/domains/zeebee/model"
)

type BorrowerLoanPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewBorrowerLoanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerLoanPageLogic {
	return &BorrowerLoanPageLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BorrowerLoanPageLogic) BorrowerLoanPage(in *borrower.LoanPageRequest) (resp *borrower.LoanPageResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	pageReq := &dto.LoanBorrowerPageRequestDto{}
	pageReq, err = copierx.Trans[borrower.LoanPageRequestData, dto.LoanBorrowerPageRequestDto](in.Data)
	if err != nil {
		return nil, err
	}

	pageReq.Filter.OrganizationId = claim.OrganizationId
	pageReq.Filter.UserId = claim.OrganizationUserId

	pageResult, err := l.svcCtx.LoanRepo.PageForBorrower(l.ctx, *pageReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	resp = &borrower.LoanPageResponse{
		Data: &borrower.LoanPageResponseData{},
	}
	resp.Data, err = copierx.Trans[gormx.PageResponse[*model.Loan], borrower.LoanPageResponseData](&pageResult,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
