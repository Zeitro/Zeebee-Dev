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

type OrganizationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewOrganizationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationDetailLogic {
	return &OrganizationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationDetailLogic) OrganizationDetail(req *types.PathIdRequest) (resp *types.OrganizationEntity, err error) {

	// Call RPC server
	rpcResp, err := l.svcCtx.OrganizationRpcClient.Detail(l.ctx, &admin.IdRequest{
		Data: req.Id,
	})
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
