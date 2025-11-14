package storagelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"zeebee-backend/app/lender/rpc/internal/global"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"

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

func (l *CreatePresignedUrlLogic) CreatePresignedUrl(in *lender.StorageCreatePresignedUrlRequest) (resp *lender.StorageCreatePresignedUrlResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	folder := global.Config.AwsS3Conf.FileFolder
	expireSeconds := global.Config.AwsS3Conf.PreSigneUrlExpireSeconds

	client := global.GetS3Client()
	typ := in.Data.Type
	suffix := strings.ToLower(in.Data.FileSuffix)
	if len(typ) == 0 {
		err = errors.New("Type Is Empty")
	}
	name := uuid.New().String()

	key := fmt.Sprintf("%s/%s/%s.%s", folder, typ, name, suffix)

	req, err := client.GetPresignUrl(l.ctx, key, expireSeconds, suffix)
	if err != nil {
		return
	}

	return &lender.StorageCreatePresignedUrlResponse{
		Data: &lender.StorageCreatePresignedUrlResponseData{
			Url: req.URL,
		},
	}, nil
}
