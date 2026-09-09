package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/schedule"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
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

func (l *DiscoverCommitLogic) DiscoverCommit(in *admin.LenderDiscoverCommitRequest) (resp *admin.LenderDiscoverRunResponse, err error) {

	var lenders []*admin.LenderDiscoverCommitLender
	if in.GetData() != nil {
		lenders = in.GetData().GetLenders()
	}
	raws := make([]domainlenderdiscover.RawLender, 0, len(lenders))
	for _, p := range lenders {
		raw, transErr := discoverRawLenderFromCommitProto(p)
		if transErr != nil {
			return nil, transErr
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
	return &admin.LenderDiscoverRunResponse{Data: data}, nil
}
