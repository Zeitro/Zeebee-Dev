package loanlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SendPreApprovalNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendPreApprovalNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPreApprovalNotificationLogic {
	return &SendPreApprovalNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendPreApprovalNotificationLogic) SendPreApprovalNotification(in *borrower.LoanSendPreApprovalNotificationRequest) (resp *proto.EmptyResponse, err error) {

	go func() {
		ctx := context.WithoutCancel(l.ctx)
		err := sendPreApprovalNotification(ctx, l.svcCtx.TenantClient, in.Data.RequestId, in.Data.LenderCompanyId, true)
		if err != nil {
			l.svcCtx.Logger.WithContext(ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	return &proto.EmptyResponse{}, nil
}
