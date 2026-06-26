package storagelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/global"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
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

func (l *CreatePresignedUrlLogic) CreatePresignedUrl(in *admin.StorageCreatePresignedUrlRequest) (resp *admin.StorageCreatePresignedUrlResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	folder := global.Config.AwsS3Conf.FileFolder
	expireSeconds := global.Config.AwsS3Conf.PreSigneUrlExpireSeconds

	client := global.GetS3Client()
	typ := in.Data.Type
	suffix := strings.ToLower(in.Data.FileSuffix)
	if len(typ) == 0 {
		err = errorx.New("Type is empty")
	}
	name := uuid.New().String()

	key := fmt.Sprintf("%s/%s/%s.%s", folder, typ, name, suffix)

	fileName := in.Data.FileName
	if len(fileName) == 0 {
		fileName = fmt.Sprintf("%s.%s", name, suffix)
	}

	req, err := client.GetPresignUrl(l.ctx, key, expireSeconds, suffix)
	if err != nil {
		return
	}

	return &admin.StorageCreatePresignedUrlResponse{
		Data: &admin.StorageCreatePresignedUrlResponseData{
			Url:         req.URL,
			FileName:    fileName,
			StorageType: string(constants.StorageTypeAwsS3),
			UploadType:  string(constants.UploadTypeManual),
		},
	}, nil
}
