package project_management_project_member

import (
	"context"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// detail
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.ProjectManagementMemberDetailRequest) (resp *types.ProjectManagementMemberDetailResponse, err error) {
	data, err := copierx.Trans[types.ProjectManagementMemberDetailRequest, borrower.ProjectManagementMemberDetailRequestData](req)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ProjectManagementProjectMemberRpcClient.Detail(l.ctx, &borrower.ProjectManagementMemberDetailRequest{Data: data})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	entity, err := copierx.Trans[proto.ProjectManagementMemberEntity, types.ProjectManagementMemberEntity](rpcResp.Data)
	if err != nil {
		return nil, err
	}
	return &types.ProjectManagementMemberDetailResponse{ProjectManagementMemberEntity: *entity}, nil
}
