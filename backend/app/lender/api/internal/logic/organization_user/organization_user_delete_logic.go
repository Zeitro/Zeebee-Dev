package organization_user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/domains/zeebee/proto"
)

type OrganizationUserDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// delete
func NewOrganizationUserDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUserDeleteLogic {
	return &OrganizationUserDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUserDeleteLogic) OrganizationUserDelete(req *types.PathIdRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &proto.PathIdRequest{
		Data: &proto.PathIdRequestData{
			Id: req.Id,
		},
	}

	_, err = l.svcCtx.OrganizationUserRpcClient.Delete(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
