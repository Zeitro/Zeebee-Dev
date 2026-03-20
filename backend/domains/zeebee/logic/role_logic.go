package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// type Role struct {
// 	model.Role
// 	Resources []*ResourceWithPivot
// }

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
		return errorx.New("Role name is reserved")
	}
	return err
}

// Page returns paginated roles
func (l *RoleLogic) Page(ctx context.Context, req dto.RolePageRequestDto) (result gormx.PageResponse[*model.Role], err error) {
	result, err = l.repos.RoleRepo.Page(ctx, req)
	return result, err
}

// Create creates a new role
func (l *RoleLogic) Create(ctx context.Context, input dto.RoleCreateRequestDto) (result *bo.RoleBo, err error) {
	if err := l.checkRoleName(input.Name); err != nil {
		return nil, err
	}

	role := &model.Role{}
	err = copier.Copy(role, &input)
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}

	// check existed
	isExisted, err := l.repos.RoleRepo.IsExisted(ctx, role)
	if err != nil {
		return nil, err
	}
	if isExisted {
		return nil, errorx.New("Role already exists")
	}

	// resourceIds := input.ResourceIds
	resourceIds := slice.Map[dto.RoleAndResourceDto, int64](input.RoleAndResources, func(i int, item dto.RoleAndResourceDto) int64 {
		return item.ResourceId
	})

	err = l.repos.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		var resources []*model.Resource
		var validResourceIds []int64
		// create role
		role, innerErr = l.repos.RoleRepo.CreateWithTx(ctx, tx, role, false)
		if innerErr != nil {
			return innerErr
		}

		// create role and resource relationship
		resources, innerErr = l.repos.ResourceRepo.GetResourcesWithTx(ctx, tx, resourceIds, input.Type)
		if innerErr != nil {
			return innerErr
		}
		validResourceIds = slice.Map[*model.Resource, int64](resources, func(i int, item *model.Resource) int64 {
			return item.Id
		})
		rolesAndResources := slice.FilterMap[dto.RoleAndResourceDto, *model.RoleAndResource](input.RoleAndResources, func(i int, item dto.RoleAndResourceDto) (*model.RoleAndResource, bool) {
			if slice.Contain(validResourceIds, item.ResourceId) {
				role.RoleAndResources[i].RoleId = role.Id
				return role.RoleAndResources[i], true
			}
			return nil, false
		})
		if len(rolesAndResources) > 0 {
			_, innerErr = l.repos.RoleAndResourceRepo.BatchCreateWithTx(ctx, tx, rolesAndResources, false)
			if innerErr != nil {
				return innerErr
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	result, err = copierx.Trans[model.Role, bo.RoleBo](role)

	return result, err
}

// Update updates a role
func (l *RoleLogic) Update(ctx context.Context, id int64, input dto.RoleUpdateRequestDto) (result *bo.RoleBo, err error) {
	if err := l.checkRoleName(input.Name); err != nil {
		return nil, err
	}
	// First get the existing record
	roleModel, err := l.repos.RoleRepo.GetRole(ctx, id, input.OrganizationId, input.Type)
	if err != nil {
		return nil, err
	}

	roleModel, err = copierx.Trans[dto.RoleUpdateRequestDto, model.Role](&input)

	// err = copier.CopyWithOption(result, data, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, err
	}

	// check existed
	isExisted, err := l.repos.RoleRepo.IsExisted(ctx, roleModel)
	if err != nil {
		return nil, err
	}
	if isExisted {
		return nil, errorx.New("Role already exists")
	}

	resourceIds := slice.Map[*model.RoleAndResource, int64](roleModel.RoleAndResources, func(i int, item *model.RoleAndResource) int64 {
		return item.ResourceId
	})
	err = l.repos.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		var resources []*model.Resource
		var validResourceIds []int64
		// update role
		roleModel, innerErr = l.repos.RoleRepo.UpdateOneWithTx(ctx, tx, roleModel)
		if innerErr != nil {
			return innerErr
		}

		// save role and resource relationship
		resources, innerErr = l.repos.ResourceRepo.GetResourcesWithTx(ctx, tx, resourceIds, roleModel.Type)
		if innerErr != nil {
			return innerErr
		}
		validResourceIds = slice.Map[*model.Resource, int64](resources, func(i int, item *model.Resource) int64 {
			return item.Id
		})
		rolesAndResources := slice.FilterMap[*model.RoleAndResource, *model.RoleAndResource](roleModel.RoleAndResources, func(i int, item *model.RoleAndResource) (*model.RoleAndResource, bool) {
			if slice.Contain(validResourceIds, item.ResourceId) {
				roleModel.RoleAndResources[i].RoleId = roleModel.Id
				return roleModel.RoleAndResources[i], true
			}
			return nil, false
		})
		oldRelations, innerErr := l.repos.RoleAndResourceRepo.GetRoleRelationsWithTx(ctx, tx, roleModel.Id)
		if innerErr != nil {
			return innerErr
		}

		updateRelations := []*model.RoleAndResource{}
		newRelations := []*model.RoleAndResource{}
		slice.ForEach(rolesAndResources, func(index int, item *model.RoleAndResource) {
			oldRelation, ok := slice.FindBy(oldRelations, func(subIndex int, subItem *model.RoleAndResource) bool {
				return subItem.ResourceId == item.ResourceId
			})
			if ok {
				rolesAndResources[index].Id = oldRelation.Id
				updateRelations = append(updateRelations, rolesAndResources[index])
			} else {
				newRelations = append(newRelations, rolesAndResources[index])
			}
		})
		removedRelations := slice.FilterMap[*model.RoleAndResource, *model.RoleAndResource](oldRelations, func(index int, item *model.RoleAndResource) (*model.RoleAndResource, bool) {
			_, ok := slice.FindBy(rolesAndResources, func(subIndex int, subItem *model.RoleAndResource) bool {
				return subItem.ResourceId == item.ResourceId
			})
			if ok {
				return nil, false
			}
			return oldRelations[index], true
		})
		if len(newRelations) > 0 {
			_, innerErr = l.repos.RoleAndResourceRepo.BatchCreateWithTx(ctx, tx, newRelations, false)
			if innerErr != nil {
				return innerErr
			}
		}
		if len(updateRelations) > 0 {
			for _, updatedRelation := range updateRelations {
				_, innerErr = l.repos.RoleAndResourceRepo.UpdateOneWithTx(ctx, tx, updatedRelation)
				if innerErr != nil {
					return innerErr
				}
			}
		}
		if len(removedRelations) > 0 {
			removedIds := slice.Map[*model.RoleAndResource, int64](removedRelations, func(i int, item *model.RoleAndResource) int64 {
				return item.Id
			})
			_, innerErr = l.repos.RoleAndResourceRepo.DeleteByIDsWithTx(ctx, tx, removedIds)
			if innerErr != nil {
				return innerErr
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	// TODO update user's permissions

	result, err = copierx.Trans[model.Role, bo.RoleBo](roleModel)
	return result, err
}

// GetRole returns a role by id
func (l *RoleLogic) GetRole(ctx context.Context, id int64, organizationId int64, roleType constants.UserType) (result *bo.RoleBo, err error) {
	role, err := l.repos.RoleRepo.GetRole(ctx, id, organizationId, roleType)
	if err != nil {
		return nil, err
	}
	result, err = copierx.Trans[model.Role, bo.RoleBo](role)
	if err != nil {
		return nil, err
	}
	slice.ForEach(result.Resources, func(index int, item *bo.ResourceBo) {
		relation, err := l.repos.RoleAndResourceRepo.GetRelation(ctx, role.Id, item.Id)
		if err != nil {
			return
		}
		item.PivotRoleDataScope = relation.DataScope
	})

	return result, nil
}

// Delete deletes a role
func (l *RoleLogic) Delete(ctx context.Context, id int64, organizationId int64, roleType constants.UserType) (err error) {
	role, err := l.repos.RoleRepo.GetRole(ctx, id, organizationId, roleType)
	if err != nil {
		return err
	}
	if role.Name == constants.AdminRoleName {
		return errorx.New("Cannot delete admin role")
	}

	// TODO remove role's users

	// TODO remove user's permissions

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
		return errorx.New("Cannot update admin role")
	}

	result.Status = status
	// Save the updated record
	_, err = l.repos.RoleRepo.UpdateOneById(ctx, id, result)
	if err != nil {
		return err
	}

	// TODO update user's permissions

	return nil
}

// All returns all roles by type
func (l *RoleLogic) All(ctx context.Context, organizationId int64, roleType constants.UserType) (result []*bo.RoleBo, err error) {
	roles, err := l.repos.RoleRepo.All(ctx, organizationId, roleType)
	if err != nil {
		return nil, err
	}
	results, err := copierx.Trans[[]*model.Role, []*bo.RoleBo](&roles)
	return *results, err
}
