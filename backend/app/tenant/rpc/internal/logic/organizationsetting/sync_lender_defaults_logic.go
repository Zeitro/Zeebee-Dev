package organizationsettinglogic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type SyncLenderDefaultsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncLenderDefaultsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncLenderDefaultsLogic {
	return &SyncLenderDefaultsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SyncLenderDefaultsLogic) SyncLenderDefaults(*tenant.OrganizationSettingSyncLenderDefaultsRequest) (*proto.EmptyResponse, error) {
	organizations, err := l.svcCtx.OrganizationRepo.FindAll(l.ctx, map[string]interface{}{"type": constants.UserTypeLender})
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.OrganizationSettingRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
		settings := make([]*model.OrganizationSetting, 0, len(organizations))
		for _, organization := range organizations {
			exists, err := l.svcCtx.OrganizationSettingRepo.IsExistedByColumnsWithTx(l.ctx, tx, map[string]interface{}{
				"organization_id": organization.Id, "purpose": constants.OrganizationPurposePreapprovalLetter,
			})
			if err != nil {
				return err
			}
			if !exists {
				settings = append(settings, &model.OrganizationSetting{OrganizationId: organization.Id, Purpose: constants.OrganizationPurposePreapprovalLetter})
			}
		}
		if len(settings) == 0 {
			return nil
		}
		_, err := l.svcCtx.OrganizationSettingRepo.BatchCreateWithTx(l.ctx, tx, settings, false)
		return err
	})
	return &proto.EmptyResponse{}, err
}
