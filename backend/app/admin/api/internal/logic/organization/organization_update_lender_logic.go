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

	// Convert types to pb request
	rpcReq := &admin.OrganizationUpdateLenderRequest{
		Data: &admin.OrganizationUpdateLenderRequestData{},
	}
	rpcReq.Data, err = copierx.Trans[types.OrganizationUpdateLenderRequest, admin.OrganizationUpdateLenderRequestData](req)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	if req.Owner != nil {
		rpcReq.Data.Owner = &admin.OrganizationCreateLenderUserEntity{
			Email:     req.Owner.Email,
			Password:  req.Owner.Password,
			FirstName: req.Owner.FirstName,
			LastName:  req.Owner.LastName,
		}
	}

	lenderMetadata, err := copierx.Trans[types.OrganizationLenderMetadataEntity, proto.OrganizationLenderMetadataEntity](req.Metadata)
	if err != nil {
		return nil, err
	}
	rpcReq.Data.Metadata = &proto.OrganizationMetadataEntity{
		Lender: lenderMetadata,
	}

	// Call RPC server
	rpcResp, err := l.svcCtx.OrganizationRpcClient.UpdateLender(l.ctx, rpcReq)
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
