package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type UserLogic struct {
	repos    *repo.AllRepos
	redisCli *redisx.Redisx
}

func NewUserLogic(repos *repo.AllRepos, redisCli *redisx.Redisx) *UserLogic {
	return &UserLogic{
		repos:    repos,
		redisCli: redisCli,
	}
}

type UserJwtToken struct {
	AccessToken        string `json:"access_token,omitempty"`
	RefreshToken       string `json:"refresh_token,omitempty"`
	AccessTokenExpire  int64  `json:"access_token_expire,omitempty"`
	RefreshTokenExpire int64  `json:"refresh_token_expire,omitempty"`
}

type JwtConfig struct {
	AccessSecret  string
	AccessExpire  int64
	RefreshSecret string
	RefreshExpire int64
}

// GenerateJwtToken generate jwt token for user
func (l *UserLogic) GenerateJwtToken(user *model.User, config JwtConfig, userType constants.UserType) (token *UserJwtToken, err error) {
	if user == nil {
		return nil, errorx.WithStack(errorx.New("User Is Empty."))
	}
	userClaim := utils.JwtUser{
		UserId:         user.Id,
		OrganizationId: user.DefaultOrganizationId,
		Username:       user.FullName,
		Role:           nil,
		UserType:       userType,
	}
	accessToken, err := utils.GenerateToken(
		userClaim,
		config.AccessSecret,
		config.AccessExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	refreshToken, err := utils.GenerateToken(
		userClaim,
		config.RefreshSecret,
		config.RefreshExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	token = &UserJwtToken{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  config.AccessExpire,
		RefreshTokenExpire: config.RefreshExpire,
	}

	return
}

// GetActiveOrganizationByUser return active organization for user
func (l *UserLogic) GetActiveOrganizationByUser(user *model.User) (currentOrganization *model.Organization, isChanged bool, err error) {
	if user == nil {
		return nil, false, errorx.New("User Is Empty")
	}
	organizations := user.Organizations
	// find current company
	currentOrganization, ok := slice.FindBy(organizations, func(index int, item *model.Organization) bool {
		return (item.Id == user.DefaultOrganizationId || user.DefaultOrganizationId == 0) && item.Status == constants.OrganizationStatusActive
	})
	// find first active company
	if !ok {
		currentOrganization, ok = slice.FindBy(organizations, func(index int, item *model.Organization) bool {
			return item.Status == constants.OrganizationStatusActive
		})
	}
	if !ok {
		// no active company
		isChanged = condition.Ternary(user.DefaultOrganizationId == 0, false, true)
		user.DefaultOrganizationId = 0
	} else if user.DefaultOrganizationId != currentOrganization.Id {
		// company changed
		isChanged = true
		user.DefaultOrganizationId = currentOrganization.Id
	}

	return currentOrganization, isChanged, err
}

// CheckAndUpdateActiveOrganizationByUser check and update active organization for user
func (l *UserLogic) CheckAndUpdateActiveOrganizationByUser(ctx context.Context, user *model.User) (currentUser *model.User, err error) {
	currentOrganization, isChanged, err := l.GetActiveOrganizationByUser(user)
	if err != nil {
		return currentUser, err
	}
	if isChanged {
		var defaultOrganizationId int64
		if currentOrganization != nil {
			defaultOrganizationId = currentOrganization.Id
		}

		user.DefaultOrganizationId = defaultOrganizationId
		err = l.repos.UserRepo.UpdateColumnsById(ctx, user.Id, map[string]interface{}{
			"default_organization_id": user.DefaultOrganizationId,
		})
		if err != nil {
			return currentUser, errorx.WithStack(err)
		}
	}
	return user, err
}

type UserLoginRequest struct {
	Email    string
	Password string
	UserType constants.UserType
}

type LoginLimiter struct {
	Rate     int
	Period   int
	RedisCli *redisx.Redisx
}

// Login
func (l *UserLogic) Login(ctx context.Context, req UserLoginRequest, limiter LoginLimiter, jwtConfig JwtConfig) (token *UserJwtToken, err error) {
	redisKey := stringx.Format(constants.UserLoginFailedTimes, req.UserType, req.Email)
	var loginFailedTimes int
	if limiter.Rate > 0 {
		loginFailedTimes, err = limiter.RedisCli.GetInt(ctx, redisKey)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		if loginFailedTimes >= limiter.Rate {
			return nil, errorx.New("Login Failed Times Exceeded, Please Try Again Later")
		}
	}

	user, err := l.repos.UserRepo.GetUserByEmail(ctx, req.Email, req.UserType)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if user == nil {
		return nil, errorx.New("User Not Found.")
	}

	// check pwd
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		// add input wrong pwd times to redis
		_, err = limiter.RedisCli.IncreaseExpire(ctx, redisKey, limiter.Period)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		return nil, errorx.WithStack(errorx.New("The Email And Password You Entered Don't Match."))
	}

	if limiter.Rate > 0 {
		err = limiter.RedisCli.Del(ctx, redisKey)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
	}

	// check status
	if user.Status == constants.UserStatusInactive {
		return nil, errorx.New("User Is Not Active")
	}
	if user.Status == constants.UserStatusUnverified {
		err = status.Error(codes.PermissionDenied, "User Is Unverified")
		return nil, err
	}

	// check status
	if user.Status == constants.UserStatusInactive {
		return nil, errorx.WithStack(errorx.New("User Is Not Active"))
	}
	if user.Status == constants.UserStatusUnverified {
		return nil, errorx.WithStack(errorx.New("User Is Unverified"))
	}

	user, err = l.CheckAndUpdateActiveOrganizationByUser(ctx, user)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	token, err = l.GenerateJwtToken(user, jwtConfig, req.UserType)
	if err != nil {
		return nil, err
	}
	return
}

type Userinfo struct {
	Id               int64
	Email            string
	FirstName        string
	LastName         string
	FullName         string
	OrganizationId   int64
	OrganizationName string
	Status           constants.UserStatus
	Organizations    []*UserOrganization
}

type UserOrganization struct {
	Id        int64
	Name      string
	Status    constants.OrganizationStatus
	IsDefault bool
}

// Info return user info
func (l *UserLogic) Info(ctx context.Context, id int64, userType constants.UserType) (userInfo *Userinfo, err error) {
	organizations := []*UserOrganization{}
	user, err := l.repos.UserRepo.GetUser(ctx, id, userType)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.WithStack(errorx.New("User Not Found"))
	}
	if len(user.Organizations) > 0 {
		err = copier.CopyWithOption(&organizations, user.Organizations, copier.Option{
			DeepCopy:      true,
			CaseSensitive: true,
		})
		if err != nil {
			return nil, errorx.Wrap(err, "Failed To Copy Organization")
		}
	}
	var currentOrganization *UserOrganization

	organizations = slice.Filter(organizations, func(i int, organization *UserOrganization) bool {
		if organization.Status != constants.OrganizationStatusActive {
			return false
		}
		relation, err := l.repos.OrganizationAndUserRepo.GetRelation(ctx, organization.Id, user.Id)
		if err != nil || relation == nil || relation.Status != constants.OrganizationAndUserStatusActive {
			return false
		}
		if organization.Id == user.DefaultOrganizationId {
			currentOrganization = organizations[i]
			organization.IsDefault = true
		}

		return true
	})

	var organizationId int64
	var organizationName string

	if currentOrganization == nil && len(organizations) > 0 {
		currentOrganization = organizations[0]
	}

	if currentOrganization != nil {
		organizationId = user.DefaultOrganizationId
		organizationName = currentOrganization.Name
	}
	userInfo = &Userinfo{
		Id:               user.Id,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		FullName:         user.FullName,
		OrganizationId:   organizationId,
		OrganizationName: organizationName,
		Status:           user.Status,
		Organizations:    organizations,
	}
	return
}

// SwitchOrganization switch user default organization
func (l *UserLogic) SwitchOrganization(ctx context.Context, userId int64, organizationId int64, userType constants.UserType) (user *model.User, err error) {
	user, err = l.repos.UserRepo.GetUser(ctx, userId, userType)
	if err != nil {
		return nil, err
	}

	organization, ok := slice.FindBy(user.Organizations, func(i int, organization *model.Organization) bool {
		return organization.Id == organizationId
	})
	if !ok {
		return nil, errorx.New("User Is Not Member Of This Organization")
	}

	if organization.Status != constants.OrganizationStatusActive {
		return nil, errorx.New("Organization Is Not Active")
	}

	relation, _ := l.repos.OrganizationAndUserRepo.GetRelation(ctx, organizationId, user.Id)
	if relation == nil || relation.Status != constants.OrganizationAndUserStatusActive {
		return nil, errorx.New("User Is Not Member Of This Organization")
	}

	user.DefaultOrganizationId = organizationId
	err = l.repos.UserRepo.UpdateColumnsById(ctx, user.Id, map[string]interface{}{
		"default_organization_id": user.DefaultOrganizationId,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	return
}

// AddOrganizationUser add user to organization
func (l *UserLogic) AddOrganizationUser(ctx context.Context, user *model.User, organizationId int64, roleIds []int64) (newUser *model.User, err error) {
	if user == nil {
		return nil, errorx.New("User Is Empty")
	}
	organizationRelations := user.OrganizationAndUsers
	if len(organizationRelations) == 0 {
		return nil, errorx.New("User Must Be Member Of An Organization")
	}
	organizationRelation := organizationRelations[0]

	err = l.repos.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		isExisted, innerErr := l.repos.UserRepo.CheckOrganizationUserExistWithTx(ctx, tx, organizationId, user.Email)
		if innerErr != nil {
			return innerErr
		}
		if isExisted {
			return errorx.New("User Is Already Member Of This Organization")
		}
		roles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, organizationId, roleIds)
		if innerErr != nil {
			return innerErr
		}
		// get existed user
		existedUser, _ := l.repos.UserRepo.GetUserByEmailWithTx(ctx, tx, user.Email, user.Type)
		if existedUser != nil && existedUser.Id > 0 {
			user.RemoveRelationships()
			user = existedUser
		}
		user.Roles = roles
		// create user
		if user.Id == 0 {
			newUser, innerErr = l.repos.UserRepo.CreateWithTx(ctx, tx, user, true)
		} else {
			innerErr = l.repos.UserRepo.UpdateRolesWithTx(ctx, tx, user, roles, nil)
			newUser = user
		}
		if innerErr != nil {
			return innerErr
		}

		organizationRelation.UserId = newUser.Id
		organizationRelation.OrganizationId = organizationId
		_, innerErr = l.repos.OrganizationAndUserRepo.CreateWithTx(ctx, tx, organizationRelation, false)
		if innerErr != nil {
			return innerErr
		}
		return innerErr
	})
	// todo update user's permissions
	return
}

// UpdateOrganizationUser update user's info in organization
func (l *UserLogic) UpdateOrganizationUser(ctx context.Context, id int64, user *model.User, organizationId int64, roleIds []int64) (newUser *model.User, err error) {
	if user == nil {
		return nil, errorx.New("User Is Empty")
	}
	oldRoles := user.Roles

	err = l.repos.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		isExisted, innerErr := l.repos.UserRepo.CheckOrganizationUserExistWithTx(ctx, tx, organizationId, user.Email)
		if innerErr != nil {
			return innerErr
		}
		if !isExisted {
			return errorx.New("User Is Not Member Of This Organization")
		}
		currentRoles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, organizationId, roleIds)
		if innerErr != nil {
			return innerErr
		}

		newRoles := slice.DifferenceWith(currentRoles, oldRoles, func(item1, item2 *model.Role) bool {
			return item1.Id == item2.Id
		})
		removedRoles := slice.DifferenceWith(oldRoles, currentRoles, func(item1, item2 *model.Role) bool {
			return item1.Id == item2.Id
		})
		innerErr = l.repos.UserRepo.UpdateRolesWithTx(ctx, tx, user, newRoles, removedRoles)
		if innerErr != nil {
			return innerErr
		}
		// update user
		user.RemoveRelationships()
		newUser, innerErr = l.repos.UserRepo.UpdateOneByIdWithTx(ctx, tx, id, user)
		if innerErr != nil {
			return innerErr
		}
		// todo check user for Lender/Borrower

		return innerErr
	})
	// todo update user's permissions
	return
}

// DeleteOrganizationUser delete user from organization
func (l *UserLogic) DeleteOrganizationUser(ctx context.Context, id int64, organizationId int64, userType constants.UserType) (err error) {
	user, err := l.repos.UserRepo.GetOrganizationUser(ctx, id, organizationId, userType)
	if err != nil {
		return err
	}
	err = l.repos.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		organization := user.Organizations[0]
		roles := user.Roles
		// remove user and role relations
		innerErr = l.repos.UserRepo.UpdateRolesWithTx(ctx, tx, user, nil, roles)
		if innerErr != nil {
			return innerErr
		}
		// remove user and organization relations
		if organization.Id == organizationId {
			innerErr = l.repos.OrganizationAndUserRepo.DeleteOrganizationAndUserRelationWithTx(ctx, tx, organization.Id, user.Id)
			if innerErr != nil {
				return innerErr
			}
		}

		return innerErr
	})
	if err != nil {
		return err
	}
	// todo update user's permissions

	return err
}

// SwitchOrganizationUserStatus switch user's status in organization
func (l *UserLogic) SwitchOrganizationUserStatus(ctx context.Context, id int64, organizationId int64, status constants.OrganizationAndUserStatus, userType constants.UserType) (err error) {
	err = status.Validate()
	if err != nil {
		return err
	}
	user, err := l.repos.UserRepo.GetOrganizationUser(ctx, id, organizationId, userType)
	if err != nil {
		return err
	}
	err = l.repos.OrganizationAndUserRepo.UpdateOrganizationAndUserRelationStatus(ctx, organizationId, user.Id, status)
	if err != nil {
		return err
	}
	// todo update user's permissions

	return nil
}

// RegisterWithOrganization register user with organization
func (l *UserLogic) RegisterWithOrganization(ctx context.Context, user *model.User, organization *model.Organization, userType constants.UserType) (newUser *model.User, err error) {
	// check user exist
	if user == nil {
		return nil, errorx.New("User Is Empty")
	}
	isExisted, err := l.repos.UserRepo.IsExistedByColumns(ctx, map[string]interface{}{
		"email": user.Email,
		"type":  userType,
	})
	if err != nil {
		return nil, err
	}
	if isExisted {
		return nil, errorx.New("User Is Already Existed")
	}
	// check organization exist
	if organization == nil {
		return nil, errorx.New("Organization Is Empty")
	}
	isExisted, err = l.repos.OrganizationRepo.IsExistedByColumns(ctx, map[string]interface{}{
		"name": organization.Name,
		"type": userType,
	})
	if err != nil {
		return nil, err
	}
	if isExisted {
		return nil, errorx.New("Organization Is Already Existed")
	}

	err = l.repos.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// create  organization
		organization.Status = constants.OrganizationStatusInactive
		organization, innerErr = l.repos.OrganizationRepo.CreateWithTx(ctx, tx, organization, false)
		if innerErr != nil {
			return innerErr
		}
		user.DefaultOrganizationId = organization.Id
		// create user
		roles := []*model.Role{
			{
				Name:           constants.AdminRoleName,
				Type:           userType,
				OrganizationId: organization.Id,
			},
			{
				Name:           "general",
				Type:           userType,
				OrganizationId: organization.Id,
			},
		}
		user.RemoveRelationships()
		user.Roles = roles
		user, innerErr = l.repos.UserRepo.CreateWithTx(ctx, tx, user, true)
		if innerErr != nil {
			return innerErr
		}
		// update organization's owner
		organization.OwnerId = user.Id
		innerErr = l.repos.OrganizationRepo.UpdateOwnerIdWithTx(ctx, tx, organization.Id, user.Id)
		if innerErr != nil {
			return innerErr
		}

		// create organization and user relation
		_, innerErr = l.repos.OrganizationAndUserRepo.CreateWithTx(ctx, tx, &model.OrganizationAndUser{
			OrganizationId: organization.Id,
			UserId:         user.Id,
			Status:         constants.OrganizationAndUserStatusActive,
		}, false)
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

// VerifyUserEmail verity user email
func (l *UserLogic) VerifyUserEmail(ctx context.Context, email string, code string, userType constants.UserType) (user *model.User, err error) {
	var token string
	if email == "" || code == "" {
		return nil, errorx.New("Email And Code Are Required")
	}
	// check code
	verifyCodeKey := stringx.Format(constants.UserRegisterVerifyCode, userType, email)
	token, err = l.redisCli.Get(ctx, verifyCodeKey)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if token != code {
		return nil, errorx.New("Invalid Or Expired Code")
	}
	err = l.repos.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// get user
		user, innerErr = l.repos.UserRepo.GetUserByEmailWithTx(ctx, tx, email, userType)
		if innerErr != nil {
			return innerErr
		}
		if user.Status != constants.UserStatusUnverified {
			return errorx.New("Email Is Already Verified")
		}
		// update user status
		innerErr = l.repos.UserRepo.UpdateColumnsByIdWithTx(ctx, tx, user.Id, map[string]interface{}{
			"status": constants.UserStatusActive,
		})
		if innerErr != nil {
			return innerErr
		}

		// update user's organization status
		innerErr = l.repos.OrganizationRepo.UpdateOwnerOrganizationStatusWithTx(ctx, tx, user.Id, constants.OrganizationStatusActive)
		if innerErr != nil {
			return innerErr
		}
		return innerErr
	})

	return
}
