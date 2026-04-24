import { formatAmount } from '@repo/shared-methods/utils';

// 编辑3rd party loan的菜单导航------------------------start------------------------
export const createLoanSteps = [
	{
		title: 'Loan detail',
		key: 'loanDetail',
	},
	{
		title: 'Property detail',
		key: 'propertyDetail',
	},
	{
		title: 'Loan terms',
		key: 'loanTerms',
	},
];
// 编辑3rd party loan的菜单导航------------------------end------------------------

// 新建zeeBee loan的菜单导航------------------------start------------------------
export const createZeeBeeSteps = [
	{
		title: 'Borrower information',
		key: 'borrowerInformation',
		children: [
			{
				title: 'Applicant information',
				key: 'applicantInformation',
				validate: ['loan_borrower.'],
			},
			{
				title: 'Borrower information',
				key: 'borrowerInformation',
				validate: ['applicants.'],
			},
		],
	},
	{
		title: 'Property & loan',
		key: 'propertyAndLoan',
		children: [
			{
				title: 'Property information',
				key: 'propertyInformation ',
			},
			{
				title: 'Loan detail',
				key: 'loanDetail ',
			},
		],
	},
	{
		title: 'Experience & Third service',
		key: 'experienceAndThirdService',
		children: [
			{
				title: 'Experience',
				key: 'experience',
			},
			{
				title: 'Third-party services',
				key: 'thirdPartyServices',
			},
		],
	},
	{
		title: 'Declaration & Demographics',
		key: 'declarationAndDemographics',
		children: [
			{
				title: 'About this property and your money for this loan',
				key: 'propertyLoanMoney',
			},
			{
				title: 'Borrower demographic information',
				key: 'borrowerDemographic',
			},
		],
	},
	{
		title: 'Review and sign',
		key: 'reviewAndSign',
	},
];
// 新建zeeBee loan的菜单导航------------------------end------------------------

export const zeeBeeDetailsTabs = [
	{
		label: 'General',
		value: 'general',
	},
	{
		label: 'Documents',
		value: 'documents',
		totalKey: 'requested_loan_document_count',
	},
];

// -------------------------------------------------------------
// zeeBee Details --> General------------------------start------------------------
// -------------------------------------------------------------
export const generalTermsFields = [
	{
		label: 'Total loan amount',
		value: 'total_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Rate',
		value: 'rate',
		render: (v) => (v ? `${v}%` : '--'),
	},
	{
		label: 'Outstanding loan amount',
		value: 'principal_amount',
		render: (v) => formatAmount(v),
	},
	{ label: 'Maturity date', value: 'maturity_date' },
	{
		label: 'Monthly Payment',
		value: 'monthly_payment',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Term',
		value: 'loan_term',
		render: (v) => `${v || '--'} Months`,
	},
];

export const generalParametersFields = [
	{
		label: 'Lien position',
		value: 'lien_position',
	},
	{
		label: 'Loan type',
		value: 'loan_type',
	},
];

export const generalList = [
	{
		title: 'Terms',
		fields: generalTermsFields,
	},
	{
		title: 'Parameters',
		fields: generalParametersFields,
	},
];
// zeeBee Details --> General------------------------end------------------------
