package logic

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

func (l *ProjectManagementLogic) AddProjectMember(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementMemberDto) error {
	return l.AddProjectMembers(ctx, actor, []dto.ProjectManagementMemberDto{input})
}

func (l *ProjectManagementLogic) AddProjectMembers(ctx context.Context, actor *dto.ProjectManagementActor, inputs []dto.ProjectManagementMemberDto) error {
	if len(inputs) == 0 {
		return errorx.New("project members are required")
	}
	for _, input := range inputs {
		if err := validateProjectManagementMemberRole(input.Role, false); err != nil {
			return err
		}
	}
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.requireProjectManagerWithTx(ctx, tx, actor, inputs[0].ProjectId); err != nil {
			return err
		}
		for _, input := range inputs {
			if input.ProjectId != inputs[0].ProjectId {
				return errorx.New("project members must belong to one project")
			}
			organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByIdWithTx(ctx, tx, input.OrganizationUserId, actor.OrganizationId)
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
			if err = l.recordProjectMemberOperation(ctx, tx, actor, input.ProjectId, constants.OperationActionProjectManagementMemberAdd, input.OrganizationUserId); err != nil {
				return err
			}
		}
		return nil
	})
}

func (l *ProjectManagementLogic) UpdateProjectMemberRole(ctx context.Context, actor *dto.ProjectManagementActor, input dto.ProjectManagementMemberDto) error {
	if err := validateProjectManagementMemberRole(input.Role, false); err != nil {
		return err
	}
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.requireProjectManagerWithTx(ctx, tx, actor, input.ProjectId); err != nil {
			return err
		}
		member, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, input.ProjectId, input.OrganizationUserId)
		if err != nil {
			return projectManagementNotFound(err, "project member not found")
		}
		if member.Role == constants.ProjectManagementMemberRoleProjectOwner {
			return errorx.New("transfer project owner instead")
		}
		if err = l.repos.ProjectManagementProjectRepo.UpdateMemberRoleWithTx(ctx, tx, member, input.Role); err != nil {
			return err
		}
		return l.recordProjectMemberOperation(ctx, tx, actor, input.ProjectId, constants.OperationActionProjectManagementMemberUpdate, input.OrganizationUserId)
	})
}

func (l *ProjectManagementLogic) TransferProjectOwner(ctx context.Context, actor *dto.ProjectManagementActor, projectId, targetOrganizationUserId int64) error {
	if projectId == 0 || targetOrganizationUserId == 0 {
		return errorx.New("project and target member are required")
	}
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		project, err := l.repos.ProjectManagementProjectRepo.LockByOrganizationWithTx(ctx, tx, projectId, actor.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "project not found")
		}
		if project.OwnerId != actor.OrganizationUserId {
			return errorx.New("only project owner can transfer ownership")
		}
		currentOwner, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, projectId, actor.OrganizationUserId)
		if err != nil {
			return projectManagementNotFound(err, "project owner member not found")
		}
		if currentOwner.Role != constants.ProjectManagementMemberRoleProjectOwner {
			return errorx.New("only project owner can transfer ownership")
		}
		if targetOrganizationUserId == actor.OrganizationUserId {
			return errorx.New("target member is already project owner")
		}
		target, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, projectId, targetOrganizationUserId)
		if err != nil {
			return projectManagementNotFound(err, "target project member not found")
		}
		targetOrganizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByIdWithTx(ctx, tx, targetOrganizationUserId, actor.OrganizationId)
		if err != nil {
			return projectManagementNotFound(err, "target organization member not found")
		}
		if err = l.repos.ProjectManagementProjectRepo.UpdateMemberRoleWithTx(ctx, tx, currentOwner, constants.ProjectManagementMemberRoleProjectAdmin); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementProjectRepo.UpdateMemberRoleWithTx(ctx, tx, target, constants.ProjectManagementMemberRoleProjectOwner); err != nil {
			return err
		}
		if err = l.repos.ProjectManagementProjectRepo.UpdateOwnerWithTx(ctx, tx, projectId, targetOrganizationUserId, targetOrganizationUser.FullName); err != nil {
			return err
		}
		return l.recordProjectMemberOperation(ctx, tx, actor, projectId, constants.OperationActionProjectManagementOwnerTransfer, targetOrganizationUserId)
	})
}

func (l *ProjectManagementLogic) RemoveProjectMember(ctx context.Context, actor *dto.ProjectManagementActor, projectId, organizationUserId int64) error {
	if projectId == 0 || organizationUserId == 0 {
		return errorx.New("project and member are required")
	}
	return l.repos.ProjectManagementProjectRepo.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.requireProjectManagerWithTx(ctx, tx, actor, projectId); err != nil {
			return err
		}
		member, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, projectId, organizationUserId)
		if err != nil {
			return projectManagementNotFound(err, "project member not found")
		}
		if member.Role == constants.ProjectManagementMemberRoleProjectOwner {
			return errorx.New("transfer project owner before removing member")
		}
		if err = l.repos.ProjectManagementProjectRepo.DeleteMemberWithTx(ctx, tx, member); err != nil {
			return err
		}
		return l.recordProjectMemberOperation(ctx, tx, actor, projectId, constants.OperationActionProjectManagementMemberRemove, organizationUserId)
	})
}

func (l *ProjectManagementLogic) requireProjectManagerWithTx(ctx context.Context, tx *gorm.DB, actor *dto.ProjectManagementActor, projectId int64) (*model.ProjectManagementProject, error) {
	project, err := l.repos.ProjectManagementProjectRepo.LockByOrganizationWithTx(ctx, tx, projectId, actor.OrganizationId)
	if err != nil {
		return nil, projectManagementNotFound(err, "project not found")
	}
	member, err := l.repos.ProjectManagementProjectRepo.LockMemberWithTx(ctx, tx, projectId, actor.OrganizationUserId)
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

func (l *ProjectManagementLogic) recordProjectMemberOperation(ctx context.Context, tx *gorm.DB, actor *dto.ProjectManagementActor, projectId int64, operation constants.OperationAction, targetOrganizationUserId int64) error {
	metadata, err := json.Marshal(map[string]int64{"project_id": projectId, "organization_user_id": targetOrganizationUserId})
	if err != nil {
		return errorx.WithStack(err)
	}
	return NewOperationLogic(l.BaseLogic).CreateOperationLog(ctx, tx, dto.OperationCreateDto{EntityId: projectId, EntityType: constants.OperationEntityTypeProjectManagementProject, OperatorId: actor.OrganizationUserId, OperatorName: actor.Name, OperatorType: constants.UserTypeBorrower, Operation: operation, Metadata: metadata, OrganizationId: actor.OrganizationId})
}
