package organization_user

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

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
	rpcReq := &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.PathIdRequest, proto.PathIdRequestData](req)
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
