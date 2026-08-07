package storage

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
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
	rpcReq := &proto.StorageCreatePresignedUrlRequest{}

	rpcReq.Data, err = copierx.Trans[types.S3PresignUrlCreateRequest, proto.StorageCreatePresignedUrlRequestData](req)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.StorageRpcClient.CreatePresignedUrl(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}

	resp, err = copierx.Trans[proto.StorageCreatePresignedUrlResponseData, types.S3PresignUrlCreateResponse](rpcResp.Data)

	return
}
