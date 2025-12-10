package logic

import (
	"context"
	"regexp"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type UserLogic struct {
	repos     *repo.AllRepos
	redisCli  *redisx.Redisx
	jwtConfig JwtConfig
}

func NewUserLogic(repos *repo.AllRepos, redisCli *redisx.Redisx, jwtConfig JwtConfig) *UserLogic {
	return &UserLogic{
		repos:     repos,
		redisCli:  redisCli,
		jwtConfig: jwtConfig,
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

// GetUserClaimFromCtx get user claim from context
func (l *UserLogic) GetUserClaimFromCtx(ctx context.Context) (claim utils.JwtUser, err error) {
	// get from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = status.Error(codes.Unauthenticated, "Not Login")
		return claim, errorx.WithStack(err)
	}
	tokens := md[constants.JwtTokenCtxKey]
	if len(tokens) == 0 {
		err = status.Error(codes.Unauthenticated, "Not Login")
		return claim, errorx.WithStack(err)
	}
	token := tokens[0]
	jwtClaim, err := utils.ParseToken(token, l.jwtConfig.AccessSecret)
	if err != nil {
		err = status.Error(codes.InvalidArgument, err.Error())
		return claim, errorx.WithStack(err)
	}
	return jwtClaim.JwtUser, nil
}

// GenerateJwtToken generate jwt token for user
func (l *UserLogic) GenerateJwtToken(user *model.User, orgUser *model.OrganizationUser, userType constants.UserType) (token *UserJwtToken, err error) {
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
	if orgUser != nil {
		userClaim.OrganizationId = orgUser.OrganizationId
		userClaim.OrganizationUserId = orgUser.Id
		userClaim.Username = orgUser.FullName
	}

	accessToken, err := utils.GenerateToken(
		userClaim,
		l.jwtConfig.AccessSecret,
		l.jwtConfig.AccessExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	refreshToken, err := utils.GenerateToken(
		userClaim,
		l.jwtConfig.RefreshSecret,
		l.jwtConfig.RefreshExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	token = &UserJwtToken{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  l.jwtConfig.AccessExpire,
		RefreshTokenExpire: l.jwtConfig.RefreshExpire,
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
func (l *UserLogic) CheckAndUpdateActiveOrganizationByUser(ctx context.Context, user *model.User) (currentUser *model.User, orgUser *model.OrganizationUser, err error) {
	currentOrganization, isChanged, err := l.GetActiveOrganizationByUser(user)
	if err != nil {
		return nil, nil, err
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
			return nil, nil, errorx.WithStack(err)
		}
	}
	if currentOrganization != nil {
		orgUser, err = l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, currentOrganization.Id, user.Id)
		if err != nil {
			return nil, nil, err
		}
	}

	return user, orgUser, err
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
func (l *UserLogic) Login(ctx context.Context, req UserLoginRequest, limiter LoginLimiter) (token *UserJwtToken, err error) {
	redisKey := stringx.Format(constants.UserLoginFailedTimes, req.UserType, req.Email)
	var loginFailedTimes int
	var orgUser *model.OrganizationUser
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

	user, orgUser, err = l.CheckAndUpdateActiveOrganizationByUser(ctx, user)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	token, err = l.GenerateJwtToken(user, orgUser, req.UserType)
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
	Status           constants.OrganizationUserStatus
	Organizations    []*UserOrganization
	IsOwner          bool
	IsAdmin          bool
}

type UserOrganization struct {
	Id        int64
	Name      string
	Status    constants.OrganizationStatus
	IsDefault bool
}

// Info return user info
func (l *UserLogic) Info(ctx context.Context, id int64, organizationId int64, userType constants.UserType) (userInfo *Userinfo, err error) {
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
	var organizationUser *model.OrganizationUser

	organizations = slice.Filter(organizations, func(i int, organization *UserOrganization) bool {
		if organization.Status != constants.OrganizationStatusActive {
			return false
		}
		orgUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, organization.Id, user.Id)
		if err != nil || orgUser == nil || orgUser.Status != constants.OrganizationUserStatusActive {
			return false
		}
		if organization.Id == user.DefaultOrganizationId {
			currentOrganization = organizations[i]
			organization.IsDefault = true
		}
		return true
	})

	organizationUser, ok := slice.FindBy(user.OrganizationUsers, func(index int, item *model.OrganizationUser) bool {
		return item.OrganizationId == organizationId
	})
	if !ok {
		err = errorx.New("Organization User Is Not Found")
		return nil, err
	}

	// var organizationId int64
	var organizationName string
	var isOwner bool
	var isAdmin bool

	if currentOrganization == nil && len(organizations) > 0 {
		currentOrganization = organizations[0]
	}

	if currentOrganization != nil {
		organizationId = user.DefaultOrganizationId
		organizationName = currentOrganization.Name

		roles, err := l.repos.RoleRepo.GetOrganizationUserRoles(ctx, organizationUser.Id)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		isAdmin = slice.ContainBy(roles, func(role *model.Role) bool {
			return role.Name == constants.RoleAdmin
		})
		isOwner = slice.ContainBy(user.Organizations, func(item *model.Organization) bool {
			return item.Id == currentOrganization.Id && item.OwnerId == user.Id
		})
	}

	userInfo = &Userinfo{
		Id:               user.Id,
		Email:            user.Email,
		FirstName:        organizationUser.FirstName,
		LastName:         organizationUser.LastName,
		FullName:         organizationUser.FullName,
		OrganizationId:   organizationId,
		OrganizationName: organizationName,
		Status:           organizationUser.Status,
		Organizations:    organizations,
		IsOwner:          isOwner,
		IsAdmin:          isAdmin,
	}
	return
}

// SwitchOrganization switch user default organization
func (l *UserLogic) SwitchOrganization(ctx context.Context, userId int64, organizationId int64, userType constants.UserType) (user *model.User, organizationUser *model.OrganizationUser, err error) {
	user, err = l.repos.UserRepo.GetUser(ctx, userId, userType)
	if err != nil {
		return nil, nil, err
	}

	organization, ok := slice.FindBy(user.Organizations, func(i int, organization *model.Organization) bool {
		return organization.Id == organizationId
	})
	if !ok {
		return nil, nil, errorx.New("User Is Not Member Of This Organization")
	}

	if organization.Status != constants.OrganizationStatusActive {
		return nil, nil, errorx.New("Organization Is Not Active")
	}

	organizationUser, _ = l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, organizationId, user.Id)
	if organizationUser == nil || organizationUser.Status != constants.OrganizationUserStatusActive {
		return nil, nil, errorx.New("User Is Not Member Of This Organization")
	}

	user.DefaultOrganizationId = organizationId
	err = l.repos.UserRepo.UpdateColumnsById(ctx, user.Id, map[string]interface{}{
		"default_organization_id": user.DefaultOrganizationId,
	})
	if err != nil {
		return nil, nil, errorx.WithStack(err)
	}
	return user, organizationUser, nil
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
		// user.Roles = roles
		user, innerErr = l.repos.UserRepo.CreateWithTx(ctx, tx, user, false)
		if innerErr != nil {
			return innerErr
		}
		// update organization's owner
		organization.OwnerId = user.Id
		innerErr = l.repos.OrganizationRepo.UpdateOwnerIdWithTx(ctx, tx, organization.Id, user.Id)
		if innerErr != nil {
			return innerErr
		}

		// create organization user
		orgUser := &model.OrganizationUser{
			OrganizationId: organization.Id,
			UserId:         user.Id,
			FirstName:      user.FirstName,
			MidName:        user.MidName,
			LastName:       user.LastName,
			Status:         constants.OrganizationUserStatusActive,
			FullName:       user.FullName,
			Roles:          roles,
		}
		// create organization and user relation
		_, innerErr = l.repos.OrganizationUserRepo.CreateWithTx(ctx, tx, orgUser, true)
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
	// Check if code exists (expired or invalid)
	if len(token) == 0 || token != code {
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

// ChangePassword change password
func (l *UserLogic) ChangePassword(ctx context.Context, input dto.UserChangePasswordRequestDto) (err error) {
	user, err := l.repos.UserRepo.GetUser(ctx, input.UserId, input.UserType)
	if err != nil {
		return err
	}
	if !l.validatePassword(input.NewPassword) {
		return errorx.New("New Password Is Invalid")
	}

	oldPassword := input.OldPassword
	if !utils.CheckPassword(oldPassword, user.PasswordHash) {
		return errorx.New("Old Password Is Incorrect")
	}
	newPasswordHash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return errorx.WithStack(err)
	}
	user.PasswordHash = newPasswordHash
	user.RemoveRelationships()

	_, err = l.repos.UserRepo.UpdateOneById(ctx, input.UserId, user)
	return err
}

func (l *UserLogic) validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	return hasLower && hasUpper && hasNumber
}
