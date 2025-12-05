package logic

import (
	"context"
	"fmt"

	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/common/constants"
	"zeebee-backend/common/copierx"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type OrganizationUserLogic struct {
	repos       *repo.AllRepos
	redisCli    *redisx.Redisx
	emailClient *email2.Email
}

func NewOrganizationUserLogic(repos *repo.AllRepos, redisCli *redisx.Redisx, emailClient *email2.Email) *OrganizationUserLogic {
	return &OrganizationUserLogic{
		repos:       repos,
		redisCli:    redisCli,
		emailClient: emailClient,
	}
}

// CreateOrganizationUser create organization user
func (l *OrganizationUserLogic) CreateOrganizationUser(ctx context.Context, createDto dto.OrganizationUserCreateRequestDto) (result *bo.OrganizationUserBo, err error) {
	isInOrganization, err := l.repos.OrganizationUserRepo.CheckOrganizationUserExistedByEmail(ctx, createDto.OrganizationId, createDto.Email, createDto.Type)
	if err != nil {
		return nil, err
	}
	if isInOrganization {
		err = errorx.New("The User Is In The Organization")
		return nil, err
	}

	user, err := l.repos.UserRepo.GetUserByEmail(ctx, createDto.Email, createDto.Type)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	var pwd string
	var organizationUser *model.OrganizationUser
	var organization *model.Organization
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// get organization
		organization, innerErr = l.repos.OrganizationRepo.GetOrganization(ctx, createDto.OrganizationId, createDto.Type)
		if innerErr != nil {
			return innerErr
		}
		if organization.Status != constants.OrganizationStatusActive {
			innerErr = errorx.New("The Organization Is Not Active")
			return innerErr
		}

		// create user if not exist
		if user == nil || user.Id == 0 {
			pwd = random.RandString(10)
			pwd = strutil.UpperFirst(pwd)
			pwdHash, innerErr := utils.HashPassword(pwd)
			if innerErr != nil {
				return errorx.WithStack(innerErr)
			}
			user = &model.User{
				Email:                 createDto.Email,
				Type:                  createDto.Type,
				Status:                constants.UserStatusUnverified,
				PasswordHash:          pwdHash,
				DefaultOrganizationId: createDto.OrganizationId,
			}
			user, innerErr = l.repos.UserRepo.CreateWithTx(ctx, tx, user, false)
			if innerErr != nil {
				return innerErr
			}
		}

		// create organization user
		organizationUser = &model.OrganizationUser{
			UserId:         user.Id,
			OrganizationId: createDto.OrganizationId,
			Phone:          createDto.Phone,
			Remark:         createDto.Remark,
			Status:         constants.OrganizationUserStatusActive,
		}

		// get organization roles
		if len(createDto.RoleIds) > 0 {
			roles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, createDto.OrganizationId, createDto.RoleIds)
			if innerErr != nil {
				return innerErr
			}
			organizationUser.Roles = roles
		}

		organizationUser, innerErr = l.repos.OrganizationUserRepo.CreateWithTx(ctx, tx, organizationUser, true)
		if innerErr != nil {
			return innerErr
		}
		return innerErr
	})
	if err != nil {
		return nil, err
	}
	// send email
	go func() {
		newCtx := context.WithoutCancel(ctx)
		err := l.sendInviteEmail(newCtx, sendInviteUserEmailInput{
			emailConfig:      createDto.EmailConfig,
			firstName:        user.FirstName,
			lastName:         user.LastName,
			targetEmail:      createDto.Email,
			userType:         createDto.Type,
			hasPwd:           len(pwd) > 0,
			pwd:              pwd,
			organizationName: organization.Name,
			operator:         createDto.Operator,
		})
		if err != nil {
			err = errorx.Wrap(err, "Send invite Email Error")
			logx.WithContext(ctx).Error(err)
		}
	}()

	result = &bo.OrganizationUserBo{}
	err = copier.CopyWithOption(result, organizationUser, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

type sendInviteUserEmailInput struct {
	emailConfig      dto.SendInviteUserEmailConfig
	firstName        string
	lastName         string
	targetEmail      string
	userType         constants.UserType
	hasPwd           bool
	pwd              string
	organizationName string
	operator         string
}

// sendInviteEmail send invite email
func (l *OrganizationUserLogic) sendInviteEmail(ctx context.Context, input sendInviteUserEmailInput) (err error) {

	var code string
	// generate invite code
	codeKey := stringx.Format(constants.UserInviteVerifyCode, input.userType, input.targetEmail)
	code, _ = l.redisCli.Get(ctx, codeKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		err = l.redisCli.SetExpire(ctx, codeKey, code, input.emailConfig.ExpireSeconds)
		if err != nil {
			return errorx.WithStack(err)
		}
	}
	// send invite email
	data := map[string]interface{}{
		"url": fmt.Sprintf("%s?code=%s&email=%s&first_name=%s&last_name=%s&oranization_name=%s",
			input.emailConfig.Url,
			code,
			input.targetEmail,
			input.firstName,
			input.lastName,
			input.organizationName,
		),
		"has_pwd":           input.hasPwd,
		"pwd":               input.pwd,
		"operator":          input.operator,
		"organization_name": input.organizationName,
	}
	content, err := email2.ParseTemplate(input.emailConfig.TplName, data)
	if err != nil {
		return errorx.WithStack(err)
	}
	msg := email2.Message{
		Subject: fmt.Sprintf("%s has invited you to join %s", input.operator, input.organizationName),
		From:    input.emailConfig.HostEmail,
		To:      []string{input.targetEmail},
		Body:    content,
		IsHTML:  true,
	}
	err = l.emailClient.SendEmail(msg)
	return errorx.WithStack(err)

}

// AllOrganizationUsers
func (l *OrganizationUserLogic) AllOrganizationUsers(ctx context.Context, organizationId int64) (result []*bo.OrganizationUserBo, err error) {
	repoResults, err := l.repos.OrganizationUserRepo.AllOrganizationUsers(ctx, organizationId)
	if err != nil {
		return nil, err
	}
	result = []*bo.OrganizationUserBo{}
	err = copier.CopyWithOption(&result, &repoResults, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return result, nil
}
