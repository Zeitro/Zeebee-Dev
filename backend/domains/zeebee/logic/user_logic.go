package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type UserLogic struct {
	repos *repo.AllRepos
}

func NewUserLogic(repos *repo.AllRepos) *UserLogic {
	return &UserLogic{
		repos: repos,
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

// SwitchOrganization
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
