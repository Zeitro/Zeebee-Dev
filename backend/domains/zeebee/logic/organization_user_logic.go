package logic

import (
	"context"
	"fmt"
	"regexp"

	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	email2 "zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/stringx"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/auth"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

// OrganizationUserLogic provides organization user business logic
type OrganizationUserLogic struct {
	*BaseLogic
	authTopicMap map[constants.UserType]string
	emailClient  *email2.Email
}

// NewOrganizationUserLogic creates a new organization user logic instance
func NewOrganizationUserLogic(baseLogic *BaseLogic, emailClient *email2.Email, authTopicMap map[constants.UserType]string) *OrganizationUserLogic {
	return &OrganizationUserLogic{
		BaseLogic:    baseLogic,
		authTopicMap: authTopicMap,
		emailClient:  emailClient,
	}
}

// Page returns organization user list
func (l *OrganizationUserLogic) Page(ctx context.Context, pageReq dto.OrganizationUserPageRequestDto) (page dto.PageResultDto[*dto.OrganizationUserEntityDto], err error) {
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
		roles, err := l.repos.RoleRepo.GetRolesByOrganizationUserId(ctx, entity.Id)
		if err != nil {
			return page, err
		}
		results.List[i].Roles = roles
	}
	pagePtr, err := copierx.Trans[gormx.PageResponse[*repo.OrganizationUserPageEntity], dto.PageResultDto[*dto.OrganizationUserEntityDto]](&results)
	if err != nil {
		return page, err
	}
	page = *pagePtr

	slice.ForEach(page.List, func(index int, item *dto.OrganizationUserEntityDto) {
		page.List[index].IsOwner = item.UserId == organizationOwnerId
	})
	return page, err
}

// CreateAdminOrganizationUser creates admin organization user
func (l *OrganizationUserLogic) CreateAdminOrganizationUser(ctx context.Context, input dto.OrganizationUserCreateAdminRequestDto) (result *dto.OrganizationUserEntityDto, err error) {
	isInOrganization, err := l.repos.OrganizationUserRepo.CheckOrganizationUserExistedByEmail(ctx, input.OrganizationId, input.Email, input.UserType)
	if err != nil {
		return nil, err
	}
	if isInOrganization {
		err = errorx.New("the user is already in the organization")
		return nil, err
	}
	user, err := l.repos.UserRepo.GetUserByEmail(ctx, input.Email, input.UserType)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if !l.validatePassword(input.Password) {
		err = errorx.New("the password is invalid")
		return nil, err
	}

	var organizationUser *model.OrganizationUser
	var organization *model.Organization
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// Get organization
		organization, innerErr = l.repos.OrganizationRepo.GetOrganizationWithTx(ctx, tx, input.OrganizationId, input.UserType, false)
		if innerErr != nil {
			return innerErr
		}
		if organization.Status != constants.OrganizationStatusActive {
			innerErr = errorx.New("the organization is not active")
			return innerErr
		}
		firstName, lastName, fullName := input.FirstName, input.LastName, stringx.FormatFullName(input.FirstName, "", input.LastName)
		// Create user if not exist
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
				FirstName:             firstName,
				LastName:              lastName,
				FullName:              fullName,
			}
			user, innerErr = l.repos.UserRepo.CreateWithTx(ctx, tx, user, false)
			if innerErr != nil {
				return innerErr
			}
		} else if user.DefaultOrganizationId == 0 {
			// Existing user without a default org: pin default to this org
			if innerErr = l.repos.UserRepo.UpdateColumnsByIdWithTx(ctx, tx, user.Id, map[string]interface{}{
				"default_organization_id": input.OrganizationId,
			}); innerErr != nil {
				return innerErr
			}
			user.DefaultOrganizationId = input.OrganizationId
		}

		// Create organization user
		organizationUser = &model.OrganizationUser{
			UserId:         user.Id,
			OrganizationId: input.OrganizationId,
			Phone:          input.Phone,
			Remark:         input.Remark,
			Status:         constants.OrganizationUserStatusActive,
			FirstName:      firstName,
			MidName:        "",
			LastName:       lastName,
			FullName:       fullName,
		}

		// Get organization roles
		if len(input.RoleIds) > 0 {
			roles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, input.OrganizationId, input.RoleIds, constants.UserTypeAdmin)
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

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganizationUser,
			Type:           organization.Type,
			OrganizationId: input.OrganizationId,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[organization.Type], msg)
	})

	result = l.transToBo(organizationUser)
	return result, nil
}

// InviteOrganizationUser invites organization user to join organization
func (l *OrganizationUserLogic) InviteOrganizationUser(ctx context.Context, createDto dto.OrganizationUserCreateRequestDto) (result *dto.OrganizationUserEntityDto, err error) {
	isInOrganization, err := l.repos.OrganizationUserRepo.CheckOrganizationUserExistedByEmail(ctx, createDto.OrganizationId, createDto.Email, createDto.UserType)
	if err != nil {
		return nil, err
	}
	if isInOrganization {
		err = errorx.New("the user is already in the organization")
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
		// Get organization
		organization, innerErr = l.repos.OrganizationRepo.GetOrganizationWithTx(ctx, tx, createDto.OrganizationId, createDto.UserType, false)
		if innerErr != nil {
			return innerErr
		}
		if organization.Status != constants.OrganizationStatusActive {
			innerErr = errorx.New("the organization is not active")
			return innerErr
		}

		// Create user if not exist
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

		// Create organization user
		organizationUser = &model.OrganizationUser{
			UserId:         user.Id,
			OrganizationId: createDto.OrganizationId,
			Phone:          createDto.Phone,
			Remark:         createDto.Remark,
			Status:         constants.OrganizationUserStatusNotActivated,
		}

		// Get organization roles
		if len(createDto.RoleIds) > 0 {
			roles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, createDto.OrganizationId, createDto.RoleIds, createDto.UserType)
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
	// Send email
	threading.GoSafe(func() {
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
			Logo:             organization.Logo,
			Address:          organization.Address,
		})
		if err != nil {
			err = errorx.Wrap(err, "send invite email error")
			logx.WithContext(ctx).Error(err)
		}
	})

	result = l.transToBo(organizationUser)

	return result, nil
}

// SendInviteEmail sends invite email to organization user
func (l *OrganizationUserLogic) SendInviteEmail(ctx context.Context, input dto.OrganizationUserSendInviteEmailRequestDto) (err error) {
	// Get organization user
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, input.Id, input.OrganizationId)
	if err != nil {
		return err
	}
	if organizationUser.Status != constants.OrganizationUserStatusNotActivated {
		err = errorx.New("the user status is not activated")
		return err
	}
	user := organizationUser.User
	if user == nil {
		err = errorx.New("user not found")
		return err
	}
	// Get organization
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, input.OrganizationId, input.UserType, false)
	if err != nil {
		return err
	}
	if organization.Status != constants.OrganizationStatusActive {
		err = errorx.New("the organization is not active")
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
		Logo:             organization.Logo,
		Address:          organization.Address,
	}
	err = l.sendInviteEmail(ctx, previousInput)
	if err != nil {
		err = errorx.Wrap(err, "send invite email error")
		return err
	}
	// go func() {
	// 	newCtx := context.WithoutCancel(ctx)
	// 	err := l.sendInviteEmail(newCtx, previousInput)
	// 	if err != nil {
	// 		err = errorx.Wrap(err, "send invite email error")
	// 		logx.WithContext(ctx).Error(err)
	// 	}
	// }()

	return nil
}

// AcceptInvitation accepts organization invitation to join organization
func (l *OrganizationUserLogic) AcceptInvitation(ctx context.Context, input dto.OrganizationUserAcceptInvitationRequestDto) (user dto.UserEntityDto, err error) {
	key := stringx.Format(constants.UserInviteCode, input.UserType, input.Email, input.OrganizationId)
	code, err := l.redisCli.Get(ctx, key)
	if err != nil {
		return user, errorx.WithStack(err)
	}
	if len(code) == 0 || code != input.Code {
		return user, errorx.New("invalid code")
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
			return errorx.New("organization user not found")
		}
		if organizationUserModel.Status != constants.OrganizationUserStatusNotActivated {
			return errorx.New("the user has been activated")
		}
		// Update organization user
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

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganizationUser,
			Type:           userModel.Type,
			OrganizationId: input.OrganizationId,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[userModel.Type], msg)
	})

	userPtr, err := copierx.Trans[model.User, dto.UserEntityDto](userModel, copierx.TimeToRFC3339Str())
	if err != nil {
		return user, err
	}

	return *userPtr, nil
}

// UpdateOrganizationUser updates organization user information
func (l *OrganizationUserLogic) UpdateOrganizationUser(ctx context.Context, input dto.OrganizationUserUpdateRequestDto) (result *dto.OrganizationUserEntityDto, err error) {
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, input.Id, input.OrganizationId)
	if err != nil {
		return nil, err
	}
	user := organizationUser.User
	err = copierx.Assign[dto.OrganizationUserUpdateRequestDto, model.OrganizationUser](&input, organizationUser)
	if err != nil {
		return nil, err
	}

	organizationUser.FullName = stringx.FormatFullName(organizationUser.FirstName, organizationUser.MidName, organizationUser.LastName)

	// var organizationUser *model.OrganizationUser
	var organization *model.Organization
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		oldRoles := organizationUser.Roles
		// Get organization
		organization, innerErr = l.repos.OrganizationRepo.GetOrganizationWithTx(ctx, tx, input.OrganizationId, input.UserType, false)
		if innerErr != nil {
			return innerErr
		}
		if organization.Status != constants.OrganizationStatusActive {
			innerErr = errorx.New("the organization is not active")
			return innerErr
		}
		// Update organization user
		organizationUser.RemoveRelationships()
		_, innerErr = l.repos.OrganizationUserRepo.UpdateOneWithTx(ctx, tx, organizationUser)
		if innerErr != nil {
			return innerErr
		}
		// Update roles
		roles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, input.OrganizationId, input.RoleIds, input.UserType)
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

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganizationUser,
			Type:           organization.Type,
			OrganizationId: input.OrganizationId,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[organization.Type], msg)
	})

	organizationUser.User = user
	result = l.transToBo(organizationUser)

	return result, nil
}

// UpdateProfile updates profile
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
	Logo        string
	Address     model.Address
	// HasPwd           bool
	Pwd              string
	OrganizationName string
	OrganizationId   int64
	Operator         string
	Code             string
}

// sendInviteEmail sends invite email
func (l *OrganizationUserLogic) sendInviteEmail(ctx context.Context, input SendInviteUserEmailInput) (err error) {
	sendTagKey := stringx.Format(constants.UserInviteEmailSendTag, input.UserType, input.TargetEmail, input.OrganizationId)
	isExisted, _ := l.redisCli.IsExisted(ctx, sendTagKey)
	if isExisted {
		err = errorx.New("the email has already been sent")
		return err
	}
	defer func() {
		// Set send tag
		_ = l.redisCli.SetExpire(ctx, sendTagKey, "1", input.EmailConfig.SendPeriodSeconds)
	}()

	var code string
	// Generate invite code
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
	// Save input
	if len(input.Pwd) > 0 {
		pwdKey := stringx.Format(constants.UserInviteUserPwd, input.UserType, input.TargetEmail, input.OrganizationId)
		// inputContent, _ := json.Marshal(input)
		err = l.redisCli.SetExpire(ctx, pwdKey, input.Pwd, input.EmailConfig.ExpireSeconds)
		if err != nil {
			return errorx.WithStack(err)
		}
	}

	// Send invite email
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
		"logo":              input.Logo,
		"address": fmt.Sprintf("%s | %s | %s, %s, %s",
			input.OrganizationName,
			input.Address.Address1,
			input.Address.City,
			input.Address.State,
			input.Address.ZipCode,
		),
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
func (l *OrganizationUserLogic) AllOrganizationUsers(ctx context.Context, organizationId int64) (result []*dto.OrganizationUserEntityDto, err error) {
	repoResults, err := l.repos.OrganizationUserRepo.AllOrganizationUsers(ctx, organizationId)
	if err != nil {
		return nil, err
	}
	result = l.batchTransToBo(repoResults)
	return result, nil
}

// Detail
func (l *OrganizationUserLogic) Detail(ctx context.Context, id int64, organizationId int64) (result *dto.OrganizationUserEntityDto, err error) {
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, id, organizationId)
	if err != nil {
		return nil, err
	}
	user := organizationUser.User
	if user == nil {
		err = errorx.New("user not existed")
		return nil, err
	}

	result = l.transToBo(organizationUser)
	return result, nil
}

// SwitchStatus
func (l *OrganizationUserLogic) SwitchStatus(ctx context.Context, id int64, organizationId int64, status constants.OrganizationUserStatus, userType constants.UserType) (err error) {
	// Check organization user
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, organizationId, userType, false)
	if err != nil {
		return err
	}
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, id, organizationId)
	if err != nil {
		return err
	}

	if status != constants.OrganizationUserStatusActive && organization.OwnerId == organizationUser.UserId {
		return errorx.New("cannot disable owner")
	}

	err = l.repos.OrganizationUserRepo.UpdateOrganizationUserStatus(ctx, id, organizationId, status)

	// Reload organization user auth
	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganizationUser,
			Type:           userType,
			OrganizationId: organizationId,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[userType], msg)
	})

	return err
}

// AssignRoles assign roles for organization user
func (l *OrganizationUserLogic) AssignRoles(ctx context.Context, id int64, organizationId int64, roleIds []int64, userType constants.UserType) (err error) {
	organizationUser, err := l.repos.OrganizationUserRepo.GetOrganizationUserById(ctx, id, organizationId)
	if err != nil {
		return err
	}
	err = l.repos.OrganizationUserRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		roles, innerErr := l.repos.RoleRepo.GetOrganizationRolesByIdsWithTx(ctx, tx, organizationId, roleIds, userType)
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

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganizationUser,
			Type:           userType,
			OrganizationId: organizationId,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[userType], msg)
	})

	return err
}

func (l *OrganizationUserLogic) Delete(ctx context.Context, id int64, organizationId int64, userType constants.UserType) (err error) {
	// Check organization user
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
			innerErr = errorx.New("cannot delete owner")
			return innerErr
		}
		innerErr = l.repos.OrganizationUserRepo.DeleteOrganizationUserWithTx(ctx, tx, organizationUser)
		return innerErr
	})

	threading.GoSafe(func() {
		subCtx := context.WithoutCancel(ctx)
		msg := auth.Msg{
			Action:         auth.ActionRoleReload,
			Scope:          auth.ScopeOrganizationUser,
			Type:           userType,
			OrganizationId: organizationId,
		}
		_ = l.publishAuthEvent(subCtx, l.authTopicMap[userType], msg)
	})

	return err
}

// transToBo
func (l *OrganizationUserLogic) transToBo(src *model.OrganizationUser) (result *dto.OrganizationUserEntityDto) {
	result, err := copierx.Trans[model.OrganizationUser, dto.OrganizationUserEntityDto](src)
	if err != nil {
		return nil
	}
	if src.User != nil {
		result.Email = src.User.Email
	}
	return result
}

// batchTransToBo batch transforms to bo
func (l *OrganizationUserLogic) batchTransToBo(src []*model.OrganizationUser) (result []*dto.OrganizationUserEntityDto) {
	result = make([]*dto.OrganizationUserEntityDto, 0, len(src))
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
