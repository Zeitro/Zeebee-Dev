package projectmanagementlogic

import (
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/data/dto"
)

func memberEntity(member *dto.ProjectManagementMemberDetailDto) (*borrower.ProjectManagementMemberEntity, error) {
	entity, err := copierx.Trans[dto.ProjectManagementMemberDetailDto, borrower.ProjectManagementMemberEntity](member)
	if err != nil {
		return nil, err
	}
	entity.CreatedAt, entity.UpdatedAt = formatEntityTime(member.CreatedAt), formatEntityTime(member.UpdatedAt)
	return entity, nil
}
