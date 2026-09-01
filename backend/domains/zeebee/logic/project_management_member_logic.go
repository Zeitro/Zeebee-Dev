package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type ProjectManagementMemberLogic struct{ *BaseLogic }

type projectManagementMemberOperation struct {
	claim                      utils.JwtUser
	projectId                  int64
	action                     constants.OperationAction
	targetOrganizationUserId   int64
	targetOrganizationUserName string
	before                     map[string]any
	after                      map[string]any
}

func NewProjectManagementMemberLogic(baseLogic *BaseLogic) *ProjectManagementMemberLogic {
	return &ProjectManagementMemberLogic{BaseLogic: baseLogic}
}

func (l *ProjectManagementMemberLogic) AddOne(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementMemberDto) error {
	return l.Add(ctx, claim, []dto.ProjectManagementMemberDto{input})
}

func (l *ProjectManagementMemberLogic) Detail(ctx context.Context, claim utils.JwtUser, projectId, memberId int64) (*dto.ProjectManagementMemberDetailDto, error) {
	if _, err := l.projectManagementLogic().RequireProject(ctx, projectId, claim.OrganizationId); err != nil {
		return nil, err
	}
	if _, err := l.projectManagementLogic().RequireMember(ctx, projectId, claim.OrganizationUserId); err != nil {
		return nil, err
	}
	member, err := l.repos.ProjectManagementProjectRepo.DetailMember(ctx, projectId, memberId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project member not found")
	}
	return copierx.Trans[repo.ProjectManagementProjectMemberEntity, dto.ProjectManagementMemberDetailDto](member)
}

func (l *ProjectManagementMemberLogic) Add(ctx context.Context, claim utils.JwtUser, inputs []dto.ProjectManagementMemberDto) error {
	if len(inputs) == 0 {
		return errorx.New("project members are required")
	}
	for _, input := range inputs {
		if err := validateProjectManagementMemberRole(input.Role, false); err != nil {
			return err
		}
	}
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, inputs[0].ProjectId); err != nil {
			return err
		}
		for _, input := range inputs {
			if input.ProjectId != inputs[0].ProjectId {
				return errorx.New("project members must belong to one project")
			}
			organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByIdWithTx(ctx, tx, input.OrganizationUserId, claim.OrganizationId)
			if err != nil {
				return projectManagementNotFound(err, "organization member not found")
			}
			if organizationUser.Status != constants.OrganizationUserStatusActive {
				return errorx.New("organization member is not active")
			}
			if _, err = l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, input.ProjectId, input.OrganizationUserId); err == nil {
				return errorx.New("project member already exists")
			} else if !isProjectManagementNotFound(err) {
				return err
			}
			member, err := copierx.Trans[dto.ProjectManagementMemberDto, model.ProjectManagementProjectMember](&input)
			if err != nil {
				return err
			}
			if err = l.repos.ProjectManagementProjectRepo.CreateMemberWithTx(ctx, tx, member); err != nil {
				return err
			}
			if err = l.recordOperation(ctx, tx, projectManagementMemberOperation{claim: claim, projectId: input.ProjectId, action: constants.OperationActionProjectManagementMemberAdd, targetOrganizationUserId: input.OrganizationUserId, targetOrganizationUserName: organizationUser.FullName, after: map[string]any{"organization_user_id": input.OrganizationUserId, "organization_user_full_name": organizationUser.FullName}}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (l *ProjectManagementMemberLogic) Update(ctx context.Context, claim utils.JwtUser, input dto.ProjectManagementMemberDto) error {
	if err := validateProjectManagementMemberRole(input.Role, false); err != nil {
		return err
	}
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, input.ProjectId); err != nil {
			return err
		}
		member, err := l.repos.ProjectManagementProjectRepo.LockMemberByIdWithTx(ctx, tx, input.ProjectId, input.MemberId)
		if err != nil {
			return projectManagementNotFound(err, "project member not found")
		}
		if member.Role == constants.ProjectManagementMemberRoleProjectOwner {
			return errorx.New("transfer project owner instead")
		}
		organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByIdWithTx(ctx, tx, member.OrganizationUserId, claim.OrganizationId)
		if err != nil && !isProjectManagementNotFound(err) {
			return projectManagementNotFound(err, "organization member not found")
		}
		organizationUserName := ""
		if organizationUser != nil {
			organizationUserName = organizationUser.FullName
		}
		before := memberOperationValues(member.OrganizationUserId, organizationUserName, member.Role)
		if err = l.repos.ProjectManagementProjectRepo.UpdateMemberRoleWithTx(ctx, tx, member, input.Role); err != nil {
			return err
		}
		after := memberOperationValues(member.OrganizationUserId, organizationUserName, input.Role)
		before, after = operationValuesDiff(before, after)
		if len(before) == 0 {
			return nil
		}
		return l.recordOperation(ctx, tx, projectManagementMemberOperation{claim: claim, projectId: input.ProjectId, action: constants.OperationActionProjectManagementMemberUpdate, targetOrganizationUserId: member.OrganizationUserId, targetOrganizationUserName: organizationUserName, before: before, after: after})
	})
}

// TransferOwner transfers project ownership to another member
func (l *ProjectManagementMemberLogic) TransferOwner(ctx context.Context, claim utils.JwtUser, projectId, targetOrganizationUserId int64) error {
	if projectId == 0 || targetOrganizationUserId == 0 {
		return errorx.New("project and target member are required")
	}
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		project, err := l.repos.ProjectManagementProjectRepo.LockByOrganizationWithTx(ctx, tx, projectId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project not found")
		}
		if project.OwnerId != claim.OrganizationUserId {
			return errorx.New("only project owner can transfer ownership")
		}
		currentOwner, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, projectId, claim.OrganizationUserId)
		if err != nil {
			return projectManagementNotFound(err, "project owner member not found")
		}
		if currentOwner.Role != constants.ProjectManagementMemberRoleProjectOwner {
			return errorx.New("only project owner can transfer ownership")
		}
		if targetOrganizationUserId == claim.OrganizationUserId {
			return errorx.New("target member is already project owner")
		}
		target, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, projectId, targetOrganizationUserId)
		if err != nil {
			return projectManagementNotFound(err, "target project member not found")
		}
		targetOrganizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByIdWithTx(ctx, tx, targetOrganizationUserId, claim.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "target organization member not found")
		}
		if err = l.repos.ProjectManagementProjectRepo.UpdateMemberRoleWithTx(ctx, tx, currentOwner, constants.ProjectManagementMemberRoleTaskExecutor); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementProjectRepo.UpdateMemberRoleWithTx(ctx, tx, target, constants.ProjectManagementMemberRoleProjectOwner); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementProjectRepo.UpdateOwnerWithTx(ctx, tx, projectId, targetOrganizationUserId, targetOrganizationUser.FullName); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementMemberOperation{claim: claim, projectId: projectId, action: constants.OperationActionProjectManagementOwnerTransfer, targetOrganizationUserId: targetOrganizationUserId, targetOrganizationUserName: targetOrganizationUser.FullName, before: memberOperationValues(currentOwner.OrganizationUserId, project.OwnerName, constants.ProjectManagementMemberRoleProjectOwner), after: memberOperationValues(targetOrganizationUserId, targetOrganizationUser.FullName, constants.ProjectManagementMemberRoleProjectOwner)})
	})
}

func (l *ProjectManagementMemberLogic) Delete(ctx context.Context, claim utils.JwtUser, projectId, memberId int64) error {
	if projectId == 0 || memberId == 0 {
		return errorx.New("project and member are required")
	}
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.projectManagementLogic().requireProjectManagerWithTx(ctx, tx, claim, projectId); err != nil {
			return err
		}
		member, err := l.repos.ProjectManagementProjectRepo.LockMemberByIdWithTx(ctx, tx, projectId, memberId)
		if err != nil {
			return projectManagementNotFound(err, "project member not found")
		}
		if member.Role == constants.ProjectManagementMemberRoleProjectOwner {
			return errorx.New("transfer project owner before removing member")
		}
		organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByIdWithTx(ctx, tx, member.OrganizationUserId, claim.OrganizationId)
		if err != nil && !isProjectManagementNotFound(err) {
			return projectManagementNotFound(err, "organization member not found")
		}
		organizationUserName := ""
		if organizationUser != nil {
			organizationUserName = organizationUser.FullName
		}
		if err = l.repos.ProjectManagementTaskRepo.DeleteUnfinishedAssigneesByProjectMemberWithTx(ctx, tx, projectId, member.OrganizationUserId); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementProjectRepo.DeleteMemberWithTx(ctx, tx, member); err != nil {
			return err
		}
		return l.recordOperation(ctx, tx, projectManagementMemberOperation{claim: claim, projectId: projectId, action: constants.OperationActionProjectManagementMemberRemove, targetOrganizationUserId: member.OrganizationUserId, targetOrganizationUserName: organizationUserName, before: memberOperationValues(member.OrganizationUserId, organizationUserName, member.Role)})
	})
}

func (l *ProjectManagementLogic) requireProjectManagerWithTx(ctx context.Context, tx *gorm.DB, claim utils.JwtUser, projectId int64) (*model.ProjectManagementProject, error) {
	project, err := l.repos.ProjectManagementProjectRepo.LockByOrganizationWithTx(ctx, tx, projectId, claim.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project not found")
	}
	member, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, projectId, claim.OrganizationUserId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project member required")
	}
	if member.Role != constants.ProjectManagementMemberRoleProjectOwner && member.Role != constants.ProjectManagementMemberRoleProjectAdmin {
		return nil, errorx.New("project member role is not allowed")
	}
	return project, nil
}

func validateProjectManagementMemberRole(role constants.ProjectManagementMemberRole, allowOwner bool) error {
	if role == constants.ProjectManagementMemberRoleProjectAdmin || role == constants.ProjectManagementMemberRoleTaskExecutor || (allowOwner && role == constants.ProjectManagementMemberRoleProjectOwner) {
		return nil
	}
	return errorx.New("project member role is not allowed")
}

func isProjectManagementNotFound(err error) bool {
	return errorx.Cause(err) == gorm.ErrRecordNotFound
}

func (l *ProjectManagementMemberLogic) recordOperation(ctx context.Context, tx *gorm.DB, input projectManagementMemberOperation) error {
	metadata, err := marshalProjectManagementOperationMetadata(projectManagementOperationMetadata{ProjectID: input.projectId, OrganizationUserID: input.targetOrganizationUserId, OrganizationUserFullName: input.targetOrganizationUserName}, input.before, input.after)
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: input.projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: input.claim.OrganizationUserId, OperatorName: input.claim.Username, OperatorType: constants.UserTypeBorrower, Operation: input.action, Metadata: metadata, OrganizationId: input.claim.OrganizationId})
}

func memberOperationValues(id int64, name string, role constants.ProjectManagementMemberRole) map[string]any {
	return map[string]any{"organization_user_id": id, "organization_user_full_name": name, "role": role}
}
