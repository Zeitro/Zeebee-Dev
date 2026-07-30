package commonlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/enums"
	"zeebee-backend/domains/zeebee/proto"
)

type GetOptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOptLogic {
	return &GetOptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOptLogic) GetOpt(in *borrower.OptionRequest) (resp *borrower.OptionResponse, err error) {

	resp = &borrower.OptionResponse{
		Data: []*proto.OptionEntity{},
	}
	err = copierx.Assign[[]enums.EnumOpt, []*proto.OptionEntity](new(enums.Get(in.Data.Name)), &resp.Data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
