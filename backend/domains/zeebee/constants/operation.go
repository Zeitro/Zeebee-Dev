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
	OperationEntityTypePricingPlanPolicy           OperationEntityType = "PricingPlanPolicy"

	OperationEntityTypeProjectManagementOrganization OperationEntityType = "ProjectManagementOrganization"
	OperationEntityTypeProjectManagementProject      OperationEntityType = "ProjectManagementProject"
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

	// Identity / security audit actions
	OperationActionUserLogin         OperationAction = "User Login"
	OperationActionUserLoginFailed   OperationAction = "User Login Failed"
	OperationActionUserPasswordReset OperationAction = "User Password Reset"
	OperationActionUserPasswordEdit  OperationAction = "User Password Edit"

	// Role authorization actions
	OperationActionRoleCreate       OperationAction = "Role Create"
	OperationActionRoleUpdate       OperationAction = "Role Update"
	OperationActionRoleDelete       OperationAction = "Role Delete"
	OperationActionRoleStatusSwitch OperationAction = "Role Status Switch"
	OperationActionRoleAssignUsers  OperationAction = "Role Assign Users"

	// Member / organization-user actions
	OperationActionMemberInvite           OperationAction = "Member Invite"
	OperationActionMemberInviteEmailSent  OperationAction = "Member Invite Email Sent"
	OperationActionMemberAcceptInvitation OperationAction = "Member Accept Invitation"
	OperationActionMemberCreate           OperationAction = "Member Create"
	OperationActionMemberDelete           OperationAction = "Member Delete"
	OperationActionMemberStatusSwitch     OperationAction = "Member Status Switch"
	OperationActionMemberUpdate           OperationAction = "Member Update"
	OperationActionMemberAssignRoles      OperationAction = "Member Assign Roles"

	// Subscription actions
	OperationActionSubscriptionCreate       OperationAction = "Subscription Create"
	OperationActionSubscriptionRenew        OperationAction = "Subscription Renew"
	OperationActionSubscriptionStatusSwitch OperationAction = "Subscription Status Switch"
	OperationActionSubscriptionSuspend      OperationAction = "Subscription Suspend"
	OperationActionSubscriptionExpireUpdate OperationAction = "Subscription Expire Update"
	OperationActionSubscriptionDefaultPlan  OperationAction = "Subscription Default Plan"

	// Pricing plan actions
	OperationActionPlanCreate        OperationAction = "Plan Create"
	OperationActionPlanUpdate        OperationAction = "Plan Update"
	OperationActionPlanDelete        OperationAction = "Plan Delete"
	OperationActionPlanStatusSwitch  OperationAction = "Plan Status Switch"
	OperationActionPlanBindResources OperationAction = "Plan Bind Resources"
	OperationActionPlanPolicyCreate  OperationAction = "Plan Policy Create"
	OperationActionPlanPolicyUpdate  OperationAction = "Plan Policy Update"
	OperationActionPlanPolicyDelete  OperationAction = "Plan Policy Delete"

	// project management
	OperationActionProjectManagementTaskTypeCreate           OperationAction = "Project Task Type Created"
	OperationActionProjectManagementTaskTypeUpdate           OperationAction = "Project Task Type Updated"
	OperationActionProjectManagementTaskTypeDelete           OperationAction = "Project Task Type Deleted"
	OperationActionProjectManagementWorkflowCreate           OperationAction = "Project Workflow Created"
	OperationActionProjectManagementWorkflowUpdate           OperationAction = "Project Workflow Updated"
	OperationActionProjectManagementWorkflowDelete           OperationAction = "Project Workflow Deleted"
	OperationActionProjectManagementWorkflowStageCreate      OperationAction = "Project Workflow Stage Created"
	OperationActionProjectManagementWorkflowStageUpdate      OperationAction = "Project Workflow Stage Updated"
	OperationActionProjectManagementWorkflowStageDelete      OperationAction = "Project Workflow Stage Deleted"
	OperationActionProjectManagementWorkflowTransitionCreate OperationAction = "Project Workflow Transition Created"
	OperationActionProjectManagementWorkflowTransitionUpdate OperationAction = "Project Workflow Transition Updated"
	OperationActionProjectManagementWorkflowTransitionDelete OperationAction = "Project Workflow Transition Deleted"
	OperationActionProjectManagementProjectCreate            OperationAction = "Project Created"
	OperationActionProjectManagementProjectUpdate            OperationAction = "Project Updated"
	OperationActionProjectManagementProjectDelete            OperationAction = "Project Deleted"
	OperationActionProjectManagementMemberAdd                OperationAction = "Project Member Added"
	OperationActionProjectManagementMemberUpdate             OperationAction = "Project Member Role Updated"
	OperationActionProjectManagementMemberRemove             OperationAction = "Project Member Removed"
	OperationActionProjectManagementOwnerTransfer            OperationAction = "Project Owner Transferred"
	OperationActionProjectManagementTaskCreate               OperationAction = "Project Task Created"
	OperationActionProjectManagementTaskUpdate               OperationAction = "Project Task Updated"
	OperationActionProjectManagementTaskDelete               OperationAction = "Project Task Deleted"
	OperationActionProjectManagementTaskAssigneeAdd          OperationAction = "Project Task Assignee Added"
	OperationActionProjectManagementTaskAssigneeRemove       OperationAction = "Project Task Assignee Removed"
	OperationActionProjectManagementAttachmentAdd            OperationAction = "Project Task Attachment Added"
	OperationActionProjectManagementAttachmentDelete         OperationAction = "Project Task Attachment Deleted"
	OperationActionProjectManagementDependencyCreate         OperationAction = "Project Task Dependency Created"
	OperationActionProjectManagementDependencyUpdate         OperationAction = "Project Task Dependency Updated"
	OperationActionProjectManagementDependencyDelete         OperationAction = "Project Task Dependency Deleted"
	OperationActionProjectManagementTaskTransition           OperationAction = "Project Task Transitioned"
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
