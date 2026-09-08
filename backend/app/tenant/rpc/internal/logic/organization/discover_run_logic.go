package organizationlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/tenant/rpc/internal/schedule"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/lenderdiscover"
)

type DiscoverRunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverRunLogic {
	return &DiscoverRunLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverRunLogic) DiscoverRun(in *tenant.LenderDiscoverRunRequest) (resp *tenant.LenderDiscoverRunResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	trigger := lenderdiscover.TriggerManual
	if in.GetData() != nil && in.GetData().GetTrigger() == lenderdiscover.TriggerScheduled {
		trigger = lenderdiscover.TriggerScheduled
	}
	if trigger != lenderdiscover.TriggerScheduled {
		return nil, errorx.New("manual discover: use Preview then Commit from the admin UI")
	}
	newCtx := context.WithoutCancel(l.ctx)
	threading.GoSafe(func() {
		rn := schedule.NewRunner(l.svcCtx)
		result, err := rn.Run(newCtx, trigger)
		if err != nil {
			l.svcCtx.Logger.WithContext(newCtx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error("fail to run discover lender")
			return
		}
		data, err := discoverRunResultToProto(result)
		if err != nil {
			l.svcCtx.Logger.WithContext(newCtx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error("fail to map discover lender result")
			return
		}
		l.svcCtx.Logger.WithContext(newCtx).WithFields(logx.Field("data", data)).Info("success to run discover lender")
	})
	return &tenant.LenderDiscoverRunResponse{}, nil
}
