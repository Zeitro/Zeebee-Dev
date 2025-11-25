package logic

import (
	"context"

	"github.com/duke-git/lancet/v2/condition"
	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/errorx"
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

func (l *UserLogic) GenerateJwtToken(user *model.User, config JwtConfig) (token *UserJwtToken, err error) {
	if user == nil {
		return nil, errorx.WithStack(errorx.New("User Is Empty."))
	}
	userClaim := utils.JwtUser{
		UserId:         user.Id,
		OrganizationId: user.DefaultOrganizationId,
		Username:       user.FullName,
		Role:           nil,
		UserType:       constants.UserTypeBorrower,
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
		return item.Id == user.DefaultOrganizationId && item.Status == constants.OrganizationStatusActive
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
