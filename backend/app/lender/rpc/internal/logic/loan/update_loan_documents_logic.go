package loanlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/domains/zeebee/proto"
)

type UpdateLoanDocumentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLoanDocumentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLoanDocumentsLogic {
	return &UpdateLoanDocumentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLoanDocumentsLogic) UpdateLoanDocuments(in *lender.LoanUpdateDocumentsRequest) (resp *proto.EmptyResponse, err error) {

	// todo: add your logic here and delete this line

	return &proto.EmptyResponse{}, nil
}
