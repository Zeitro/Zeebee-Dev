package commonlogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/enums"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type GetOptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOptLogic {
	return &GetOptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOptLogic) GetOpt(in *lender.OptionRequest) (resp *lender.OptionResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	resp = &lender.OptionResponse{
		Data: []*lender.OptionEntity{},
	}
	optsPtr, err := copierx.Trans[[]enums.EnumOpt, []*lender.OptionEntity](new(enums.Get(in.Data.Name)))
	if err != nil {
		return nil, err
	}
	resp.Data = *optsPtr

	return resp, nil
}
