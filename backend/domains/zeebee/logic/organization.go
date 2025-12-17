package logic

import (
	"context"

	"gorm.io/gorm"

	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/bo"
	"zeebee-backend/domains/zeebee/data/dto"
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
// The function ensures data consistency by performing all operations within a single transaction.
// Note:
//   - The function performs the following operations within a transaction:
//     1. Creates or retrieves an EmailAuthorization record
//     2. Updates organization email settings if the scope is contact email
//     3. Updates authorization status if the email has changed
//   - All database operations are performed using the repository pattern with transaction support
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
		if input.EmailScope == constants.EmailScopeContact {
			saveColumns["contact_email"] = input.Email
			saveColumns["email_auth_strategy"] = input.EmailAuthStrategy
		}

		// Update organization email settings if there are columns to update
		if len(saveColumns) > 0 {
			innerErr = l.repos.OrganizationRepo.UpdateColumnsByIdWithTx(ctx, tx, input.OrganizationId, saveColumns)
			if innerErr != nil {
				return innerErr
			}
		}

		// If this is not a new record and the email has changed, update authorization status
		if !isNew && emailAuthorization.Email != input.Email {
			saveColumns = map[string]interface{}{
				"is_authorized": false,
			}
			innerErr = l.repos.EmailAuthorizationRepo.UpdateColumnsByIdWithTx(ctx, tx, emailAuthorization.Id, saveColumns)
			if innerErr != nil {
				return innerErr
			}
		}

		return innerErr
	})

	return err
}
