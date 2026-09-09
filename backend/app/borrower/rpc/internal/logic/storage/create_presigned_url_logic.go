package storagelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
)

type CreatePresignedUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePresignedUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePresignedUrlLogic {
	return &CreatePresignedUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePresignedUrlLogic) CreatePresignedUrl(in *proto.StorageCreatePresignedUrlRequest) (resp *proto.StorageCreatePresignedUrlResponse, err error) {

	folder := global.Config.AwsS3Conf.FileFolder
	expireSeconds := global.Config.AwsS3Conf.PreSigneUrlExpireSeconds

	client := global.GetS3Client()
	typ := in.Data.Type
	suffix := strings.ToLower(in.Data.FileSuffix)
	if len(typ) == 0 {
		err = errorx.New("type is empty")
	}
	name := uuid.New().String()

	key := fmt.Sprintf("%s/%s/%s.%s", folder, typ, name, suffix)

	fileName := in.Data.FileName
	if len(fileName) == 0 {
		fileName = fmt.Sprintf("%s.%s", name, suffix)
	}

	// key := folder + "/" + in.Data.Type + "/" + name
	req, err := client.GetPresignUrl(l.ctx, key, expireSeconds, suffix)
	if err != nil {
		return
	}

	return &proto.StorageCreatePresignedUrlResponse{
		Data: &proto.StorageCreatePresignedUrlResponseData{
			Url:         req.URL,
			FileName:    fileName,
			StorageType: string(constants.StorageTypeAwsS3),
			UploadType:  string(constants.UploadTypeManual),
		},
	}, nil
}
