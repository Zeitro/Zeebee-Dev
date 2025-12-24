package loanlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SendPreApprovalNotificationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendPreApprovalNotificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPreApprovalNotificationLogic {
	return &SendPreApprovalNotificationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendPreApprovalNotificationLogic) SendPreApprovalNotification(in *borrower.LoanSendPreApprovalNotificationRequest) (resp *borrower.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	go func() {
		ctx := context.WithoutCancel(l.ctx)
		err := sendPreApprovalNotification(ctx, l.svcCtx.OrganizationRepo, in.Data.RequestId, in.Data.LenderCompanyId, true)
		if err != nil {
			l.WithContext(ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	return &borrower.EmptyResponse{}, nil
}
