package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserAcceptInvitationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// accept invitation
func NewOrganizationUserAcceptInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserAcceptInvitationLogic {
	return &OrganizationUserAcceptInvitationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserAcceptInvitationLogic) OrganizationUserAcceptInvitation(req *types.OrganizationUserAcceptInvitationRequest) (resp *types.UserLoginResponse, err error) {
	data, err := copierx.Trans[types.OrganizationUserAcceptInvitationRequest, lender.OrganizationUserAcceptInvitationRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcReq := &lender.OrganizationUserAcceptInvitationRequest{
		Data: data,
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.AcceptInvitation(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[lender.UserLoginResponseData, types.UserLoginResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
