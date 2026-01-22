package storage

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
)

type CreateAwsS3PresignUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create aws s3 presign url
func NewCreateAwsS3PresignUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAwsS3PresignUrlLogic {
	return &CreateAwsS3PresignUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAwsS3PresignUrlLogic) CreateAwsS3PresignUrl(req *types.S3PresignUrlCreateRequest) (resp *types.S3PresignUrlCreateResponse, err error) {
	rpcReq := &admin.StorageCreatePresignedUrlRequest{
		Data: &admin.StorageCreatePresignedUrlRequestData{
			Type:       req.Type,
			FileSuffix: req.FileSuffix,
		},
	}
	rpcResp, err := l.svcCtx.StorageRpcClient.CreatePresignedUrl(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.S3PresignUrlCreateResponse{
		Url: rpcResp.Data.Url,
	}, nil
}
