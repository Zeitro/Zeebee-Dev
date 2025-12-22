package logic

import (
	"context"
	"fmt"
	"regexp"

	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
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

// Page return organization user list
func (l *OrganizationUserLogic) Page(ctx context.Context, pageReq dto.OrganizationUserPageRequestDto) (page bo.PageResultBo[*bo.OrganizationUserBo], err error) {
	organizationId := pageReq.Filter.OrganizationId
	organization, err := l.repos.OrganizationRepo.FindOneById(ctx, organizationId)
	if err != nil {
		return page, err
	}
	organizationOwnerId := organization.OwnerId

	results, err := l.repos.OrganizationUserRepo.Page(ctx, pageReq)
	if err != nil {
		return page, err
	}

	for i, entity := range results.List {
		roles, err := l.repos.RoleRepo.GetOrganizationUserRoles(ctx, entity.Id)
		if err != nil {
			return page, err
		}
		results.List[i].Roles = roles
	}

	err = copier.CopyWithOption(&page, results, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return page, errorx.GetFailedTransRespErr(err)
	}

	slice.ForEach(page.List, func(index int, item *bo.OrganizationUserBo) {
		page.List[index].IsOwner = item.UserId == organizationOwnerId
	})
	return page, err
}

// CreateAdminOrganizationUser Create admin organization user
func (l *OrganizationUserLogic) CreateAdminOrganizationUser(ctx context.Context, input dto.OrganizationUserCreateAdminRequestDto) (result *bo.OrganizationUserBo, err error) {
	isInOrganization, err := l.repos.OrganizationUserRepo.CheckOrganizationUserExistedByEmail(ctx, input.OrganizationId, input.Email, input.UserType)
	if err != nil {
		return nil, err
	}
	if isInOrganization {
		err = errorx.New("The User Is In The Organization")
		return nil, err
	}
	user, err := l.repos.UserRepo.GetUserByEmail(ctx, input.Email, input.UserType)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if l.validatePassword(input.Password) {
		err = errorx.New("The Password Is Invalid")
		return nil, err
	}

	var organizationUser *model.OrganizationUser
	var organization *model.Organization
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// get organization
		organization, innerErr = l.repos.OrganizationRepo.GetOrganizationWithTx(ctx, tx, input.OrganizationId, input.UserType, false)
		if innerErr != nil {
			return innerErr
		}
		if organization.Status != constants.OrganizationStatusActive {
			innerErr = errorx.New("The Organization Is Not Active")
			return innerErr
		}
		// create user if not exist
		if user == nil || user.Id == 0 {

			pwdHash, innerErr := utils.HashPassword(input.Password)
			if innerErr != nil {
				return errorx.WithStack(innerErr)
			}
			user = &model.User{
				Email:                 input.Email,
				Type:                  input.UserType,
				Status:                constants.UserStatusActive,
				PasswordHash:          pwdHash,
				DefaultOrganizationId: input.OrganizationId,
			}
			user, innerErr = l.repos.UserRepo.CreateWithTx(ctx, tx, user, false)
			if innerErr != nil {
				return innerErr
			}
		}

		// create organization user
		organizationUser = &model.OrganizationUser{
			UserId:         user.Id,
			OrganizationId: input.OrganizationId,
			Phone:          input.Phone,
			Remark:         input.Remark,
			Status:         constants.OrganizationUserStatusActive,
		}

		// get organization roles
		if len(input.RoleIds) > 0 {
			roles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, input.OrganizationId, input.RoleIds)
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
	result = l.transToBo(organizationUser)
	return result, nil
}

// InviteOrganizationUser Invite organization user to join organization
func (l *OrganizationUserLogic) InviteOrganizationUser(ctx context.Context, createDto dto.OrganizationUserCreateRequestDto) (result *bo.OrganizationUserBo, err error) {
	isInOrganization, err := l.repos.OrganizationUserRepo.CheckOrganizationUserExistedByEmail(ctx, createDto.OrganizationId, createDto.Email, createDto.UserType)
	if err != nil {
		return nil, err
	}
	if isInOrganization {
		err = errorx.New("The User Is In The Organization")
		return nil, err
	}

	user, err := l.repos.UserRepo.GetUserByEmail(ctx, createDto.Email, createDto.UserType)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	var pwd string
	var organizationUser *model.OrganizationUser
	var organization *model.Organization
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// get organization
		organization, innerErr = l.repos.OrganizationRepo.GetOrganizationWithTx(ctx, tx, createDto.OrganizationId, createDto.UserType, false)
		if innerErr != nil {
			return innerErr
		}
		if organization.Status != constants.OrganizationStatusActive {
			innerErr = errorx.New("The Organization Is Not Active")
			return innerErr
		}

		// create user if not exist
		if user == nil || user.Id == 0 {
			pwd = random.RandString(10) + random.RandNumeral(5)
			pwd = strutil.UpperFirst(pwd)
			pwdHash, innerErr := utils.HashPassword(pwd)
			if innerErr != nil {
				return errorx.WithStack(innerErr)
			}
			user = &model.User{
				Email:                 createDto.Email,
				Type:                  createDto.UserType,
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
			Status:         constants.OrganizationUserStatusNotActivated,
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
		err := l.sendInviteEmail(newCtx, SendInviteUserEmailInput{
			EmailConfig:      createDto.EmailConfig,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			TargetEmail:      createDto.Email,
			UserType:         createDto.UserType,
			Pwd:              pwd,
			OrganizationName: organization.Name,
			OrganizationId:   organizationUser.OrganizationId,
			Operator:         createDto.Operator,
		})
		if err != nil {
			err = errorx.Wrap(err, "Send invite Email Error")
			logx.WithContext(ctx).Error(err)
		}
	}()

	result = l.transToBo(organizationUser)

	return result, nil
}

// SendInviteEmail send invite email to organization user
func (l *OrganizationUserLogic) SendInviteEmail(ctx context.Context, input dto.OrganizationUserSendInviteEmailRequestDto) (err error) {
	// get organization user
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, input.Id, input.OrganizationId)
	if err != nil {
		return err
	}
	if organizationUser.Status != constants.OrganizationUserStatusNotActivated {
		err = errorx.New("The User Status Is Not Activated")
		return err
	}
	user := organizationUser.User
	if user == nil {
		err = errorx.New("User Not Found")
		return err
	}
	// get  organization
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, input.OrganizationId, input.UserType, false)
	if err != nil {
		return err
	}
	if organization.Status != constants.OrganizationStatusActive {
		err = errorx.New("The Organization Is Not Active")
		return err
	}

	inputKey := stringx.Format(constants.UserInviteUserPwd, input.UserType, user.Email, input.OrganizationId)
	pwd, err := l.redisCli.Get(ctx, inputKey)
	if err != nil {
		return err
	}
	previousInput := SendInviteUserEmailInput{
		EmailConfig:      input.EmailConfig,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		TargetEmail:      user.Email,
		UserType:         input.UserType,
		OrganizationName: organization.Name,
		OrganizationId:   input.OrganizationId,
		Operator:         input.Operator,
		Pwd:              pwd,
	}
	go func() {
		newCtx := context.WithoutCancel(ctx)
		err := l.sendInviteEmail(newCtx, previousInput)
		if err != nil {
			err = errorx.Wrap(err, "Send invite Email Error")
			logx.WithContext(ctx).Error(err)
		}
	}()

	return nil
}

// AcceptInvitation accept organization invitation to join organization
func (l *OrganizationUserLogic) AcceptInvitation(ctx context.Context, input dto.OrganizationUserAcceptInvitationRequestDto) (user bo.UserBo, err error) {
	key := stringx.Format(constants.UserInviteCode, input.UserType, input.Email, input.OrganizationId)
	code, err := l.redisCli.Get(ctx, key)
	if err != nil {
		return user, errorx.WithStack(err)
	}
	if len(code) == 0 || code != input.Code {
		return user, errorx.New("Invalid Code")
	}
	var organizationUserModel *model.OrganizationUser
	var userModel *model.User
	err = l.repos.OrganizationUserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		var ok bool
		userModel, innerErr = l.repos.UserRepo.GetUserByEmailWithTx(ctx, tx, input.Email, input.UserType)
		if innerErr != nil {
			return innerErr
		}
		organizationUserModel, ok = slice.FindBy(userModel.OrganizationUsers, func(index int, item *model.OrganizationUser) bool {
			return input.OrganizationId == item.OrganizationId
		})
		if !ok {
			return errorx.New("Organization User Is Not Found")
		}
		if organizationUserModel.Status != constants.OrganizationUserStatusNotActivated {
			return errorx.New("The User Has Been Activated")
		}
		// update organization user
		innerErr = l.repos.OrganizationUserRepo.UpdateColumnsByIdWithTx(ctx, tx, organizationUserModel.Id, map[string]interface{}{
			"status":     constants.OrganizationUserStatusActive,
			"first_name": input.FirstName,
			"last_name":  input.LastName,
			"full_name":  stringx.FormatFullName(input.FirstName, "", input.LastName),
		})
		if innerErr != nil {
			return innerErr
		}

		userModel.DefaultOrganizationId = input.OrganizationId

		if len(userModel.FirstName) == 0 {
			userModel.FirstName = input.FirstName
			userModel.LastName = input.LastName
			userModel.FullName = stringx.FormatFullName(input.FirstName, "", input.LastName)
		}
		if userModel.Status == constants.UserStatusUnverified {
			userModel.Status = constants.UserStatusActive
		}
		userModel.RemoveRelationships()
		userModel, innerErr = l.repos.UserRepo.UpdateOneByIdWithTx(ctx, tx, userModel.Id, userModel)
		if innerErr != nil {
			return innerErr
		}

		return innerErr
	})
	if err != nil {
		return user, err
	}
	err = copier.CopyWithOption(&user, userModel, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	return user, nil
}

// UpdateOrganizationUser update organization user information
func (l *OrganizationUserLogic) UpdateOrganizationUser(ctx context.Context, reqDto dto.OrganizationUserUpdateRequestDto) (result *bo.OrganizationUserBo, err error) {
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, reqDto.Id, reqDto.OrganizationId)
	if err != nil {
		return nil, err
	}
	// saveModel := &model.OrganizationUser{}
	err = copier.CopyWithOption(organizationUser, reqDto, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	organizationUser.FullName = stringx.FormatFullName(organizationUser.FirstName, organizationUser.MidName, organizationUser.LastName)

	// var organizationUser *model.OrganizationUser
	var organization *model.Organization
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		oldRoles := organizationUser.Roles
		// get organization
		organization, innerErr = l.repos.OrganizationRepo.GetOrganizationWithTx(ctx, tx, reqDto.OrganizationId, reqDto.UserType, false)
		if innerErr != nil {
			return innerErr
		}
		if organization.Status != constants.OrganizationStatusActive {
			innerErr = errorx.New("The Organization Is Not Active")
			return innerErr
		}
		// update organization user
		organizationUser.RemoveRelationships()
		_, innerErr = l.repos.OrganizationUserRepo.UpdateOneWithTx(ctx, tx, organizationUser)
		if innerErr != nil {
			return innerErr
		}
		// update roles
		roles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, reqDto.OrganizationId, reqDto.RoleIds)
		if innerErr != nil {
			return innerErr
		}

		newRoles := []*model.Role{}
		removedRoles := []*model.Role{}
		newRoles = slice.DifferenceWith(roles, oldRoles, func(item1, item2 *model.Role) bool {
			return item1.Id == item2.Id
		})
		removedRoles = slice.DifferenceWith(oldRoles, roles, func(item1, item2 *model.Role) bool {
			return item1.Id == item2.Id
		})
		innerErr = l.repos.OrganizationUserRepo.UpdateRolesWithTx(ctx, tx, organizationUser, newRoles, removedRoles)
		if innerErr != nil {
			return innerErr
		}
		organizationUser.Roles = roles

		return innerErr
	})
	if err != nil {
		return nil, err
	}
	// todo update user's permissions

	result = l.transToBo(organizationUser)

	return result, nil
}

// UpdateProfile update profile
func (l *OrganizationUserLogic) UpdateProfile(ctx context.Context, input dto.OrganizationUserProfileEditRequestDto) (err error) {
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserByUserId(ctx, input.OrganizationId, input.UserId)
	if err != nil {
		return err
	}
	organizationUser.FirstName = input.FirstName
	organizationUser.LastName = input.LastName
	organizationUser.MidName = input.MidName
	organizationUser.FullName = stringx.FormatFullName(input.FirstName, input.MidName, input.LastName)
	organizationUser.Phone = input.Phone
	_, err = l.repos.OrganizationUserRepo.UpdateOneById(ctx, organizationUser.Id, organizationUser)
	return err
}

type SendInviteUserEmailInput struct {
	EmailConfig dto.SendInviteUserEmailConfig
	FirstName   string
	LastName    string
	TargetEmail string
	UserType    constants.UserType
	// HasPwd           bool
	Pwd              string
	OrganizationName string
	OrganizationId   int64
	Operator         string
	Code             string
}

// sendInviteEmail send invite email
func (l *OrganizationUserLogic) sendInviteEmail(ctx context.Context, input SendInviteUserEmailInput) (err error) {

	sendTagKey := stringx.Format(constants.UserInviteEmailSendTag, input.UserType, input.TargetEmail, input.OrganizationId)
	isExisted, _ := l.redisCli.IsExisted(ctx, sendTagKey)
	if isExisted {
		err = errorx.New("The Email Has Been Sent")
		return err
	}
	defer func() {
		// set send tag
		_ = l.redisCli.SetExpire(ctx, sendTagKey, "1", input.EmailConfig.SendPeriodSeconds)
	}()

	var code string
	// generate invite code
	codeKey := stringx.Format(constants.UserInviteCode, input.UserType, input.TargetEmail, input.OrganizationId)
	code, _ = l.redisCli.Get(ctx, codeKey)
	if len(code) == 0 {
		code = random.RandNumeral(10)
		err = l.redisCli.SetExpire(ctx, codeKey, code, input.EmailConfig.ExpireSeconds)
		if err != nil {
			return errorx.WithStack(err)
		}
		input.Code = code
	}
	// save input
	if len(input.Pwd) > 0 {
		pwdKey := stringx.Format(constants.UserInviteUserPwd, input.UserType, input.TargetEmail, input.OrganizationId)
		// inputContent, _ := json.Marshal(input)
		err = l.redisCli.SetExpire(ctx, pwdKey, input.Pwd, input.EmailConfig.ExpireSeconds)
		if err != nil {
			return errorx.WithStack(err)
		}
	}

	// send invite email
	data := map[string]interface{}{
		"url": fmt.Sprintf("%s?code=%s&email=%s&first_name=%s&last_name=%s&organization_name=%s&organization_id=%v",
			input.EmailConfig.Url,
			code,
			input.TargetEmail,
			input.FirstName,
			input.LastName,
			input.OrganizationName,
			input.OrganizationId,
		),
		"has_pwd":           len(input.Pwd) > 0,
		"pwd":               input.Pwd,
		"operator":          input.Operator,
		"organization_name": input.OrganizationName,
	}
	content, err := email2.ParseTemplate(input.EmailConfig.TplName, data)
	if err != nil {
		return errorx.WithStack(err)
	}
	msg := email2.Message{
		Subject: fmt.Sprintf("%s has invited you to join %s", input.Operator, input.OrganizationName),
		From:    input.EmailConfig.HostEmail,
		To:      []string{input.TargetEmail},
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
	result = l.batchTransToBo(repoResults)
	return result, nil
}

// Detail
func (l *OrganizationUserLogic) Detail(ctx context.Context, id int64, organizationId int64) (result *bo.OrganizationUserBo, err error) {
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, id, organizationId)
	if err != nil {
		return nil, err
	}
	user := organizationUser.User
	if user == nil {
		err = errorx.New("User Is Not Existed")
		return nil, err
	}
	result = &bo.OrganizationUserBo{}
	err = copier.CopyWithOption(result, organizationUser, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	result.Email = user.Email
	return result, nil
}

// SwitchStatus
func (l *OrganizationUserLogic) SwitchStatus(ctx context.Context, id int64, organizationId int64, status constants.OrganizationUserStatus, userType constants.UserType) (err error) {
	// check organization user
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, organizationId, userType, false)
	if err != nil {
		return err
	}
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, id, organizationId)
	if err != nil {
		return err
	}

	if status != constants.OrganizationUserStatusActive && organization.OwnerId == organizationUser.UserId {
		return errorx.New("Cannot Disable Owner")
	}

	err = l.repos.OrganizationUserRepo.UpdateOrganizationUserStatus(ctx, id, organizationId, status)

	// todo update user's permissions
	return err
}

// AssignRoles assign roles for organization user
func (l *OrganizationUserLogic) AssignRoles(ctx context.Context, id int64, organizationId int64, roleIds []int64) (err error) {
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, id, organizationId)
	if err != nil {
		return err
	}
	err = l.repos.OrganizationUserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		roles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, organizationId, roleIds)
		if innerErr != nil {
			return innerErr
		}
		newRoles := []*model.Role{}
		removedRoles := []*model.Role{}
		newRoles = slice.DifferenceWith(roles, organizationUser.Roles, func(item1, item2 *model.Role) bool {
			return item1.Id == item2.Id
		})
		removedRoles = slice.DifferenceWith(organizationUser.Roles, roles, func(item1, item2 *model.Role) bool {
			return item1.Id == item2.Id
		})
		innerErr = l.repos.OrganizationUserRepo.UpdateRolesWithTx(ctx, tx, organizationUser, newRoles, removedRoles)

		return innerErr
	})
	// todo update user's permissions
	return err
}

func (l *OrganizationUserLogic) Delete(ctx context.Context, id int64, organizationId int64, userType constants.UserType) (err error) {
	// check organization user
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, organizationId, userType, false)
	if err != nil {
		return err
	}
	err = l.repos.OrganizationUserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		organizationUser, innerErr := l.repos.OrganizationUserRepo.GetOrganizationUserByIdWithTx(ctx, tx, id, organizationId)
		if innerErr != nil {
			return innerErr
		}
		if organization.OwnerId == organizationUser.UserId {
			innerErr = errorx.New("Cannot Delete Owner")
			return innerErr
		}
		innerErr = l.repos.OrganizationUserRepo.DeleteOrganizationUserWithTx(ctx, tx, organizationUser)
		return innerErr
	})
	// todo delete user's permissions
	return err
}

// transToBo
func (l *OrganizationUserLogic) transToBo(src *model.OrganizationUser) (result *bo.OrganizationUserBo) {
	if src == nil {
		return nil
	}
	result = &bo.OrganizationUserBo{}
	_ = copier.CopyWithOption(result, src, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters: []copier.TypeConverter{
			copierx.TimeToRFC3339Str(),
		},
	})
	if src.User != nil {
		result.Email = src.User.Email
	}
	return result
}

// batchTransToBo batch trans to bo
func (l *OrganizationUserLogic) batchTransToBo(src []*model.OrganizationUser) (result []*bo.OrganizationUserBo) {
	result = make([]*bo.OrganizationUserBo, 0, len(src))
	if len(src) == 0 {
		return result
	}
	for i := range src {
		result = append(result, l.transToBo(src[i]))
	}
	return result
}

func (l *OrganizationUserLogic) validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	return hasLower && hasUpper && hasNumber
}
