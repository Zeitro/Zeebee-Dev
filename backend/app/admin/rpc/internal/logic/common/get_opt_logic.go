package commonlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/enums"
	"zeebee-backend/domains/zeebee/proto"

	"zeebee-backend/common/errorx"
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

func (l *GetOptLogic) GetOpt(in *admin.OptionRequest) (resp *admin.OptionResponse, err error) {

	resp = &admin.OptionResponse{
		Data: []*proto.OptionEntity{},
	}
	err = copier.Copy(&resp.Data, new(enums.Get(in.Data.Name)))
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
