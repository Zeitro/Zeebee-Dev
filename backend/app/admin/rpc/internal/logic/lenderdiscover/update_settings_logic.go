package lenderdiscover

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
)

type UpdateSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSettingsLogic {
	return &UpdateSettingsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateSettingsLogic) UpdateSettings(in *admin.LenderDiscoverUpdateSettingsRequest) (resp *admin.LenderDiscoverSettingsResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	if in.GetData() == nil {
		return nil, errorx.New("Request data is required")
	}
	d := in.GetData()
	if d.IntervalDays < 1 {
		return nil, errorx.New("interval_days must be at least 1")
	}
	row, err := l.svcCtx.LenderDiscoverSettingsRepo.GetSingleton(l.ctx)
	if err != nil {
		return nil, err
	}
	cols := map[string]interface{}{
		"schedule_enabled":     d.ScheduleEnabled,
		"interval_days":        d.IntervalDays,
		"scope_state":          d.ScopeState,
		"scope_extra":          d.ScopeExtra,
		"user_prompt":          d.UserPrompt,
		"openai_model":         d.OpenaiModel,
		"updated_by_user_id":   d.UpdatedByUserId,
	}
	err = l.svcCtx.LenderDiscoverSettingsRepo.UpdateById(l.ctx, row.Id, cols)
	if err != nil {
		return nil, err
	}
	updated, err := l.svcCtx.LenderDiscoverSettingsRepo.GetSingleton(l.ctx)
	if err != nil {
		return nil, err
	}
	return &admin.LenderDiscoverSettingsResponse{Data: settingsToProto(updated)}, nil
}
