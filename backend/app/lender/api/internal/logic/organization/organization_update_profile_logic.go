package organization

import (
	"context"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrganizationUpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update profile
func NewOrganizationUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrganizationUpdateProfileLogic {
	return &OrganizationUpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrganizationUpdateProfileLogic) OrganizationUpdateProfile(req *types.OrganizationUpdateProfileRequest) (resp *types.OrganizationProfileEntity, err error) {
	// Convert API request to proto request
	rpcReqData, err := copierx.Trans[types.OrganizationUpdateProfileRequest, lender.OrganizationUpdateProfileRequestData](
		req,
		copierx.TimeToRFC3339Str(),
	)
	if err != nil {
		return nil, err
	}
	lenderMetadata, err := copierx.Trans[types.OrganizationLenderMetadataEntity, proto.OrganizationLenderMetadataEntity](
		req.Metadata,
	)
	if err != nil {
		return nil, err
	}
	if rpcReqData.Metadata == nil {
		rpcReqData.Metadata = &proto.OrganizationMetadataEntity{}
	}
	rpcReqData.Metadata.Lender = lenderMetadata

	// Build RPC request
	rpcReq := &lender.OrganizationUpdateProfileRequest{
		Data: rpcReqData,
	}

	// Call RPC client to update organization profile
	rpcResp, err := l.svcCtx.OrganizationRpcClient.UpdateProfile(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	// Convert proto response to API types
	resp, err = copierx.Trans[proto.OrganizationProfileEntity, types.OrganizationProfileEntity](
		rpcResp.Data,
		copierx.TimeToRFC3339Str(),
		copierx.BytesToStrMap(),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
