package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type OrganizationSettingLogic struct {
	repos *repo.AllRepos
}

// NewOrganizationSettingLogic creates a new instance of OrganizationSettingLogic.
// It initializes the logic with the provided repository dependencies.
func NewOrganizationSettingLogic(repos *repo.AllRepos) *OrganizationSettingLogic {
	return &OrganizationSettingLogic{
		repos: repos,
	}
}

// DetailByPurpose retrieves the details of an organization setting by its purpose and organization id.
func (l *OrganizationSettingLogic) DetailByPurpose(ctx context.Context, organizationId int64, purpose string) (result *dto.OrganizationSettingEntityDto, err error) {
	organizationSetting, err := l.repos.OrganizationSettingRepo.DetailByPurpose(ctx, organizationId, purpose, true)
	if err != nil {
		return nil, err
	}

	result, err = copierx.Trans[model.OrganizationSetting, dto.OrganizationSettingEntityDto](organizationSetting)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SwitchStatus updates the status of an organization setting.
func (l *OrganizationSettingLogic) SwitchStatus(ctx context.Context, id int64, organizationId int64, enabled bool) (err error) {
	organizationSetting, err := l.repos.OrganizationSettingRepo.Detail(ctx, id, true)
	if err != nil {
		return err
	}
	if organizationSetting.OrganizationId != organizationId {
		return errorx.New("Not own the organization setting")
	}

	if enabled {
		switch organizationSetting.Purpose {
		case constants.OrganizationPurposePreapprovalLetter:
			err = l.checkPreapprovalTemplate(ctx, organizationSetting.OrganizationId)
		default:
			err = errorx.Errorf("Invalid purpose: %s", organizationSetting.Purpose)
		}
	}
	if err != nil {
		return err
	}
	err = l.repos.OrganizationSettingRepo.UpdateColumnsById(ctx, id, map[string]interface{}{
		"enabled": enabled,
	})
	return err
}

// InitLenderOrganizationSettingWithTx initializes lender organization setting
func (l *OrganizationSettingLogic) InitLenderOrganizationSettingWithTx(ctx context.Context, tx *gorm.DB, organization *model.Organization) (err error) {
	if organization == nil {
		return errorx.New("Organization is empty")
	}

	if organization.Type != constants.UserTypeLender {
		return nil
	}
	var organizationSettings []*model.OrganizationSetting
	// Init preapproval letter template
	organizationSettings = append(organizationSettings, &model.OrganizationSetting{
		OrganizationId: organization.Id,
		Purpose:        constants.OrganizationPurposePreapprovalLetter,
		Enabled:        false,
		Setting:        nil,
	})
	_, err = l.repos.OrganizationSettingRepo.BatchCreateWithTx(ctx, tx, organizationSettings, false)

	return err
}

// checkPreapprovalTemplate checks preapproval template
func (l *OrganizationSettingLogic) checkPreapprovalTemplate(ctx context.Context, organization int64) (err error) {
	organizationTemplate, err := l.repos.OrganizationTemplateRepo.DetailByPurpose(ctx, organization, constants.OrganizationPurposePreapprovalLetter, false)
	if err != nil {
		return err
	}
	if organizationTemplate.HasConfigured {
		return nil
	}
	variableGroups := organizationTemplate.VariableGroupSchema
	for _, variableGroup := range variableGroups {
		for _, variable := range variableGroup.Variables {
			err = variable.ValidateRequiredField()
			if err != nil {
				err = errorx.New("Please configure organization preapproval letter template")
				return err
			}
		}
	}
	// Update HasConfigured
	organizationTemplate.HasConfigured = true
	_, err = l.repos.OrganizationTemplateRepo.UpdateOneById(ctx, organizationTemplate.Id, organizationTemplate)
	return err
}
