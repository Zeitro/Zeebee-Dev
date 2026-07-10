package storage

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

type CreateAwsS3PresignUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create aws s3 presign url
func NewCreateAwsS3PresignUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAwsS3PresignUrlLogic {
	return &CreateAwsS3PresignUrlLogic{
		Logger: svcCtx.Logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAwsS3PresignUrlLogic) CreateAwsS3PresignUrl(req *types.S3PresignUrlCreateRequest) (resp *types.S3PresignUrlCreateResponse, err error) {
	rpcReq := &admin.StorageCreatePresignedUrlRequest{}

	rpcReq.Data, err = copierx.Trans[types.S3PresignUrlCreateRequest, admin.StorageCreatePresignedUrlRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.StorageRpcClient.CreatePresignedUrl(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[admin.StorageCreatePresignedUrlResponseData, types.S3PresignUrlCreateResponse](rpcResp.Data)

	return
}
