package constants

type OperationEntityType string

const (
	OperationEntityTypeBorrowerPropertyTransaction OperationEntityType = "BorrowerPropertyTransaction"
	OperationEntityTypeLoan                        OperationEntityType = "Loan"
	OperationEntityTypeUser                        OperationEntityType = "User"
	OperationEntityTypeRole                        OperationEntityType = "Role"
	OperationEntityTypeOrganizationUser            OperationEntityType = "OrganizationUser"
	OperationEntityTypeOrganizationSubscription    OperationEntityType = "OrganizationSubscription"
	OperationEntityTypePricingPlan                 OperationEntityType = "PricingPlan"
)

type OperationAction string

const (
	// Borrower property transaction actions
	OperationActionBorrowerPropertyTransactionCreate          OperationAction = "Transaction created"
	OperationActionBorrowerPropertyTransactionUpdateBaseInfo  OperationAction = "Basic information updated"
	OperationActionBorrowerPropertyTransactionUpdateStatus    OperationAction = "Status changed"
	OperationActionBorrowerPropertyTransactionUpdateChecklist OperationAction = "Checklist changed"
	OperationActionBorrowerPropertyTransactionAddDocument     OperationAction = "Document created"
	OperationActionBorrowerPropertyTransactionDelete          OperationAction = "Transaction deleted"
	OperationActionBorrowerPropertyTransactionUploadDocument  OperationAction = "Documents uploaded"
	OperationActionBorrowerPropertyTransactionDeleteDocument  OperationAction = "Document deleted"
	OperationActionBorrowerPropertyTransactionSyncToDrive     OperationAction = "Synced to Google Drive"

	// Loan actions
	OperationActionLoanCreate       OperationAction = "Loan Create"
	OperationActionLoanUpdate       OperationAction = "Loan Update"
	OperationActionLoanStatusChange OperationAction = "Loan Status Change"

	// Identity / security audit actions (NFR-003)
	OperationActionUserLogin         OperationAction = "User Login"
	OperationActionUserLoginFailed   OperationAction = "User Login Failed"
	OperationActionUserPasswordReset OperationAction = "User Password Reset"
	OperationActionUserPasswordEdit  OperationAction = "User Password Edit"

	// Role authorization actions (NFR-003: 角色授权变更)
	OperationActionRoleCreate       OperationAction = "Role Create"
	OperationActionRoleUpdate       OperationAction = "Role Update"
	OperationActionRoleDelete       OperationAction = "Role Delete"
	OperationActionRoleStatusSwitch OperationAction = "Role Status Switch"
	OperationActionRoleAssignUsers  OperationAction = "Role Assign Users"

	// Member / organization-user actions (NFR-003: 成员授权变更 / 邀请 / 加入)
	OperationActionMemberInvite           OperationAction = "Member Invite"
	OperationActionMemberInviteEmailSent  OperationAction = "Member Invite Email Sent"
	OperationActionMemberAcceptInvitation OperationAction = "Member Accept Invitation"
	OperationActionMemberCreate           OperationAction = "Member Create"
	OperationActionMemberDelete           OperationAction = "Member Delete"
	OperationActionMemberStatusSwitch     OperationAction = "Member Status Switch"
	OperationActionMemberUpdate           OperationAction = "Member Update"
	OperationActionMemberAssignRoles      OperationAction = "Member Assign Roles"

	// Subscription actions (NFR-003: 订阅变更)
	OperationActionSubscriptionCreate       OperationAction = "Subscription Create"
	OperationActionSubscriptionRenew        OperationAction = "Subscription Renew"
	OperationActionSubscriptionStatusSwitch OperationAction = "Subscription Status Switch"
	OperationActionSubscriptionSuspend      OperationAction = "Subscription Suspend"
	OperationActionSubscriptionExpireUpdate OperationAction = "Subscription Expire Update"
	OperationActionSubscriptionDefaultPlan  OperationAction = "Subscription Default Plan"

	// Pricing plan actions (NFR-003: 定价计划变更)
	OperationActionPlanCreate         OperationAction = "Plan Create"
	OperationActionPlanUpdate         OperationAction = "Plan Update"
	OperationActionPlanDelete         OperationAction = "Plan Delete"
	OperationActionPlanStatusSwitch   OperationAction = "Plan Status Switch"
	OperationActionPlanBindResources  OperationAction = "Plan Bind Resources"
)

type StorageType string

const (
	StorageTypeAwsS3       StorageType = "Aws S3"
	StorageTypeGoogleDrive StorageType = "Google Drive"
)

func (s StorageType) IsValid() bool {
	switch s {
	case StorageTypeAwsS3, StorageTypeGoogleDrive:
		return true
	}
	return false
}

type UploadType string

const (
	UploadTypeManual      UploadType = "Manual"
	UploadTypeGoogleDrive UploadType = "Google Drive"
)

func (s UploadType) IsValid() bool {
	switch s {
	case UploadTypeManual, UploadTypeGoogleDrive:
		return true
	}
	return false
}
