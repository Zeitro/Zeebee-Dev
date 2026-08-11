package template_loader

import (
	"context"
	"errors"
	"strings"

	tenantclient "zeebee-backend/app/tenant/rpc/client"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/proto"
	"zeebee-backend/domains/zeebee/templater/loader"
)

type TenantLoader struct {
	systemTemplateLoader loader.TemplateLoader
	tenantClient         *tenantclient.Client
}

func New(baseTemplateLogic *logic.BaseTemplateLogic, tenantClient *tenantclient.Client) *TenantLoader {
	return &TenantLoader{systemTemplateLoader: loader.NewDefaultTemplateLoader(baseTemplateLogic, nil), tenantClient: tenantClient}
}

func (l *TenantLoader) LoadSystemTemplate(ctx context.Context, purpose string, organizationType string) (string, map[string]interface{}, error) {
	return l.systemTemplateLoader.LoadSystemTemplate(ctx, purpose, organizationType)
}

func (l *TenantLoader) LoadOrganizationTemplate(ctx context.Context, purpose string, organizationId int64) (string, map[string]interface{}, error) {
	tpl, err := l.tenantClient.OrganizationTemplate.DetailByPurpose(ctx, &tenant.OrganizationTemplateDetailByPurposeRequest{Data: &tenant.OrganizationTemplateDetailByPurposeRequestData{
		UserType: string(constants.UserTypeLender), OrganizationId: organizationId, Purpose: purpose,
	}})
	if err != nil {
		return "", nil, err
	}
	if tpl.Data == nil || tpl.Data.Status != string(constants.CommonStatusActive) {
		return "", nil, loader.ErrTemplateNotActive
	}
	if !tpl.Data.HasConfigured {
		return "", nil, errors.New("Template is not configured")
	}
	return tpl.Data.Content, fixedValues(tpl.Data.VariableGroupSchema), nil
}

func fixedValues(groups []*proto.TemplateVariableGroupEntity) map[string]interface{} {
	data := make(map[string]interface{})
	for _, group := range groups {
		for _, variable := range group.Variables {
			if variable.Source == string(constants.TemplateVariableSourceFixed) || variable.Source == string(constants.TemplateVariableSourceTemplate) {
				data[strings.Trim(variable.Key, ".")] = variable.FixedValue
			}
		}
	}
	return data
}
