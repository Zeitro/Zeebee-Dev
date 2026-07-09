package loanapplicantlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type BorrowerNamesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewBorrowerNamesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerNamesLogic {
	return &BorrowerNamesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BorrowerNamesLogic) BorrowerNames(in *proto.EmptyRequest) (resp *borrower.BorrowerNamesResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	names, err := l.svcCtx.LoanApplicantRepo.FindStrDistinctColumn(l.ctx, "full_name", map[string]interface{}{"borrower_organization_id": claim.OrganizationId})
	if err != nil {
		return nil, err
	}
	return &borrower.BorrowerNamesResponse{Data: names}, nil
}
