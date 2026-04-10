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

type AcceptInvitationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// accept invitation
func NewAcceptInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AcceptInvitationLogic {
	return &AcceptInvitationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AcceptInvitationLogic) AcceptInvitation(req *types.OrganizationUserAcceptInvitationRequest) (resp *types.UserLoginResponse, err error) {
	rpcReq := &borrower.OrganizationUserAcceptInvitationRequest{
		Data: &borrower.OrganizationUserAcceptInvitationRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationUserAcceptInvitationRequest, borrower.OrganizationUserAcceptInvitationRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.AcceptInvitation(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	resp = &types.UserLoginResponse{}
	resp, err = copierx.Trans[borrower.UserLoginResponseData, types.UserLoginResponse](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}

	return resp, nil
}
