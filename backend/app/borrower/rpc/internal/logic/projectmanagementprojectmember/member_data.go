package projectmanagementprojectmemberlogic

import (
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

func memberEntity(member *dto.ProjectManagementMemberDetailDto) (*proto.ProjectManagementMemberEntity, error) {
	entity, err := copierx.Trans[dto.ProjectManagementMemberDetailDto, proto.ProjectManagementMemberEntity](member, copierx.TimeToRFC3339Str())
	if err != nil {
		return nil, err
	}
	return entity, nil
}
