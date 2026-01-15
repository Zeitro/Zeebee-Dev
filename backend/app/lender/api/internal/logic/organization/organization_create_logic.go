package organization

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create
func NewOrganizationCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationCreateLogic {
	return &OrganizationCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationCreateLogic) OrganizationCreate(req *types.OrganizationCreateRequest) (resp *types.OrganizationEntity, err error) {
	// Convert API request to proto request
	rpcReqData, err := copierx.Trans[types.OrganizationCreateRequest, lender.OrganizationCreateRequestData](
		req,
	)
	if err != nil {
		return nil, err
	}

	lenderMetadata, err := copierx.Trans[types.OrganizationLenderMetadataEntity, lender.OrganizationLenderMetadataEntity](
		req.Metadata,
	)
	if err != nil {
		return nil, err
	}
	if rpcReqData.Metadata == nil {
		rpcReqData.Metadata = &lender.OrganizationMetadataEntity{}
	}
	rpcReqData.Metadata.Lender = lenderMetadata

	// Build RPC request
	rpcReq := &lender.OrganizationCreateRequest{
		Data: rpcReqData,
	}

	// Call RPC client to create organization
	rpcResp, err := l.svcCtx.OrganizationRpcClient.Create(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// Convert proto response to API types
	resp, err = copierx.Trans[lender.OrganizationEntity, types.OrganizationEntity](
		rpcResp.Data,
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
