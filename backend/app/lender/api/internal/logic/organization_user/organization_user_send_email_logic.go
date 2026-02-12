package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"

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
	rpcReq := &lender.OrganizationUserSendInvitationEmailRequest{
		Data: &lender.OrganizationUserSendInvitationEmailRequestData{
			Id: req.Id,
		},
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.SendInvitationEmail(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
