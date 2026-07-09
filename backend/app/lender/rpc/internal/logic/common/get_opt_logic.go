package commonlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/enums"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

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

func (l *GetOptLogic) GetOpt(in *lender.OptionRequest) (resp *lender.OptionResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()
	resp = &lender.OptionResponse{
		Data: []*proto.OptionEntity{},
	}
	optsPtr, err := copierx.Trans[[]enums.EnumOpt, []*proto.OptionEntity](new(enums.Get(in.Data.Name)))
	if err != nil {
		return nil, err
	}
	resp.Data = *optsPtr

	return resp, nil
}
