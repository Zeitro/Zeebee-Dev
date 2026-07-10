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

type CheckInvitationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// check invitation
func NewCheckInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInvitationLogic {
	return &CheckInvitationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckInvitationLogic) CheckInvitation(req *types.OrganizationUserCheckInvitationRequest) (resp *types.OrganizationUserCheckInvitationResponse, err error) {
	rpcReq := &borrower.OrganizationUserCheckInvitationRequest{
		Data: &borrower.OrganizationUserCheckInvitationRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationUserCheckInvitationRequest, borrower.OrganizationUserCheckInvitationRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.CheckInvitation(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	resp = &types.OrganizationUserCheckInvitationResponse{}
	resp, err = copierx.Trans[borrower.OrganizationUserCheckInvitationResponseData, types.OrganizationUserCheckInvitationResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
