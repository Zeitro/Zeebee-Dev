package logic

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type RoleLogic struct {
	repos *repo.AllRepos
}

func NewRoleLogic(repos *repo.AllRepos) *RoleLogic {
	return &RoleLogic{
		repos: repos,
	}
}

func (l *RoleLogic) checkRoleName(roleName string) (err error) {
	if roleName == constants.AdminRoleName {
		return errorx.New("Role Name is Reserved")
	}
	return err
}

// Page returns paginated roles
func (l *RoleLogic) Page(ctx context.Context, req *repo.RolePageReq) (result gormx.PageResponse[*model.Role], err error) {
	result, err = l.repos.RoleRepo.Page(ctx, req)
	return result, err
}

// Create creates a new role
func (l *RoleLogic) Create(ctx context.Context, data *model.Role) (role *model.Role, err error) {
	if err := l.checkRoleName(data.Name); err != nil {
		return nil, err
	}

	// check existed
	isExisted, err := l.repos.RoleRepo.IsExisted(ctx, data)
	if err != nil {
		return nil, err
	}
	if isExisted {
		return nil, errorx.New("Role Already Exists")
	}

	role, err = l.repos.RoleRepo.Create(ctx, data)

	return role, nil
}

// Update updates a role
func (l *RoleLogic) Update(ctx context.Context, id int64, data *model.Role) (role *model.Role, err error) {
	if err := l.checkRoleName(data.Name); err != nil {
		return nil, err
	}
	// First get the existing record
	result, err := l.repos.RoleRepo.GetRole(ctx, id, data.OrganizationId, data.Type)
	if err != nil {
		return nil, err
	}

	// Update the fields
	err = copier.CopyWithOption(result, data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, err
	}

	// check existed
	isExisted, err := l.repos.RoleRepo.IsExisted(ctx, result)
	if err != nil {
		return nil, err
	}
	if isExisted {
		return nil, errorx.New("Role Already Exists")
	}

	// Save the updated record
	role, err = l.repos.RoleRepo.UpdateOneById(ctx, id, result)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetRole returns a role by id
func (l *RoleLogic) GetRole(ctx context.Context, id int64, organizationId int64, roleType constants.UserType) (role *model.Role, err error) {
	role, err = l.repos.RoleRepo.GetRole(ctx, id, organizationId, roleType)
	return role, err
}

// Delete deletes a role
func (l *RoleLogic) Delete(ctx context.Context, id int64, organizationId int64, roleType constants.UserType) (err error) {
	role, err := l.repos.RoleRepo.GetRole(ctx, id, organizationId, roleType)
	if err != nil {
		return err
	}
	if role.Name == constants.AdminRoleName {
		return errorx.New("Cannot Delete Admin Role")
	}

	// todo remove role's users

	// todo remove role's permissions

	_, err = l.repos.RoleRepo.DeleteByID(ctx, id)
	return err
}

// UpdateStatus updates a role's status
func (l *RoleLogic) UpdateStatus(ctx context.Context, id int64, organizationId int64, roleType constants.UserType, status constants.CommonStatus) (err error) {
	// First get the existing record
	result, err := l.repos.RoleRepo.GetRole(ctx, id, organizationId, roleType)
	if err != nil {
		return err
	}
	if result.Name == constants.AdminRoleName {
		return errorx.New("Cannot Update Admin Role")
	}

	result.Status = status
	// Save the updated record
	_, err = l.repos.RoleRepo.UpdateOneById(ctx, id, result)
	return err
}

// All returns all roles by type
func (l *RoleLogic) All(ctx context.Context, organizationId int64, roleType constants.UserType) (result []*model.Role, err error) {
	result, err = l.repos.RoleRepo.All(ctx, organizationId, roleType)
	return result, err
}
