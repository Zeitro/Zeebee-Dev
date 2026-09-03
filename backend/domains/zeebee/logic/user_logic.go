package logic

import (
	"context"
	"errors"
	"regexp"
	"time"

	"zeebee-backend/common/timex"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
)

// UserLogic provides user business logic
type UserLogic struct {
	*BaseLogic
	// jwtConfig    dto.JwtConfig
	authTopicMap map[constants.UserType]string
}

// NewUserLogic creates a new user logic instance
func NewUserLogic(baseLogic *BaseLogic, authTopicMap map[constants.UserType]string) *UserLogic {
	return &UserLogic{
		BaseLogic: baseLogic,
		// jwtConfig:    jwtConfig,
		authTopicMap: authTopicMap,
	}
}

// GetUserClaimFromCtx parses the caller's JWT from gRPC metadata and returns
// the claim. The shared jwtConfig parses the token (single unified key). The
// platform user_type MUST be present in metadata and is checked against the
// type embedded in the token — an absent or mismatched user_type (a token
// minted for one platform and replayed against another) is rejected as an
// invalid token.
func (l *UserLogic) GetUserClaimFromCtx(ctx context.Context) (claim utils.JwtUser, err error) {
	// Get from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = status.Error(codes.Unauthenticated, "not login")
		return claim, errorx.WithStack(err)
	}
	tokens := md[utils.JwtTokenCtxKey]
	if len(tokens) == 0 {
		err = status.Error(codes.Unauthenticated, "not login")
		return claim, errorx.WithStack(err)
	}
	token := tokens[0]
	jwtConf := l.jwtConfig
	jwtClaim, err := utils.ParseToken(token, jwtConf.AccessSecret)
	if err != nil {
		err = status.Error(codes.Unauthenticated, "invalid token")
		return claim, errorx.WithStack(err)
	}
	// Defense-in-depth: verify the platform type declared in metadata matches
	// the type embedded in the token. The api auth middleware already enforces
	// this at the edge; re-checking here catches a token replayed across
	// platforms via a path that bypasses the middleware.
	ctxUserTypes, ok := md[utils.UserTypeCtxKey]
	if !ok || len(ctxUserTypes) == 0 {
		err = status.Error(codes.Unauthenticated, "invalid token")
		return claim, errorx.WithStack(err)
	}
	if ctxUserTypes[0] != string(jwtClaim.UserType) {
		err = status.Error(codes.Unauthenticated, "invalid token")
		return claim, errorx.WithStack(err)
	}
	return jwtClaim.JwtUser, nil
}

// GenerateJwtToken generate jwt token for user
func (l *UserLogic) GenerateJwtToken(user *model.User, organizationUser *model.OrganizationUser) (token *dto.UserJwtTokenDto, err error) {
	if user == nil {
		return nil, errorx.WithStack(errorx.New("user is empty"))
	}
	jwtConfig := l.jwtConfig
	userClaim := utils.JwtUser{
		UserId:         user.Id,
		OrganizationId: user.DefaultOrganizationId,
		Username:       user.FullName,
		Role:           nil,
		UserType:       user.Type,
	}
	if organizationUser != nil {
		userClaim.OrganizationId = organizationUser.OrganizationId
		userClaim.OrganizationUserId = organizationUser.Id
		userClaim.Username = organizationUser.FullName
	}

	accessToken, err := utils.GenerateToken(
		userClaim,
		jwtConfig.AccessSecret,
		jwtConfig.AccessExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	refreshToken, err := utils.GenerateToken(
		userClaim,
		jwtConfig.RefreshSecret,
		jwtConfig.RefreshExpire,
	)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	token = &dto.UserJwtTokenDto{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  jwtConfig.AccessExpire,
		RefreshTokenExpire: jwtConfig.RefreshExpire,
	}
	return
}

func (l *UserLogic) GenerateJwtTokenById(ctx context.Context, userId int64, organizationId int64, userType constants.UserType) (token *dto.UserJwtTokenDto, err error) {
	user, err := l.repos.UserRepo.GetUser(ctx, userId, userType, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, organizationId, userId)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	token, err = l.GenerateJwtToken(user, organizationUser)
	return token, err
}

// RefreshToken refresh token
func (l *UserLogic) RefreshToken(ctx context.Context, refreshToken string) (token *dto.UserJwtTokenDto, err error) {
	// Verify the refresh token signature with the single shared secret.
	claim, err := utils.ParseToken(refreshToken, l.jwtConfig.RefreshSecret)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	user, err := l.repos.UserRepo.GetUser(ctx, claim.UserId, claim.UserType, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if user.Status != constants.UserStatusActive {
		return nil, errorx.New("user is not active")
	}

	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, claim.OrganizationId, user.Id)
	if err != nil {
		return nil, err
	}
	if organizationUser.Status != constants.OrganizationUserStatusActive {
		return nil, errorx.New("organization user is not active")
	}
	token, err = l.GenerateJwtToken(user, organizationUser)
	return token, err
}

// ResetPasswordWithEmail reset password with email
func (l *UserLogic) ResetPasswordWithEmail(ctx context.Context, input dto.UserResetPasswordByEmailRequestDto) (token *dto.UserJwtTokenDto, err error) {
	email := input.Email
	code := input.Code
	newPwd := input.NewPassword
	if email == "" || code == "" || newPwd == "" {
		return nil, errorx.New("email, code and new password are required")
	}
	redisKey := stringx.Format(constants.ResetPwdCode, input.UserType, email)
	cacheToken, err := l.redisCli.Get(ctx, redisKey)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if cacheToken == "" || cacheToken != code {
		return nil, errorx.WithStack(errorx.New("invalid or expired code"))
	}

	user, err := l.repos.UserRepo.GetUserByEmail(ctx, email, input.UserType)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New("user not found")
	}

	// Hash pwd
	pwdHash, err := utils.HashPassword(newPwd)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	// Update pwd
	err = l.repos.UserRepo.UpdateUserPassword(ctx, user.Id, pwdHash)
	if err != nil {
		return nil, err
	}

	user, organizationUser, err := l.CheckAndUpdateActiveOrganizationByUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err = l.GenerateJwtToken(user, organizationUser)
	if err != nil {
		return nil, err
	}
	// Delete redis key
	_ = l.redisCli.Del(ctx, redisKey)
	return token, nil
}

// GetActiveOrganizationByUser return active organization for user
func (l *UserLogic) GetActiveOrganizationByUser(user *model.User) (currentOrganization *model.Organization, isChanged bool, err error) {
	if user == nil {
		return nil, false, errorx.New("user is empty")
	}
	organizations := user.Organizations
	organizationUsers := user.OrganizationUsers

	// Find current organization
	currentOrganization, ok := slice.FindBy(organizations, func(index int, item *model.Organization) bool {
		isActiveOrganization := (item.Id == user.DefaultOrganizationId || user.DefaultOrganizationId == 0) && item.Status == constants.OrganizationStatusActive
		if !isActiveOrganization {
			return false
		}
		_, ok := slice.FindBy(organizationUsers, func(index int, organizationUser *model.OrganizationUser) bool {
			return organizationUser.OrganizationId == item.Id && organizationUser.Status == constants.OrganizationUserStatusActive
		})

		return ok
	})
	// Find first active organization
	if !ok {
		currentOrganization, ok = slice.FindBy(organizations, func(index int, item *model.Organization) bool {
			isActiveOrganization := item.Status == constants.OrganizationStatusActive
			if !isActiveOrganization {
				return false
			}
			_, ok := slice.FindBy(organizationUsers, func(index int, organizationUser *model.OrganizationUser) bool {
				return organizationUser.OrganizationId == item.Id && organizationUser.Status == constants.OrganizationUserStatusActive
			})

			return ok
		})
	}
	if !ok {
		// No active company
		isChanged = condition.Ternary(user.DefaultOrganizationId == 0, false, true)
		user.DefaultOrganizationId = 0
	} else if user.DefaultOrganizationId != currentOrganization.Id {
		// Company changed
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

// Login
func (l *UserLogic) Login(ctx context.Context, req dto.UserLoginRequestDto, limiter dto.LoginLimiterDto) (token *dto.UserJwtTokenDto, err error) {
	redisKey := stringx.Format(constants.UserLoginFailedTimes, req.UserType, req.Email)
	var loginFailedTimes int
	var orgUser *model.OrganizationUser

	lockPeriod := condition.Ternary(limiter.LimiterItem.LockPeriod > 0, limiter.LimiterItem.LockPeriod, limiter.LimiterItem.Period)
	lockTime := timex.DurationToString(time.Second*time.Duration(lockPeriod), true)
	tooManyAttemptError := errorx.Errorf("too many failed login attempts. please try again in %s", lockTime)
	if limiter.LimiterItem.Rate > 0 {
		loginFailedTimes, err = limiter.RedisCli.GetInt(ctx, redisKey)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		if loginFailedTimes >= limiter.LimiterItem.Rate {
			return nil, tooManyAttemptError
		}
	}

	user, err := l.repos.UserRepo.GetUserByEmail(ctx, req.Email, req.UserType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.New("user not found")
		}
		return nil, errorx.WithStack(err)
	}
	if user == nil {
		return nil, errorx.New("user not found")
	}

	// Check pwd
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		var failedTimes int64
		// Add input wrong pwd times to redis
		failedTimes, err = limiter.RedisCli.IncreaseExpire(ctx, redisKey, limiter.LimiterItem.Period)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		leftTimes := limiter.LimiterItem.Rate - int(failedTimes)
		// Lock login after reach rate
		if leftTimes <= 0 {
			err = limiter.RedisCli.Expire(ctx, redisKey, lockPeriod)
			return nil, tooManyAttemptError
		}
		return nil, errorx.Errorf(
			"the email and password you entered don't match. %d attempts remaining",
			leftTimes,
		)
	}

	if limiter.LimiterItem.Rate > 0 {
		err = limiter.RedisCli.Del(ctx, redisKey)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
	}

	// Check status
	if user.Status == constants.UserStatusInactive {
		return nil, errorx.New("user is not active")
	}
	if user.Status == constants.UserStatusUnverified {
		err = status.Error(codes.PermissionDenied, "user is unverified")
		return nil, err
	}

	user, orgUser, err = l.CheckAndUpdateActiveOrganizationByUser(ctx, user)
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	token, err = l.GenerateJwtToken(user, orgUser)
	if err != nil {
		return nil, err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:             auth.ActionRoleReload,
			Scope:              auth.ScopeOrganizationUser,
			Type:               req.UserType,
			OrganizationId:     orgUser.OrganizationId,
			OrganizationUserId: orgUser.Id,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[req.UserType], msg)
	})

	return
}

// Info return user info
func (l *UserLogic) Info(ctx context.Context, id int64, organizationId int64, userType constants.UserType) (userInfo *dto.UserinfoDto, err error) {
	organizations := []*dto.UserOrganizationDto{}
	user, err := l.repos.UserRepo.GetUser(ctx, id, userType, true)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.WithStack(errorx.New("user not found"))
	}
	if len(user.Organizations) > 0 {
		organizationsPtr := &[]*dto.UserOrganizationDto{}
		err = copierx.Assign[[]*model.Organization, []*dto.UserOrganizationDto](&user.Organizations, organizationsPtr)
		if err != nil {
			return nil, err
		}
		organizations = *organizationsPtr
	}
	var currentOrganization *dto.UserOrganizationDto
	var organizationUser *model.OrganizationUser

	organizations = slice.Filter(organizations, func(i int, organization *dto.UserOrganizationDto) bool {
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
	var ok bool
	organizationUser, ok = slice.FindBy(user.OrganizationUsers, func(index int, item *model.OrganizationUser) bool {
		return item.OrganizationId == organizationId && item.IsDel == 0
	})
	if !ok {
		err = errorx.New("organization user not found")
		return nil, err
	}

	var organizationName string
	var currentOrganizationId = organizationId
	var isOwner bool
	var isAdmin bool

	if currentOrganization == nil && len(organizations) > 0 {
		currentOrganization = organizations[0]
	}

	if currentOrganization != nil {
		currentOrganizationId = user.DefaultOrganizationId
		organizationName = currentOrganization.Name

		roles, err := l.repos.RoleRepo.GetRolesByOrganizationUserId(ctx, organizationUser.Id)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		isAdmin = slice.ContainBy(roles, func(role *model.Role) bool {
			return role.Name == constants.AdminRoleName
		})
		isOwner = slice.ContainBy(user.Organizations, func(item *model.Organization) bool {
			return item.Id == currentOrganization.Id && item.OwnerId == user.Id
		})
	}

	userInfo = &dto.UserinfoDto{
		Id:                 user.Id,
		OrganizationUserId: organizationUser.Id,
		Email:              user.Email,
		FirstName:          organizationUser.FirstName,
		LastName:           organizationUser.LastName,
		FullName:           organizationUser.FullName,
		OrganizationId:     currentOrganizationId,
		OrganizationName:   organizationName,
		Status:             organizationUser.Status,
		Organizations:      organizations,
		IsOwner:            isOwner,
		IsAdmin:            isAdmin,
	}
	return
}

// SwitchOrganization switch user default organization
func (l *UserLogic) SwitchOrganization(ctx context.Context, userId int64, organizationId int64, userType constants.UserType) (user *dto.UserEntityDto, err error) {
	userModel, err := l.repos.UserRepo.GetUser(ctx, userId, userType, true)
	if err != nil {
		return nil, err
	}

	organization, ok := slice.FindBy(userModel.Organizations, func(i int, organization *model.Organization) bool {
		return organization.Id == organizationId
	})
	if !ok {
		return nil, errorx.New("User is not member of this organization")
	}

	if organization.Status != constants.OrganizationStatusActive {
		return nil, errorx.New("organization is not active")
	}

	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, organizationId, userModel.Id)
	if err != nil {
		return nil, err
	}
	if organizationUser == nil || organizationUser.Status != constants.OrganizationUserStatusActive {
		return nil, errorx.New("user is not member of this organization")
	}

	userModel.DefaultOrganizationId = organizationId
	err = l.repos.UserRepo.UpdateColumnsById(ctx, userModel.Id, map[string]interface{}{
		"default_organization_id": userModel.DefaultOrganizationId,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	user, err = copierx.Trans[model.User, dto.UserEntityDto](userModel)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// RegisterWithOrganization creates a new user and organization relationship in a single transaction.
// It validates both user and organization inputs, checks for existence, and creates the necessary
// relationships between them. The function handles both cases of creating a new organization
// or using an existing one.
func (l *UserLogic) RegisterWithOrganization(ctx context.Context, userInput dto.UserRegisterRequestDto, organizationInput dto.OrganizationRegisterRequestDto, userType constants.UserType, isManual bool) (newUser *model.User, organizationUser *model.OrganizationUser, err error) {
	var user *model.User
	var organization *model.Organization
	// Check user exist
	if len(userInput.FirstName) == 0 {
		return nil, nil, errorx.New("user is empty")
	}
	isExisted, err := l.repos.UserRepo.IsExistedByColumns(ctx, map[string]interface{}{
		"email": userInput.Email,
		"type":  userType,
	})
	if err != nil {
		return nil, nil, err
	}
	if isExisted {
		return nil, nil, errorx.New("user already exists")
	}
	user, err = copierx.Trans[dto.UserRegisterRequestDto, model.User](&userInput)
	if err != nil {
		return nil, nil, err
	}

	// Check organization exist
	if len(organizationInput.Name) == 0 {
		return nil, nil, errorx.New("organization name is empty")
	}
	isExisted, err = l.repos.OrganizationRepo.IsExistedByColumns(ctx, map[string]interface{}{
		"name": organizationInput.Name,
		"type": userType,
	})
	if err != nil {
		return nil, nil, err
	}
	if isExisted {
		return nil, nil, errorx.New("organization already exists")
	}
	organization, err = copierx.Trans[dto.OrganizationRegisterRequestDto, model.Organization](&organizationInput)
	if err != nil {
		return nil, nil, err
	}

	err = l.repos.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		organization.Status = condition.Ternary(isManual, organizationInput.Status, constants.OrganizationStatusNotActivated)
		// organization.Status = constants.OrganizationStatusNotActivated
		organization, innerErr = l.repos.OrganizationRepo.CreateWithTx(ctx, tx, organization, false)
		if innerErr != nil {
			return innerErr
		}

		user.DefaultOrganizationId = organization.Id
		// Create user
		roles := []*model.Role{
			{
				Name:           constants.AdminRoleName,
				Type:           userType,
				OrganizationId: organization.Id,
			},
		}
		user.RemoveRelationships()
		user, innerErr = l.repos.UserRepo.CreateWithTx(ctx, tx, user, false)
		if innerErr != nil {
			return innerErr
		}
		// Update organization's owner
		organization.OwnerId = user.Id
		innerErr = l.repos.OrganizationRepo.UpdateOwnerIdWithTx(ctx, tx, organization.Id, user.Id)
		if innerErr != nil {
			return innerErr
		}
		user.DefaultOrganizationId = organization.Id

		// Create organization user
		organizationUser = &model.OrganizationUser{
			OrganizationId: organization.Id,
			UserId:         user.Id,
			FirstName:      user.FirstName,
			MidName:        user.MidName,
			LastName:       user.LastName,
			Status:         condition.Ternary(isManual, constants.OrganizationUserStatusActive, constants.OrganizationUserStatusNotActivated),
			FullName:       user.FullName,
			Roles:          roles,
		}
		// Create organization and user relation
		organizationUser, innerErr = l.repos.OrganizationUserRepo.CreateWithTx(ctx, tx, organizationUser, true)
		if innerErr != nil {
			return innerErr
		}

		// Sync organization template
		organizationTemplateLogic := NewOrganizationTemplateLogic(l.BaseLogic)
		innerErr = organizationTemplateLogic.SyncByOrganizationIdWithTx(ctx, tx, organization)
		if innerErr != nil {
			return innerErr
		}
		// Init lender organization setting
		organizationSettingLogic := NewOrganizationSettingLogic(l.BaseLogic)
		innerErr = organizationSettingLogic.InitLenderOrganizationSettingWithTx(ctx, tx, organization)
		if innerErr != nil {
			return innerErr
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganization,
			Type:           userInput.Type,
			OrganizationId: organization.Id,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[userInput.Type], msg)
	})

	return user, organizationUser, nil
}

// ValidatePortalPassword enforces the same rules as lender portal password changes.
func (l *UserLogic) ValidatePortalPassword(password string) bool {
	return l.validatePassword(password)
}

// CreatePrimaryOwnerForLenderOrgWithTx creates the lender owner user for an organization
// inserted without one (e.g. OpenAI discover). Caller must use the same transaction for the org update.
func (l *UserLogic) CreatePrimaryOwnerForLenderOrgWithTx(ctx context.Context, tx *gorm.DB, organization *model.Organization, userInput dto.UserRegisterRequestDto) (err error) {
	if organization == nil {
		return errorx.New("organization is empty")
	}
	if organization.Type != constants.UserTypeLender {
		return errorx.New("invalid organization type")
	}
	if organization.OwnerId != 0 {
		return errorx.New("organization already has an owner")
	}
	if organization.Source != constants.OrganizationSourceSearch {
		return errorx.New("owner onboarding is only allowed for search-sourced lenders")
	}
	if len(userInput.Email) == 0 || len(userInput.PasswordHash) == 0 {
		return errorx.New("owner email and password are required")
	}
	if len(userInput.FirstName) == 0 || len(userInput.LastName) == 0 {
		return errorx.New("owner first and last name are required")
	}
	isExisted, err := l.repos.UserRepo.IsExistedByColumnsWithTx(ctx, tx, map[string]interface{}{
		"email": userInput.Email,
		"type":  constants.UserTypeLender,
	})
	if err != nil {
		return err
	}
	if isExisted {
		return errorx.New("user already exists")
	}

	user, err := copierx.Trans[dto.UserRegisterRequestDto, model.User](&userInput)
	if err != nil {
		return err
	}
	user.PasswordHash = userInput.PasswordHash
	user.DefaultOrganizationId = organization.Id
	user.RemoveRelationships()

	user, err = l.repos.UserRepo.CreateWithTx(ctx, tx, user, false)
	if err != nil {
		return err
	}

	err = l.repos.OrganizationRepo.UpdateOwnerIdWithTx(ctx, tx, organization.Id, user.Id)
	if err != nil {
		return err
	}
	organization.OwnerId = user.Id

	roles := []*model.Role{
		{
			Name:           constants.AdminRoleName,
			Type:           constants.UserTypeLender,
			OrganizationId: organization.Id,
		},
	}
	organizationUser := &model.OrganizationUser{
		OrganizationId: organization.Id,
		UserId:         user.Id,
		FirstName:      user.FirstName,
		MidName:        user.MidName,
		LastName:       user.LastName,
		Status:         constants.OrganizationUserStatusActive,
		FullName:       user.FullName,
		Roles:          roles,
	}
	_, err = l.repos.OrganizationUserRepo.CreateWithTx(ctx, tx, organizationUser, true)
	return err
}

// PublishRoleReloadForOrganization notifies Casbin consumers after roles change for an org.
func (l *UserLogic) PublishRoleReloadForOrganization(ctx context.Context, organizationId int64, userType constants.UserType) {
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganization,
			Type:           userType,
			OrganizationId: organizationId,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[userType], msg)
	})
}

// VerifyUserEmail verity user email
func (l *UserLogic) VerifyUserEmail(ctx context.Context, email string, code string, userType constants.UserType) (user *model.User, err error) {
	var token string
	if email == "" || code == "" {
		return nil, errorx.New("email and code are required")
	}
	// Check code
	verifyCodeKey := stringx.Format(constants.UserRegisterVerifyCode, userType, email)
	token, err = l.redisCli.Get(ctx, verifyCodeKey)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	// Check if code exists (expired or invalid)
	if len(token) == 0 || token != code {
		return nil, errorx.New("invalid or expired code")
	}
	defer func() {
		if err == nil {
			_ = l.redisCli.Del(ctx, verifyCodeKey)
		}
	}()

	var activatedOrgIds []int64
	err = l.repos.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// Get user
		user, innerErr = l.repos.UserRepo.GetUserByEmailWithTx(ctx, tx, email, userType)
		if innerErr != nil {
			return innerErr
		}
		if user.Status != constants.UserStatusUnverified {
			return errorx.New("email is already verified")
		}
		// Update user status
		innerErr = l.repos.UserRepo.UpdateColumnsByIdWithTx(ctx, tx, user.Id, map[string]interface{}{
			"status": constants.UserStatusActive,
		})
		if innerErr != nil {
			return innerErr
		}
		// Update user's organization status
		activatedOrgIds, innerErr = l.handleOwnedOrganizationWithTx(ctx, tx, user)
		if innerErr != nil {
			return innerErr
		}

		// Update organization user status
		innerErr = l.repos.OrganizationUserRepo.ActiveByUserIdWithTx(ctx, tx, user.Id)
		if innerErr != nil {
			return innerErr
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Subscribe default plan for newly activated organizations
	subLogic := NewOrganizationSubscriptionLogic(l.BaseLogic, l.authTopicMap)
	for _, orgId := range activatedOrgIds {
		_ = subLogic.SubscribeDefaultPlan(ctx, orgId, userType)
	}

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganization,
			Type:           userType,
			OrganizationId: user.DefaultOrganizationId,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[userType], msg)
	})

	return
}

// ChangePassword change password
func (l *UserLogic) ChangePassword(ctx context.Context, input dto.UserChangePasswordRequestDto) (err error) {
	user, err := l.repos.UserRepo.GetUser(ctx, input.UserId, input.UserType, false)
	if err != nil {
		return err
	}
	if !l.validatePassword(input.NewPassword) {
		return errorx.New("new password is invalid")
	}

	oldPassword := input.OldPassword
	if !utils.CheckPassword(oldPassword, user.PasswordHash) {
		return errorx.New("old password is incorrect")
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

// handleOwnedOrganization activates owned organizations and returns IDs of newly activated orgs.
func (l *UserLogic) handleOwnedOrganizationWithTx(ctx context.Context, tx *gorm.DB, user *model.User) (activatedOrgIds []int64, err error) {
	if user == nil {
		return nil, errorx.New("user is empty")
	}

	// Find not-activated orgs before status update
	notActivatedOrgs, err := l.repos.OrganizationRepo.FindByOwnerIdAndStatusWithTx(ctx, tx, user.Id, user.Type, constants.OrganizationStatusNotActivated)
	if err != nil {
		return nil, err
	}

	// Update owned organization status
	err = l.repos.OrganizationRepo.UpdateOwnerOrganizationStatusWithTx(ctx, tx, user.Id, constants.OrganizationStatusActive)
	if err != nil {
		return nil, err
	}
	// Get owned organization ids
	organizations, err := l.repos.OrganizationRepo.GetOrganizationsByOwnerIdWithTx(ctx, tx, user.Id, user.Type)
	if err != nil {
		return nil, err
	}
	// Sync owned organization templates
	organizationTemplateLogic := NewOrganizationTemplateLogic(l.BaseLogic)
	for index := range organizations {
		err = organizationTemplateLogic.SyncByOrganizationIdWithTx(ctx, tx, organizations[index])
		if err != nil {
			return nil, err
		}
	}

	for _, org := range notActivatedOrgs {
		activatedOrgIds = append(activatedOrgIds, org.Id)
	}
	return activatedOrgIds, nil
}
