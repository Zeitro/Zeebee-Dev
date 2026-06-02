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
		// CONTACT_EMAIL: 'Setting:System:ContactEmail',
	},
	PERSONAL: {
		CHANGE_PASSWORD: 'Personal:ChangePassword',
	},
	GMAIL_AUTH: 'GmailAuth',
	MANAGEMENT_BORROWER: {
		ADD: 'Management:Borrower:Add',
		EDIT: 'Management:Borrower:Edit',
		DELETE: 'Management:Borrower:Delete',
	},
	MANAGEMENT_BORROWING_ENTITY: {
		ADD: 'Management:BorrowingEntity:Add',
		EDIT: 'Management:BorrowingEntity:Edit',
		DELETE: 'Management:BorrowingEntity:Delete',
	},
	PROPERTY: {
		ADD: 'Property:Add',
		EDIT: 'Property:Edit',
		DELETE: 'Property:Delete',
		DETAILS: 'Property:Details',
	},
	LOAN: {
		LOAN_PAGE: 'Loan',
		ADD_3RD_PARTY: 'Loan:Create3rdParty',
		ADD_ZEEBEE: 'Loan:Create',
		EDIT_3RD_PARTY: 'Loan:Edit3rdParty',
		EDIT_ZEEBEE_DRAFT: 'Loan:EditDraftZeeBee',
		DELETE: 'Loan:Delete',
		IMPORT: 'Loan:ImportLoans',
		DETAILS: 'Loan:Details',
		CHANGE_PENDING_STATUS: 'Loan:ChangePendingStatus',
		CHANGE_STAGE: 'Loan:ChangeStage',
		QUOTE: 'Loan:Quote',
	},
};
