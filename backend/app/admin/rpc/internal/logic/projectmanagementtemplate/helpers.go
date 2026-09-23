package projectmanagementtemplatelogic

import (
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func templateEntity(template *model.ProjectManagementTemplate) (*proto.ProjectManagementTemplateEntity, error) {
	data, err := copierx.Trans[model.ProjectManagementTemplate, proto.ProjectManagementTemplateEntity](template, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
	if err != nil {
		return nil, err
	}
	stages, err := copierx.Trans[[]*model.ProjectManagementStage, []*proto.ProjectManagementStageEntity](&template.Stages, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	data.Stages = *stages
	return data, nil
}

func templateResponse(template *model.ProjectManagementTemplate) (*admin.ProjectManagementTemplateDetailResponse, error) {
	data, err := templateEntity(template)
	if err != nil {
		return nil, err
	}
	return &admin.ProjectManagementTemplateDetailResponse{Data: data}, nil
}

func templateWorkflowResponse(workflow *model.ProjectManagementWorkflow) (*admin.ProjectManagementTemplateWorkflowDetailResponse, error) {
	data, err := copierx.Trans[model.ProjectManagementWorkflow, proto.ProjectManagementWorkflowEntity](workflow, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return &admin.ProjectManagementTemplateWorkflowDetailResponse{Data: data}, nil
}
