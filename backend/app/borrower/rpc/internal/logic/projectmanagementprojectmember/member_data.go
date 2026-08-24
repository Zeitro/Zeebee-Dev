package projectmanagementprojectmemberlogic

import (
	"zeebee-backend/app/borrower/rpc/pb/borrower"
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

func memberCandidateEntity(candidate *dto.ProjectManagementMemberCandidateEntityDto) (*borrower.ProjectManagementMemberCandidateEntity, error) {
	entity, err := copierx.Trans[dto.ProjectManagementMemberCandidateEntityDto, borrower.ProjectManagementMemberCandidateEntity](candidate)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
