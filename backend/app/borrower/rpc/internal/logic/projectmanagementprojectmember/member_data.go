package projectmanagementprojectmemberlogic

import (
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

func memberEntity(member *dto.ProjectManagementMemberDetailDto) (*proto.ProjectManagementMemberEntity, error) {
	entity, err := copierx.Trans[dto.ProjectManagementMemberDetailDto, proto.ProjectManagementMemberEntity](member)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(member.CreatedAt), formatEntityTime(member.UpdatedAt)
	return entity, nil
}
