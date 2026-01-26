package basetemplatelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/admin/rpc/internal/logic"
	"zeebee-backend/app/admin/rpc/internal/svc"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
)

type ExistsByTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewExistsByTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExistsByTagLogic {
	return &ExistsByTagLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.NewBaseLogic(ctx, svcCtx),
	}
}

// ExistsByTag checks if a template tag exists within the same type and category
func (l *ExistsByTagLogic) ExistsByTag(in *admin.ExistsByTagRequest) (resp *admin.ExistsByTagResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()

	// Check tag existence (proto only has tag and type, category check removed)
	exists, err := l.baseLogic.BaseTemplateLogic.ExistsByTag(
		l.ctx,
		in.Data.Tag,
		constants.UserType(in.Data.Type),
	)
	if err != nil {
		return nil, err
	}

	return &admin.ExistsByTagResponse{
		Data: &admin.ExistsByTagResponseData{
			Exists: exists,
		},
	}, nil
}
