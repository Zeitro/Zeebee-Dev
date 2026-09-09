package commonlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/uuid"
	"zeebee-backend/domains/zeebee/proto"
)

type GetUuidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUuidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUuidLogic {
	return &GetUuidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUuidLogic) GetUuid(in *proto.EmptyRequest) (resp *borrower.GetUuidResponse, err error) {

	uid, err := uuid.GetSonyflakeIDString()
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &borrower.GetUuidResponse{
		Data: uid,
	}, nil
}
