package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserSendEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// send email
func NewOrganizationUserSendEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserSendEmailLogic {
	return &OrganizationUserSendEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserSendEmailLogic) OrganizationUserSendEmail(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.OrganizationUserSendInvitationEmailRequest{
		Data: &borrower.OrganizationUserSendInvitationEmailRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.PathIdRequest, borrower.OrganizationUserSendInvitationEmailRequestData](req)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.SendInvitationEmail(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	resp = &types.EmptyResponse{}
	return resp, nil
}
