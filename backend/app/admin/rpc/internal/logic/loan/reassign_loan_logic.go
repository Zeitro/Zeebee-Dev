package loanlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type ReassignLoanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewReassignLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReassignLoanLogic {
	return &ReassignLoanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *ReassignLoanLogic) ReassignLoan(in *admin.ReassignLoanRequest) (resp *proto.EmptyResponse, err error) {

	claim, err := l.baseLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if in.Data == nil {
		return nil, errorx.New("Data is empty")
	}

	err = l.baseLogic.LoanLogic.ReassignLender(l.ctx, in.Data.Id, in.Data.LenderOrganizationId, claim)
	if err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}
