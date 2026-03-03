package logic

import (
	"context"
	"regexp"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type JwtConfig struct {
	AccessSecret  string
	AccessExpire  int64
	RefreshSecret string
	RefreshExpire int64
}

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

// GetUserClaimFromCtx get user claim from context
func (l *UserLogic) GetUserClaimFromCtx(ctx context.Context) (claim utils.JwtUser, err error) {
	// Get from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = status.Error(codes.Unauthenticated, "Not login")
		return claim, errorx.WithStack(err)
	}
	tokens := md[utils.JwtTokenCtxKey]
	if len(tokens) == 0 {
		err = status.Error(codes.Unauthenticated, "Not login")
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
func (l *UserLogic) GenerateJwtToken(user *model.User, organizationUser *model.OrganizationUser, userType constants.UserType) (token *dto.UserJwtTokenDto, err error) {
	if user == nil {
		return nil, errorx.WithStack(errorx.New("User is empty"))
	}
	userClaim := utils.JwtUser{
		UserId:         user.Id,
		OrganizationId: user.DefaultOrganizationId,
		Username:       user.FullName,
		Role:           nil,
		UserType:       userType,
	}
	if organizationUser != nil {
		userClaim.OrganizationId = organizationUser.OrganizationId
		userClaim.OrganizationUserId = organizationUser.Id
		userClaim.Username = organizationUser.FullName
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
	token = &dto.UserJwtTokenDto{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpire:  l.jwtConfig.AccessExpire,
		RefreshTokenExpire: l.jwtConfig.RefreshExpire,
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
	token, err = l.GenerateJwtToken(user, organizationUser, userType)
	return token, err
}

// RefreshToken refresh token
func (l *UserLogic) RefreshToken(ctx context.Context, refreshToken string, userType constants.UserType) (token *dto.UserJwtTokenDto, err error) {
	// Parse refresh token
	claim, err := utils.ParseToken(refreshToken, l.jwtConfig.RefreshSecret)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if claim.UserType != userType {
		return nil, errorx.New("User type does not match")
	}

	user, err := l.repos.UserRepo.GetUser(ctx, claim.UserId, claim.UserType, false)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if user.Status != constants.UserStatusActive {
		return nil, errorx.New("User is not active")
	}

	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, claim.OrganizationId, user.Id)
	if err != nil {
		return nil, err
	}
	if organizationUser.Status != constants.OrganizationUserStatusActive {
		return nil, errorx.New("Organization user is not active")
	}
	token, err = l.GenerateJwtToken(user, organizationUser, claim.UserType)
	return token, err
}

// ResetPasswordWithEmail reset password with email
func (l *UserLogic) ResetPasswordWithEmail(ctx context.Context, input dto.UserResetPasswordByEmailRequestDto) (token *dto.UserJwtTokenDto, err error) {
	email := input.Email
	code := input.Code
	newPwd := input.NewPassword
	if email == "" || code == "" || newPwd == "" {
		return nil, errorx.New("Email, code and new password are required")
	}
	redisKey := stringx.Format(constants.ResetPwdCode, input.UserType, email)
	cacheToken, err := l.redisCli.Get(ctx, redisKey)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if cacheToken == "" || cacheToken != code {
		return nil, errorx.WithStack(errorx.New("Invalid or expired code"))
	}

	user, err := l.repos.UserRepo.GetUserByEmail(ctx, email, input.UserType)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.New("User not found")
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
		return nil, errorx.WithStack(err)
	}

	token, err = l.GenerateJwtToken(user, organizationUser, input.UserType)
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
		return nil, false, errorx.New("User is empty")
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
	if limiter.Rate > 0 {
		loginFailedTimes, err = limiter.RedisCli.GetInt(ctx, redisKey)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		if loginFailedTimes >= limiter.Rate {
			return nil, errorx.New("Login failed times exceeded, please try again later")
		}
	}

	user, err := l.repos.UserRepo.GetUserByEmail(ctx, req.Email, req.UserType)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	if user == nil {
		return nil, errorx.New("User not found")
	}

	// Check pwd
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		// Add input wrong pwd times to redis
		_, err = limiter.RedisCli.IncreaseExpire(ctx, redisKey, limiter.Period)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
		return nil, errorx.WithStack(errorx.New("The email and password you entered don't match."))
	}

	if limiter.Rate > 0 {
		err = limiter.RedisCli.Del(ctx, redisKey)
		if err != nil {
			return nil, errorx.WithStack(err)
		}
	}

	// Check status
	if user.Status == constants.UserStatusInactive {
		return nil, errorx.New("User is not active")
	}
	if user.Status == constants.UserStatusUnverified {
		err = status.Error(codes.PermissionDenied, "User is unverified")
		return nil, err
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

// Info return user info
func (l *UserLogic) Info(ctx context.Context, id int64, organizationId int64, userType constants.UserType) (userInfo *dto.UserinfoDto, err error) {
	organizations := []*dto.UserOrganizationDto{}
	user, err := l.repos.UserRepo.GetUser(ctx, id, userType, true)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorx.WithStack(errorx.New("User not found"))
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

	organizationUser, ok := slice.FindBy(user.OrganizationUsers, func(index int, item *model.OrganizationUser) bool {
		return item.OrganizationId == organizationId && item.IsDel == 0
	})
	if !ok {
		err = errorx.New("Organization user not found")
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
		OrganizationId:     organizationId,
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
		return nil, errorx.New("Organization is not active")
	}

	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, organizationId, userModel.Id)
	if err != nil {
		return nil, err
	}
	if organizationUser == nil || organizationUser.Status != constants.OrganizationUserStatusActive {
		return nil, errorx.New("User is not member of this organization")
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
		return nil, nil, errorx.New("User Is Empty")
	}
	isExisted, err := l.repos.UserRepo.IsExistedByColumns(ctx, map[string]interface{}{
		"email": userInput.Email,
		"type":  userType,
	})
	if err != nil {
		return nil, nil, err
	}
	if isExisted {
		return nil, nil, errorx.New("User already exists")
	}
	user, err = copierx.Trans[dto.UserRegisterRequestDto, model.User](&userInput)
	if err != nil {
		return nil, nil, err
	}

	// Check organization exist
	if len(organizationInput.Name) == 0 {
		return nil, nil, errorx.New("Organization name is empty")
	}
	isExisted, err = l.repos.OrganizationRepo.IsExistedByColumns(ctx, map[string]interface{}{
		"name": organizationInput.Name,
		"type": userType,
	})
	if err != nil {
		return nil, nil, err
	}
	if isExisted {
		return nil, nil, errorx.New("Organization already exists")
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
		eg, _ := errgroup.WithContext(ctx)
		// Sync organization template
		eg.Go(func() error {
			organizationTemplateLogic := NewOrganizationTemplateLogic(l.repos)
			return organizationTemplateLogic.SyncByOrganizationIdWithTx(ctx, tx, organization)
		})
		// Init lender organization setting
		eg.Go(func() error {
			organizationSettingLogic := NewOrganizationSettingLogic(l.repos)
			return organizationSettingLogic.InitLenderOrganizationSettingWithTx(ctx, tx, organization)
		})
		innerErr = eg.Wait()
		return errorx.WithStack(innerErr)
	})
	if err != nil {
		return nil, nil, err
	}

	return user, organizationUser, nil
}

// VerifyUserEmail verity user email
func (l *UserLogic) VerifyUserEmail(ctx context.Context, email string, code string, userType constants.UserType) (user *model.User, err error) {
	var token string
	if email == "" || code == "" {
		return nil, errorx.New("Email and code are required")
	}
	// Check code
	verifyCodeKey := stringx.Format(constants.UserRegisterVerifyCode, userType, email)
	token, err = l.redisCli.Get(ctx, verifyCodeKey)
	if err != nil {
		return nil, errorx.WithStack(err)
	}
	// Check if code exists (expired or invalid)
	if len(token) == 0 || token != code {
		return nil, errorx.New("Invalid or expired code")
	}
	defer func() {
		if err == nil {
			_ = l.redisCli.Del(ctx, verifyCodeKey)
		}
	}()

	err = l.repos.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// Get user
		user, innerErr = l.repos.UserRepo.GetUserByEmailWithTx(ctx, tx, email, userType)
		if innerErr != nil {
			return innerErr
		}
		if user.Status != constants.UserStatusUnverified {
			return errorx.New("Email is already verified")
		}
		// Update user status
		innerErr = l.repos.UserRepo.UpdateColumnsByIdWithTx(ctx, tx, user.Id, map[string]interface{}{
			"status": constants.UserStatusActive,
		})
		if innerErr != nil {
			return innerErr
		}
		// Update user's organization status
		innerErr = l.handleOwnedOrganizationWithTx(ctx, tx, user)
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
	return
}

// ChangePassword change password
func (l *UserLogic) ChangePassword(ctx context.Context, input dto.UserChangePasswordRequestDto) (err error) {
	user, err := l.repos.UserRepo.GetUser(ctx, input.UserId, input.UserType, false)
	if err != nil {
		return err
	}
	if !l.validatePassword(input.NewPassword) {
		return errorx.New("New password is invalid")
	}

	oldPassword := input.OldPassword
	if !utils.CheckPassword(oldPassword, user.PasswordHash) {
		return errorx.New("Old password is incorrect")
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

// handleOwnedOrganization
func (l *UserLogic) handleOwnedOrganizationWithTx(ctx context.Context, tx *gorm.DB, user *model.User) (err error) {
	if user == nil {
		return errorx.New("User is empty")
	}
	// Update owned organization status
	err = l.repos.OrganizationRepo.UpdateOwnerOrganizationStatusWithTx(ctx, tx, user.Id, constants.OrganizationStatusActive)
	if err != nil {
		return err
	}
	// Get owned organization ids
	organizations, err := l.repos.OrganizationRepo.GetOrganizationsByOwnerIdWithTx(ctx, tx, user.Id, user.Type)
	if err != nil {
		return err
	}
	// Sync owned organization templates
	organizationTemplateLogic := NewOrganizationTemplateLogic(l.repos)
	for index := range organizations {
		err = organizationTemplateLogic.SyncByOrganizationIdWithTx(ctx, tx, organizations[index])
		if err != nil {
			return err
		}
	}

	return err
}
