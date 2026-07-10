export const PERMISSIONS = {
	SETTING_USER: {
		ADD: 'Setting:User:Add',
		EDIT: 'Setting:User:Edit',
		DELETE: 'Setting:User:Delete',
		CHANGE_STATUS: 'Setting:User:ChangeStatus',
		SEND_INVITATION: 'Setting:User:SendInvitation',
	},
	SETTING_ROLE: {
		ADD: 'Setting:Role:Add',
		EDIT: 'Setting:Role:Edit',
		DELETE: 'Setting:Role:Delete',
		ASSIGN: 'Setting:Role:Assign',
	},
	SETTING_SYSTEM: {
		PAGE: 'Setting:System',
		PRE_APPROVED_LETTER: 'Setting:System:SendPreApprovalLetter',
	},
	SETTING_PRE_APPROVED_LETTER: {
		SAVE: 'Setting:PreApprovedLetter:Save',
		DOWNLOAD: 'Setting:PreApprovedLetter:Download',
	},
	SETTING_ORGANIZATION_PROFILE: {
		SAVE: 'Setting:OrganizationProfile:Save',
	},
	LOAN: {
		DETAILS: 'Loan:Details',
		CHANGE_STAGE: 'Loan:ChangeStage',
		CHANGE_PENDING_STATUS: 'Loan:ChangePendingStatus',
		EDIT: 'Loan:Details:Edit',
		UNDERWRITING: 'Loan:Details:Underwriting',
		QUOTING: 'Loan:Details:Quoting',
		PROCESSING: 'Loan:Details:Processing',
		CLOSING: 'Loan:Details:Closing',
	},
};
