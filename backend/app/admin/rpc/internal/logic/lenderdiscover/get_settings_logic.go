package lenderdiscover

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
)

type GetSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSettingsLogic {
	return &GetSettingsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetSettingsLogic) GetSettings(_ *admin.EmptyRequest) (resp *admin.LenderDiscoverSettingsResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	row, err := l.svcCtx.LenderDiscoverSettingsRepo.GetSingleton(l.ctx)
	if err != nil {
		return nil, err
	}
	return &admin.LenderDiscoverSettingsResponse{Data: settingsToProto(row)}, nil
}
