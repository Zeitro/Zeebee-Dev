package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/email/client"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/repo"
)

type OrganizationLogic struct {
	repos *repo.AllRepos
}

// NewOrganizationLogic creates a new instance of OrganizationLogic.
// It initializes the logic with the provided repository dependencies.
func NewOrganizationLogic(repos *repo.AllRepos) *OrganizationLogic {
	return &OrganizationLogic{
		repos: repos,
	}
}

// Detail retrieves organization details by ID and type.
// It fetches the organization from the repository and transforms it to the business object format.
func (l *OrganizationLogic) Detail(ctx context.Context, id int64, organizationType constants.UserType) (result *bo.OrganizationBo, err error) {
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, id, organizationType, true)
	if err != nil {
		return
	}
	result, err = copierx.Trans[model.Organization, bo.OrganizationBo](organization, copierx.TimeToRFC3339Str())
	if err != nil {
		return
	}
	return result, nil
}

// SaveEmail saves email authorization information for an organization.
// It handles the creation or update of email authorization records and organization email settings within a database transaction.
func (l *OrganizationLogic) SaveEmail(ctx context.Context, input dto.OrganizationSaveEmailRequestDto) (err error) {
	// Begin a database transaction to ensure atomicity of all operations
	err = l.repos.OrganizationRepo.Transaction(ctx, func(tx *gorm.DB) error {
		var innerErr error

		// Create a new EmailAuthorization model with the input data
		modelData := &model.EmailAuthorization{
			OrganizationId:    input.OrganizationId,
			UserId:            input.UserId,
			EmailAuthStrategy: input.EmailAuthStrategy,
			Email:             input.Email,
			EmailScope:        input.EmailScope,
		}

		// Get or create the email authorization record within the transaction
		emailAuthorization, isNew, innerErr := l.repos.EmailAuthorizationRepo.GetOrCreateEmailAuthorizationWithTx(ctx, tx, modelData)
		if innerErr != nil {
			return innerErr
		}

		// Prepare update columns for organization if the email scope is contact
		saveColumns := map[string]interface{}{}
		saveColumns["email_auth_strategy"] = input.EmailAuthStrategy
		if input.EmailScope == constants.EmailScopeContact {
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
		emailAuthorization.Email = input.Email
		emailAuthorization.EmailAuthStrategy = input.EmailAuthStrategy
		emailAuthorization.EmailScope = input.EmailScope
		if !isNew && emailAuthorization.Email != input.Email {
			emailAuthorization.IsAuthorized = false
		}
		if !isNew {
			_, innerErr = l.repos.EmailAuthorizationRepo.UpdateOneByIdWithTx(ctx, tx, emailAuthorization.Id, emailAuthorization)
		}

		return innerErr
	})

	return err
}

// GenerateAuthCodeUrl generates an authentication URL for an organization or user.
// It validates the organization status, retrieves email authorization based on the strategy,
// updates email provider information, and generates the final authentication URL.
func (l *OrganizationLogic) GenerateAuthCodeUrl(ctx context.Context, input dto.OrganizationGenerateAuthCodeUrlRequestDto) (url string, err error) {
	if input.AuthCodeGenerator == nil {
		return "", errorx.New("Auth Code Url Generator is Nil")
	}
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, input.OrganizationId, input.OrganizationType, false)
	if err != nil {
		return "", err
	}
	if organization.Status != constants.OrganizationStatusActive {
		return "", errorx.New("The Organization Is Not Active")
	}
	var emailAuthorization *model.EmailAuthorization
	switch organization.EmailAuthStrategy {
	case constants.EmailAuthStrategyOrganization:
		emailAuthorization, err = l.repos.EmailAuthorizationRepo.GetEmailAuthorization(ctx, input.OrganizationId, 0, input.EmailScope)
	case constants.EmailAuthStrategyUser:
		emailAuthorization, err = l.repos.EmailAuthorizationRepo.GetEmailAuthorization(ctx, input.OrganizationId, input.UserId, input.EmailScope)
		// err = errorx.New("Not Support The Strategy Yet")
	default:
		err = errorx.New("Invalid Email Auth Strategy")

	}
	if err != nil {
		return "", err
	}
	if emailAuthorization == nil || emailAuthorization.Id == 0 {
		return "", errorx.New("Email Authorization Not Found")
	}

	emailAuthorization.EmailProvider = input.EmailProvider
	emailAuthorization.EmailAuthStrategy = organization.EmailAuthStrategy
	_, err = l.repos.EmailAuthorizationRepo.UpdateOneById(ctx, emailAuthorization.Id, emailAuthorization)
	if err != nil {
		return "", err
	}
	url = input.AuthCodeGenerator.GenerateAuthCodeUrl(input.EmailScope)

	return url, nil
}

// EmailAuthorize handles email authorization for organizations.
// It creates an email client based on the provider, validates the email authorization,
// and updates the authorization status and token.
func (l *OrganizationLogic) EmailAuthorize(ctx context.Context, input dto.OrganizationEmailAuthorizeRequestDto) (email string, emailClient client.Client, err error) {
	// var emailClient client.Client
	switch input.EmailProvider {
	case constants.EmailProviderGmail:
		// Create a new Gmail client with OAuth2 configuration and authorization code
		emailClient, err = client.NewGmailClient(input.Oauth2Config, input.Code)
	default:
		// Return error for unsupported email providers
		err = errorx.New("Invalid Email Provider")
	}
	if err != nil {
		return "", nil, err
	}

	// Get organization details for email authorization
	organization, err := l.repos.OrganizationRepo.GetOrganization(ctx, input.OrganizationId, input.OrganizationType, false)
	if err != nil {
		return
	}
	var emailAuthorization *model.EmailAuthorization
	// Handle different email authorization strategies
	switch organization.EmailAuthStrategy {
	case constants.EmailAuthStrategyOrganization:
		// Get organization-level email authorization
		emailAuthorization, err = l.repos.EmailAuthorizationRepo.GetEmailAuthorization(ctx, input.OrganizationId, 0, input.EmailScope)
	case constants.EmailAuthStrategyUser:
		// Get user-level email authorization
		emailAuthorization, err = l.repos.EmailAuthorizationRepo.GetEmailAuthorization(ctx, input.OrganizationId, input.UserId, input.EmailScope)
	default:
		// Return error for invalid authorization strategy
		err = errorx.New("Invalid Email Auth Strategy")
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
		return "", nil, errorx.New("Email Not Match")
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
