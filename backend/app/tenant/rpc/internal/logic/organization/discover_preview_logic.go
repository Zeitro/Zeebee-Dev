package organizationlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/schedule"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/lenderdiscover"
	"zeebee-backend/domains/zeebee/proto"
)

type DiscoverPreviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiscoverPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiscoverPreviewLogic {
	return &DiscoverPreviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiscoverPreviewLogic) DiscoverPreview(in *proto.EmptyRequest) (resp *tenant.LenderDiscoverPreviewResponse, err error) {
	rn := schedule.NewRunner(l.svcCtx)
	cands, err := rn.Preview(context.WithoutCancel(l.ctx))
	if err != nil {
		return nil, err
	}
	list := make([]*tenant.LenderDiscoverCandidateEntity, 0, len(cands))
	var newCount, existedCount int32
	for i := range cands {
		candidate, err := discoverCandidateToProto(&cands[i])
		if err != nil {
			return nil, err
		}
		list = append(list, candidate)
		if cands[i].DedupeStatus == lenderdiscover.DedupeNew {
			newCount++
		} else {
			existedCount++
		}
	}
	return &tenant.LenderDiscoverPreviewResponse{Data: &tenant.LenderDiscoverPreviewResponseData{
		Candidates: list, NewCount: newCount, ExistedCount: existedCount,
	}}, nil
}
