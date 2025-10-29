import { formatAmount } from '@/lib/format';

// 创建或编辑loan时的步骤
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

export const tinOptions = [
	{ label: 'SSN', value: 'SSN' },
	{ label: 'ITIN', value: 'ITIN' },
	{ label: 'EIN', value: 'EIN' },
];

export const ethnicityOptions = [
	{ label: 'Hispanic or Latino', value: 'Hispanic or Latino' },
	{ label: 'Mexican', value: 'Mexican' },
	{ label: 'Cuban', value: 'Cuban' },
	{
		label: 'Other Hispanic or Latino',
		value: 'Other Hispanic or Latino',
		extras: {
			paramsKey: 'other_hispanic_or_latino_name',
			eventType: 'input',
			placeholder: 'Enter ethnicity, for example: Colombian, Dominican, Nicaraguan and so on',
		},
	},
	{ label: 'Not Hispanic or Latino', value: 'Not Hispanic or Latino' },
	{
		label: 'I do not wish to provide this information',
		value: 'Hidden',
		isRadio: true,
	},
];

export const raceOptions = [
	{
		label: 'American Indian or Alaska Native',
		value: 'American Indian or Alaska Native',
		extras: {
			paramsKey: 'american_indian_or_alaska_native_name',
			eventType: 'input',
			placeholder: 'Name of enrolled or principal tribe',
		},
	},
	{
		label: 'Asian',
		value: 'Asian',
		extras: {
			eventType: 'checkbox',
			paramsKey: 'asian_values',
			options: [
				{
					label: 'Asian Indian',
					value: 'Asian Indian',
				},
				{
					label: 'Chinese',
					value: 'Chinese',
				},
				{
					label: 'Filipino',
					value: 'Filipino',
				},
				{
					label: 'Japanese',
					value: 'Japanese',
				},
				{
					label: 'Korean',
					value: 'Korean',
				},
				{
					label: 'Vietnamese',
					value: 'Vietnamese',
				},
				{
					label: 'Other Asian',
					value: 'Other Asian',
					extras: {
						paramsKey: 'asian_other_name',
						eventType: 'input',
						placeholder: 'Enter race, For example: Hmong, Lootian, Thai and so on',
					},
				},
			],
		},
	},
	{ label: 'Black or African American', value: 'Black or African American' },
	{
		label: 'Native Hawaiian or Other Pacific Islander',
		value: 'Native Hawaiian or Other Pacific Islander',
	},
	{
		label: 'Native Hawaiian Guamanian or Chamorro Samoan',
		value: 'Native Hawaiian Guamanian or Chamorro Samoan',
	},
	{
		label: 'Other Pacific Islander',
		value: 'Other Pacific Islander',
		extras: {
			paramsKey: 'other_pacific_islander_name',
			eventType: 'input',
			placeholder: 'Enter race, For example: Fijian, Tongan and so on',
		},
	},
	{ label: 'White', value: 'White' },
	{
		label: 'I do not wish to provide this information',
		value: 'Hidden',
		isRadio: true,
	},
];

// loan 在pending状态下的可选状态
export const loanPendingStatus = [
	{
		label: 'Quoting',
		value: 'quoting',
	},
	{
		label: 'Processing',
		value: 'processing',
	},
	{
		label: 'Closing',
		value: 'closing',
	},
];

export const stageColor = {
	Pending: 'var(--base-primary)',
	Active: 'var(--base-success)',
	Archived: 'var(--base-error)',
};

export const pendingStatusColor = {
	Quoting: {
		color: 'var(--green-medium)',
		bg: 'var(--green-light)',
		border: 'var(--green-light-100)',
	},
	Processing: {
		color: 'var(--blue-dark)',
		bg: 'var(--blue-light-2)',
		border: 'var(--blue-light)',
	},
	UnderWriting: {
		color: '#eb2f96',
		bg: '#fff0f6',
		border: '#ffadd2',
	},
	Closing: {
		color: 'var(--purple-dark)',
		bg: 'var(--purple-light)',
		border: 'var(--purple-medium)',
	},
};

// zeeBee Details字段显示
export const pendingFields = [
	{ label: 'ID', value: 'code' },
	{ label: 'Target closing date', value: 'closing_date' },
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
];

export const activeFields = [
	{ label: 'ID', value: 'code' },
	{ label: 'Maturity date', value: 'maturity_date' },
	{
		label: 'Total loan amount',
		value: 'total_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Outstanding loan amount',
		value: 'total_principal_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Rate',
		value: 'rate',
		render: (v) => (v ? `${v}%` : '--'),
	},
];

export const archivedFields = [
	{ label: 'ID', value: 'code' },
	{ label: 'Paid off date', value: 'paid_off_date' },
	{
		label: 'Total loan amount',
		value: 'total_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Principal balance:',
		value: 'principal_amount:',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Rate',
		value: 'rate',
		render: (v) => (v ? `${v}%` : '--'),
	},
];

export const loanDetailsFields = {
	Pending: pendingFields,
	Active: activeFields,
	Archived: archivedFields,
};

// zeeBee Details --> General
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
// loan status
// export const loanStatusOptions = [
// 	{
// 		label: 'Pending',
// 		value: 'pending',
// 	},
// 	{
// 		label: 'Active',
// 		value: 'active',
// 	},
// 	{
// 		label: 'Archived',
// 		value: 'archived',
// 	},
// ];

// // loan purpose
// export const loanPurposeOptions = [
// 	{
// 		label: 'Purchase',
// 		value: 'purchase',
// 	},
// 	{
// 		label: 'Rate-term refi',
// 		value: 'rateTermRefi',
// 	},
// 	{
// 		label: 'Cash-out refi',
// 		value: 'cashOutRefi',
// 	},
// ];
// // loan type
// export const loanTypeOptions = [
// 	{
// 		label: 'Bridge',
// 		value: 'bridge',
// 	},
// 	{
// 		label: 'Rehab',
// 		value: 'rehab',
// 	},
// 	{
// 		label: 'Ground up construction',
// 		value: 'groundUpConstruction',
// 	},
// ];
