package projectmanagementtemplatelogic

import (
	"testing"

	"github.com/stretchr/testify/require"

	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func TestProjectManagementTemplatePageRequestConversion(t *testing.T) {
	pageReq, err := copierx.Trans[admin.ProjectManagementTemplatePageRequestData, dto.ProjectManagementTemplatePageRequestDto](&admin.ProjectManagementTemplatePageRequestData{
		Page:   2,
		Limit:  20,
		Filter: &admin.ProjectManagementTemplatePageRequestFilter{Keyword: "template", Status: string(constants.CommonStatusActive)},
		Sorter: &proto.CommonPageSorter{Column: "name", Sort: "asc"},
	})
	require.NoError(t, err)
	require.Equal(t, int64(2), pageReq.Page)
	require.Equal(t, int64(20), pageReq.Limit)
	require.Equal(t, "template", pageReq.Filter.Keyword)
	require.Equal(t, constants.CommonStatusActive, pageReq.Filter.Status)
	require.Equal(t, "name", pageReq.Sorter.Column)
	require.Equal(t, "asc", pageReq.Sorter.Sort)
}

func TestProjectManagementTemplatePageListConversion(t *testing.T) {
	list, err := copierx.Trans[[]*model.ProjectManagementTemplate, []*proto.ProjectManagementTemplateEntity](&[]*model.ProjectManagementTemplate{{Name: "Template", Stages: []*model.ProjectManagementStage{{Name: "Pending"}}}})
	require.NoError(t, err)
	require.Len(t, *list, 1)
	require.Equal(t, "Template", (*list)[0].Name)
	require.Len(t, (*list)[0].Stages, 1)
	require.Equal(t, "Pending", (*list)[0].Stages[0].Name)
}
