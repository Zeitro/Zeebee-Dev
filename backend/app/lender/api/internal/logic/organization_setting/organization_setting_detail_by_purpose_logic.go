package organization_setting

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationSettingDetailByPurposeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Detail by purpose
func NewOrganizationSettingDetailByPurposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationSettingDetailByPurposeLogic {
	return &OrganizationSettingDetailByPurposeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationSettingDetailByPurposeLogic) OrganizationSettingDetailByPurpose(req *types.OrganizationSettingDetailByPurposeRequest) (resp *types.OrganizationSettingDetailResponse, err error) {
	rpcReq := &lender.OrganizationSettingDetailByPurposeRequest{}

	rpcReq.Data, err = copierx.Trans[types.OrganizationSettingDetailByPurposeRequest, lender.OrganizationSettingDetailByPurposeRequestData](req)
	if err != nil {
		return
	}

	rpcResp, err := l.svcCtx.OrganizationSettingRpcClient.DetailByPurpose(l.ctx, rpcReq)
	if err != nil {
		return
	}
	resp, err = copierx.Trans[proto.OrganizationSettingEntity, types.OrganizationSettingDetailResponse](rpcResp.Data, copierx.BytesToStrMap())

	return resp, err
}
