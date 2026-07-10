package loanlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"
	"zeebee-backend/domains/zeebee/repo"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type LoanDocumentOptsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoanDocumentOptsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoanDocumentOptsLogic {
	return &LoanDocumentOptsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoanDocumentOptsLogic) LoanDocumentOpts(in *lender.LoanDocumentOptsRequest) (resp *lender.LoanDocumentOptsResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	opts := []repo.LoanDocumentOption{}
	resultOpts := []*proto.LoanDocumentOptEntity{}

	opts, err = l.svcCtx.LoanDocumentRepo.GetOptsByLoanId(l.ctx, in.Data.Id, []constants.LoanDocumentCategory{constants.LoanDocumentCategoryBorrower, constants.LoanDocumentCategoryExtend})
	if err != nil {
		return nil, err
	}

	resultOptsPtr, err := copierx.Trans[[]repo.LoanDocumentOption, []*proto.LoanDocumentOptEntity](&opts)
	if err != nil {
		return nil, err
	}
	resultOpts = *resultOptsPtr

	return &lender.LoanDocumentOptsResponse{
		Data: resultOpts,
	}, nil
}
