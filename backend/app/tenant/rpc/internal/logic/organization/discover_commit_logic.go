package organizationlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/schedule"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/lenderdiscover"
)

type DiscoverCommitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverCommitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverCommitLogic {
	return &DiscoverCommitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverCommitLogic) DiscoverCommit(in *tenant.LenderDiscoverCommitRequest) (resp *tenant.LenderDiscoverRunResponse, err error) {
	if in.Data == nil {
		err = errorx.New("request data is nil")
		return
	}
	lenders := in.GetData().GetLenders()
	raws := make([]lenderdiscover.RawLender, 0, len(lenders))
	for _, lender := range lenders {
		raw, err := discoverRawLenderFromCommitProto(lender)
		if err != nil {
			return nil, err
		}
		raws = append(raws, raw)
	}
	rn := schedule.NewRunner(l.svcCtx)
	result, err := rn.Commit(context.WithoutCancel(l.ctx), raws)
	if err != nil {
		return nil, err
	}
	data, err := discoverRunResultToProto(result)
	if err != nil {
		return nil, err
	}
	return &tenant.LenderDiscoverRunResponse{Data: data}, nil
}
