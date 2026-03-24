package organizationsettinglogic

import (
	"context"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type DetailByPurposeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewDetailByPurposeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailByPurposeLogic {
	return &DetailByPurposeLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *DetailByPurposeLogic) DetailByPurpose(in *lender.OrganizationSettingDetailByPurposeRequest) (resp *lender.OrganizationSettingDetailResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	if in.Data == nil {
		return nil, errorx.New("Data is required")
	}

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.baseLogic.OrganizationSettingLogic.DetailByPurpose(l.ctx, claim.OrganizationId, in.Data.Purpose)
	if err != nil {
		return nil, err
	}
	data, err := copierx.Trans[dto.OrganizationSettingEntityDto, lender.OrganizationSettingEntity](result, copierx.TimeToRFC3339Str(), copierx.StrMapToBytes())
	if err != nil {
		return nil, err
	}

	return &lender.OrganizationSettingDetailResponse{
		Data: data,
	}, nil
}
