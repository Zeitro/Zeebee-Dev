package organizationlogic

import (
	"context"

	"zeebee-backend/app/admin/rpc/internal/schedule"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	domainlenderdiscover "zeebee-backend/domains/zeebee/lenderdiscover"
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

func (l *DiscoverPreviewLogic) DiscoverPreview(_ *proto.EmptyRequest) (resp *admin.LenderDiscoverPreviewResponse, err error) {

	rn := schedule.NewRunner(l.svcCtx)
	cands, err := rn.Preview(context.WithoutCancel(l.ctx))
	if err != nil {
		return nil, err
	}
	var protoList []*admin.LenderDiscoverCandidateEntity
	var newN, existedN int32
	for i := range cands {
		c := &cands[i]
		ent, transErr := discoverCandidateToProto(c)
		if transErr != nil {
			return nil, transErr
		}
		protoList = append(protoList, ent)
		if c.DedupeStatus == domainlenderdiscover.DedupeNew {
			newN++
		} else {
			existedN++
		}
	}
	return &admin.LenderDiscoverPreviewResponse{
		Data: &admin.LenderDiscoverPreviewResponseData{
			Candidates:   protoList,
			NewCount:     newN,
			ExistedCount: existedN,
		},
	}, nil
}
