package commonlogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/lender/rpc/internal/logic"
	"zeebee-backend/app/lender/rpc/internal/svc"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SyncDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	baseLogic *logic.BaseLogic
}

func NewSyncDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDataLogic {
	return &SyncDataLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SyncDataLogic) SyncData(in *proto.EmptyRequest) (resp *proto.EmptyResponse, err error) {
	defer func() {
		if err != nil {
			l.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
				logx.Field("service_name", l.svcCtx.Config.Name),
			).Error(err)
		}
	}()
	organizations, err := l.svcCtx.OrganizationRepo.FindAll(l.ctx, map[string]interface{}{
		"type": constants.UserTypeLender,
	})
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.OrganizationSettingRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		var innerErr error
		var isExisted bool
		organizationSettings := make([]*model.OrganizationSetting, 0, len(organizations))
		for _, organization := range organizations {
			isExisted, innerErr = l.svcCtx.OrganizationSettingRepo.IsExistedByColumnsWithTx(l.ctx, tx, map[string]interface{}{
				"organization_id": organization.Id,
				"purpose":         constants.OrganizationPurposePreapprovalLetter,
			})
			if innerErr != nil {
				return innerErr
			}
			if isExisted {
				continue
			}
			setting := &model.OrganizationSetting{
				OrganizationId: organization.Id,
				Purpose:        constants.OrganizationPurposePreapprovalLetter,
				Setting:        nil,
			}
			organizationSettings = append(organizationSettings, setting)
		}
		if len(organizationSettings) == 0 {
			return nil
		}
		_, innerErr = l.svcCtx.OrganizationSettingRepo.BatchCreateWithTx(l.ctx, tx, organizationSettings, false)

		return innerErr
	})

	return &proto.EmptyResponse{}, errorx.WithStack(err)
}
