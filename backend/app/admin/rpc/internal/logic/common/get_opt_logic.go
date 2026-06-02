package commonlogic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
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

func (l *GetOptLogic) GetOpt(in *admin.OptionRequest) (resp *admin.OptionResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	resp = &admin.OptionResponse{
		Data: []*admin.OptionEntity{},
	}
	err = copier.Copy(&resp.Data, new(enums.Get(in.Data.Name)))
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return resp, nil
}
