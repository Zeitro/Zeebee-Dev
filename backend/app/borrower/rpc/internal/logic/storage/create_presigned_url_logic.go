package storagelogic

import (
	"context"

	"github.com/google/uuid"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type CreatePresignedUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePresignedUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePresignedUrlLogic {
	return &CreatePresignedUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePresignedUrlLogic) CreatePresignedUrl(in *borrower.StorageCreatePresignedUrlRequest) (resp *borrower.StorageCreatePresignedUrlResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	folder := global.Config.AwsS3Conf.FileFolder
	expireSeconds := global.Config.AwsS3Conf.PreSigneUrlExpireSeconds

	client := global.GetS3Client()
	typ := in.Data.Type
	if len(typ) == 0 {
		typ = "pdf"
	}
	name := uuid.New().String()
	key := folder + "/" + in.Data.Type + "/" + name
	req, err := client.GetPresignUrl(l.ctx, key, expireSeconds)
	if err != nil {
		return
	}

	return &borrower.StorageCreatePresignedUrlResponse{
		Data: &borrower.StorageCreatePresignedUrlResponseData{
			Url: req.URL,
		},
	}, nil
}
