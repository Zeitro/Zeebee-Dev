package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/admin/rpc/internal/schedule"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
)

type DiscoverRunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverRunLogic {
	return &DiscoverRunLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DiscoverRunLogic) DiscoverRun(in *admin.LenderDiscoverRunRequest) (resp *admin.LenderDiscoverRunResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	trigger := domainlenderdiscover.TriggerManual
	if in.GetData() != nil && in.GetData().GetTrigger() == domainlenderdiscover.TriggerScheduled {
		trigger = domainlenderdiscover.TriggerScheduled
	}
	if trigger != domainlenderdiscover.TriggerScheduled {
		return nil, errorx.New("manual discover: use Preview then Commit from the admin UI")
	}
	newCtx := context.WithoutCancel(l.ctx)
	threading.GoSafe(func() {
		var err error
		defer func() {
			if err != nil {
				l.svcCtx.Logger.WithContext(newCtx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error("fail to run discover lender")
			}
		}()
		rn := schedule.NewRunner(l.svcCtx)
		result, err := rn.Run(newCtx, trigger)
		if err != nil {
			return
		}
		data, err := discoverRunResultToProto(result)
		if err != nil {
			return
		}
		l.svcCtx.Logger.WithContext(newCtx).WithFields(
			logx.Field("data", data),
		).Info("success ro run discover lender")
	})
	return &admin.LenderDiscoverRunResponse{}, nil
}
