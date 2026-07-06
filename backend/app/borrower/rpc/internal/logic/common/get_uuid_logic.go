package commonlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/uuid"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUuidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUuidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUuidLogic {
	return &GetUuidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUuidLogic) GetUuid(in *proto.EmptyRequest) (resp *borrower.GetUuidResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	uid, err := uuid.GetSonyflakeIDString()
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &borrower.GetUuidResponse{
		Data: uid,
	}, nil
}
