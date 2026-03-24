package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/email/client"
	"zeebee-backend/domains/zeebee/model"
)

// OrganizationEmailLogic provides organization email business logic
type OrganizationEmailLogic struct {
	*BaseLogic
}

// NewOrganizationEmailLogic creates a new organization email logic instance
func NewOrganizationEmailLogic(baseLogic *BaseLogic) *OrganizationEmailLogic {
	return &OrganizationEmailLogic{
		BaseLogic: baseLogic,
	}
}

// SaveEmail saves or updates an organization's email configuration with the provided details.
// It handles both new email creation and existing email updates within a database transaction.
// The method also manages related email authorization records based on the auth strategy changes.
func (l *OrganizationEmailLogic) SaveEmail(ctx context.Context, input dto.OrganizationEmailSaveEmailRequestDto) (err error) {
	// Begin a database transaction to ensure atomicity of all operations
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error
		// var organization *model.Organization
		var oldAuthStrategy constants.EmailAuthStrategy
		var oldEmail string

		// Get and validate the organization exists
		_, innerErr = l.repos.OrganizationRepo.GetOrganizationWithTx(ctx, tx, input.OrganizationId, input.OrganizationType, false)
		if innerErr != nil {
			return innerErr
		}

		// Create a new OrganizationEmail model with the input data
		modelData := &model.OrganizationEmail{
			OrganizationId: input.OrganizationId,
			AuthStrategy:   input.AuthStrategy,
			Email:          input.Email,
			Purpose:        input.Purpose,
		}

		// Get or create the organization email record
		organizationEmail, isNew, innerErr := l.repos.OrganizationEmailRepo.GetOrCreateOrganizationEmailWithTx(ctx, tx, modelData)
		if innerErr != nil {
			return innerErr
		}
		oldAuthStrategy = organizationEmail.AuthStrategy
		oldEmail = organizationEmail.Email

		// Prepare columns to update in organization table
		saveColumns := map[string]interface{}{}
		if input.Purpose == constants.EmailPurposeContact {
			saveColumns["contact_email"] = input.Email
		}

		// Update organization email settings if there are columns to update
		if len(saveColumns) > 0 {
			innerErr = l.repos.OrganizationRepo.UpdateColumnsByIdWithTx(ctx, tx, input.OrganizationId, saveColumns)
			if innerErr != nil {
				return innerErr
			}
		}
		// Update email authorization record
		organizationEmail.Email = input.Email
		organizationEmail.AuthStrategy = input.AuthStrategy
		organizationEmail.Purpose = input.Purpose
		if !isNew {
			_, innerErr = l.repos.OrganizationEmailRepo.UpdateOneByIdWithTx(ctx, tx, organizationEmail.Id, organizationEmail)
		}
		// If auth strategy changed, delete old authorization record
		if input.AuthStrategy != oldAuthStrategy {
			innerErr = l.repos.EmailAuthorizationRepo.DeleteEmailAuthorizationByOrganizationEmailIdWithTx(ctx, tx, organizationEmail.Id, input.Purpose)
		}
		// Update email authorization's authorization status if email changed for organization auth strategy
		if !isNew && input.AuthStrategy == constants.EmailAuthStrategyOrganization && input.Email != oldEmail {
			emailAuthorization, _ := l.repos.EmailAuthorizationRepo.GetEmailAuthorizationByOrganizationEmailIdWithTx(ctx, tx, organizationEmail.Id, input.Purpose)
			if emailAuthorization != nil && emailAuthorization.Id > 0 {
				emailAuthorization.Email = input.Email
				emailAuthorization.IsAuthorized = false
				_, innerErr = l.repos.EmailAuthorizationRepo.UpdateOneByIdWithTx(ctx, tx, emailAuthorization.Id, emailAuthorization)
			}
		}
		return innerErr
	})
	return err
}

// GenerateAuthCodeUrl generates an authentication URL for an organization or user.
// It validates the organization status, retrieves email authorization based on the strategy,
// updates email provider information, and generates the final authentication URL.
func (l *OrganizationEmailLogic) GenerateAuthCodeUrl(ctx context.Context, input dto.OrganizationEmailGenerateAuthCodeUrlRequestDto) (url string, err error) {
	if input.AuthCodeGenerator == nil {
		return "", errorx.New("Auth code url generator is nil")
	}
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, input.OrganizationId, input.OrganizationType, false)
	if err != nil {
		return "", err
	}
	if organization.Status != constants.OrganizationStatusActive {
		return "", errorx.New("The organization is not active")
	}

	organizationEmail, err := l.repos.OrganizationEmailRepo.GetOrganizationEmailByOrganizationId(ctx, input.OrganizationId, input.Purpose)
	if err != nil {
		return "", err
	}

	var emailAuthorization *model.EmailAuthorization
	switch organizationEmail.AuthStrategy {
	case constants.EmailAuthStrategyOrganization:
		emailAuthorization, err = l.repos.EmailAuthorizationRepo.GetEmailAuthorizationByOrganizationEmailId(ctx, organizationEmail.Id, 0, organizationEmail.Purpose)
		if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
			return "", err
		}
		if emailAuthorization == nil || emailAuthorization.Id == 0 {
			if emailAuthorization == nil || emailAuthorization.Id == 0 {
				emailAuthorization = &model.EmailAuthorization{
					OrganizationEmailId: organizationEmail.Id,
					UserId:              0,
					Email:               organizationEmail.Email,
					Purpose:             organizationEmail.Purpose,
					EmailProvider:       input.EmailProvider,
					IsAuthorized:        false,
				}
				emailAuthorization, err = l.repos.EmailAuthorizationRepo.Create(ctx, emailAuthorization)
			}
		}
	case constants.EmailAuthStrategyUser:
		emailAuthorization, err = l.repos.EmailAuthorizationRepo.GetEmailAuthorizationByOrganizationEmailId(ctx, input.OrganizationId, input.UserId, organizationEmail.Purpose)
		if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
			return "", err
		}
		if emailAuthorization == nil || emailAuthorization.Id == 0 {
			var user *model.User
			user, err = l.repos.UserRepo.GetUser(ctx, input.UserId, input.OrganizationType, false)
			if err != nil {
				return "", err
			}
			if emailAuthorization == nil || emailAuthorization.Id == 0 {
				emailAuthorization = &model.EmailAuthorization{
					OrganizationEmailId: organizationEmail.Id,
					UserId:              input.UserId,
					Email:               user.Email,
					Purpose:             organizationEmail.Purpose,
					EmailProvider:       input.EmailProvider,
					IsAuthorized:        false,
				}
				emailAuthorization, err = l.repos.EmailAuthorizationRepo.Create(ctx, emailAuthorization)
			}
		}

	default:
		err = errorx.New("Invalid email auth strategy")
	}
	if err != nil {
		return "", err
	}
	if emailAuthorization == nil || emailAuthorization.Id == 0 {
		return "", errorx.New("Email authorization not found")
	}

	emailAuthorization.EmailProvider = input.EmailProvider
	// emailAuthorization.EmailAuthStrategy = organization.EmailAuthStrategy
	_, err = l.repos.EmailAuthorizationRepo.UpdateOneById(ctx, emailAuthorization.Id, emailAuthorization)
	if err != nil {
		return "", err
	}
	url = input.AuthCodeGenerator.GenerateAuthCodeUrl(input.Purpose)

	return url, nil
}

// EmailAuthorize handles the email authorization process for an organization.
// It supports different email providers (currently Gmail) and authorization strategies (organization-level and user-level).
// The function validates the email provider, retrieves the appropriate email authorization, verifies the authenticated email,
// refreshes the token, and updates the authorization status.
func (l *OrganizationEmailLogic) EmailAuthorize(ctx context.Context, input dto.OrganizationEmailAuthorizeRequestDto) (email string, emailClient client.Client, err error) {
	// var emailClient client.Client
	switch input.EmailProvider {
	case constants.EmailProviderGmail:
		// Create a new Gmail client with OAuth2 configuration and authorization code
		emailClient, err = client.NewGmailClient(input.Oauth2Config, input.Code)
	default:
		// Return error for unsupported email providers
		err = errorx.New("Invalid email provider")
	}
	if err != nil {
		return "", nil, err
	}

	// Get organization details for email authorization
	organizationEmail, err := l.repos.OrganizationEmailRepo.GetOrganizationEmailByOrganizationId(ctx, input.OrganizationId, input.Purpose)
	if err != nil {
		return
	}
	var emailAuthorization *model.EmailAuthorization
	// Handle different email authorization strategies
	switch organizationEmail.AuthStrategy {
	case constants.EmailAuthStrategyOrganization:
		// Get organization-level email authorization
		emailAuthorization, err = l.repos.EmailAuthorizationRepo.GetEmailAuthorizationByOrganizationEmailId(ctx, organizationEmail.Id, 0, input.Purpose)
	case constants.EmailAuthStrategyUser:
		// Get user-level email authorization
		emailAuthorization, err = l.repos.EmailAuthorizationRepo.GetEmailAuthorizationByOrganizationEmailId(ctx, organizationEmail.Id, input.UserId, input.Purpose)
	default:
		// Return error for invalid authorization strategy
		err = errorx.New("Invalid email auth strategy")
	}

	if err != nil {
		return "", nil, err
	}

	// Get the authenticated email from the email client
	authEmail, err := emailClient.GetEmail(ctx)
	if err != nil {
		return "", nil, err
	}

	// Verify that the authenticated email matches the stored authorization email
	if emailAuthorization.Email != authEmail {
		return "", nil, errorx.New("Email does not match")
	}

	// Refresh the authorization token
	token, err := emailClient.RefreshToken(ctx)
	if err != nil {
		return "", nil, err
	}

	// Update the authorization status and token
	emailAuthorization.IsAuthorized = true
	emailAuthorization.Token = &model.EmailToken{
		Token: token,
	}

	// Save the updated authorization to the repository
	_, err = l.repos.EmailAuthorizationRepo.UpdateOneById(ctx, emailAuthorization.Id, emailAuthorization)
	if err != nil {
		return "", nil, err
	}

	return authEmail, emailClient, nil
}
