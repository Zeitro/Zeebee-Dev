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

type Role struct {
	model.Role
	Resources []*ResourceWithPivot
}

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
		return nil, errorx.New("Role Already Exists")
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

	result = l.transToBo(role)

	return result, nil
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

	resourceIds := slice.Map[*model.RoleAndResource, int64](data.RoleAndResources, func(i int, item *model.RoleAndResource) int64 {
		return item.ResourceId
	})
	err = l.repos.RoleRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		var resources []*model.Resource
		var validResourceIds []int64
		// create role
		role, innerErr = l.repos.RoleRepo.UpdateOneWithTx(ctx, tx, data)
		if innerErr != nil {
			return innerErr
		}

		// save role and resource relationship
		resources, innerErr = l.repos.ResourceRepo.GetResourcesWithTx(ctx, tx, resourceIds, data.Type)
		if innerErr != nil {
			return innerErr
		}
		validResourceIds = slice.Map[*model.Resource, int64](resources, func(i int, item *model.Resource) int64 {
			return item.Id
		})
		rolesAndResources := slice.FilterMap[*model.RoleAndResource, *model.RoleAndResource](data.RoleAndResources, func(i int, item *model.RoleAndResource) (*model.RoleAndResource, bool) {
			if slice.Contain(validResourceIds, item.ResourceId) {
				data.RoleAndResources[i].RoleId = role.Id
				return data.RoleAndResources[i], true
			}
			return nil, false
		})
		oldRelations, innerErr := l.repos.RoleAndResourceRepo.GetRoleRelationsWithTx(ctx, tx, role.Id)
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
	// todo update user's permissions

	return role, nil
}

// GetRole returns a role by id
func (l *RoleLogic) GetRole(ctx context.Context, id int64, organizationId int64, roleType constants.UserType) (result *Role, err error) {
	role, err := l.repos.RoleRepo.GetRole(ctx, id, organizationId, roleType)
	if err != nil {
		return nil, err
	}
	result = &Role{}
	err = copier.CopyWithOption(result, role, copier.Option{
		DeepCopy: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	slice.ForEach(result.Resources, func(index int, item *ResourceWithPivot) {
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
		return errorx.New("Cannot Delete Admin Role")
	}

	// todo remove role's users

	// todo remove user's permissions

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
	if err != nil {
		return err
	}

	// todo update user's permissions

	return nil
}

// All returns all roles by type
func (l *RoleLogic) All(ctx context.Context, organizationId int64, roleType constants.UserType) (result []*bo.RoleBo, err error) {
	roles, err := l.repos.RoleRepo.All(ctx, organizationId, roleType)
	if err != nil {
		return nil, err
	}
	result = l.batchTransToBo(roles)
	return result, err
}

// transToBo
func (l *RoleLogic) transToBo(src *model.Role) (result *bo.RoleBo) {
	if src == nil {
		return nil
	}
	result = &bo.RoleBo{}
	_ = copier.CopyWithOption(result, src, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	return result
}

// batchTransToBo batch trans to bo
func (l *RoleLogic) batchTransToBo(src []*model.Role) (result []*bo.RoleBo) {
	result = make([]*bo.RoleBo, 0, len(src))
	if len(src) == 0 {
		return result
	}
	for i := range src {
		result = append(result, l.transToBo(src[i]))
	}
	return result
}
