package organization

import (
	"context"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationCreateLenderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// lender create
func NewOrganizationCreateLenderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationCreateLenderLogic {
	return &OrganizationCreateLenderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationCreateLenderLogic) OrganizationCreateLender(req *types.OrganizationCreateLenderRequest) (resp *types.OrganizationEntity, err error) {

	// Convert types to pb request
	rpcReq := &admin.OrganizationCreateLenderRequest{
		Data: &admin.OrganizationCreateLenderRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationCreateLenderRequest, admin.OrganizationCreateLenderRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	lenderMetadata, err := copierx.Trans[types.OrganizationLenderMetadataEntity, proto.OrganizationLenderMetadataEntity](req.Organization.Metadata)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.Organization.Metadata = &proto.OrganizationMetadataEntity{
		Lender: lenderMetadata,
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.OrganizationRpcClient.CreateLender(l.ctx, rpcReq)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Convert pb response to types
	resp, err = copierx.Trans[proto.OrganizationEntity, types.OrganizationEntity](rpcResp.Data, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
