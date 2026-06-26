package organization_user

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUserCheckInvitationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// check invitation
func NewOrganizationUserCheckInvitationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserCheckInvitationLogic {
	return &OrganizationUserCheckInvitationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserCheckInvitationLogic) OrganizationUserCheckInvitation(req *types.OrganizationUserCheckInvitationRequest) (resp *types.OrganizationUserCheckInvitationResponse, err error) {
	rpcReq := &lender.OrganizationUserCheckInvitationRequest{
		Data: &lender.OrganizationUserCheckInvitationRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationUserCheckInvitationRequest, lender.OrganizationUserCheckInvitationRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.OrganizationUserRpcClient.CheckInvitation(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	resp = &types.OrganizationUserCheckInvitationResponse{}
	resp, err = copierx.Trans[lender.OrganizationUserCheckInvitationResponseData, types.OrganizationUserCheckInvitationResponse](rpcResp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
