package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUpdateLenderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender update
func NewOrganizationUpdateLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUpdateLenderLogic {
	return &OrganizationUpdateLenderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUpdateLenderLogic) OrganizationUpdateLender(req *types.OrganizationUpdateLenderRequest) (resp *types.OrganizationEntity, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Convert types to pb request
	rpcReq := &admin.OrganizationUpdateLenderRequest{
		Data: &admin.OrganizationUpdateLenderRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationUpdateLenderRequest, admin.OrganizationUpdateLenderRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.OrganizationRpcClient.UpdateLender(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[admin.OrganizationEntity, types.OrganizationEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
