import {HomeOutlined} from '@ant-design/icons';
import {formatAmount} from '@repo/shared-methods/utils';
// start------dashboard-----
export const loansSourceTypeOptions = [
	{value: 'All', label: 'All loans'},
	{value: 'ZeeBee', label: 'ZeeBee loans'},
	{value: '3rdParty', label: '3rd party loans'},
];
export const maturityAlertsStatusColor = {
	Overdue: {
		color: '#F5222D',
		background: '#FFF1F0',
		borderColor: '#FFA39E',
	},
	'Due Soon': {
		color: '#FA8C16',
		background: '#FFF7E6',
		borderColor: '#FFD591',
	},
};
export const filterTimeCountTips =
	'Based on the selected number of days, perform a global page filter for loans that match this time period.';

// loan最近天数筛选
export const filterTimeCount = [
	{value: 'All', label: 'All dates'},
	{value: 7, label: 'Last 7 days'},
	{value: 15, label: 'Last 15 days'},
	{value: 30, label: 'Last 30 days'},
];
// end------dashboard-----
// start------create zeeBee loan/loan details-----
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

// loan details

export const createZeeBeeSteps = [
	{
		title: 'Borrower information',
		key: 'borrowerInformation',
		// children: [
		// 	{
		// 		title: 'Applicant information',
		// 		key: 'applicantInformation',
		// 		validate: ['loan_borrower.'],
		// 	},
		// 	{
		// 		title: 'Borrower information',
		// 		key: 'borrowerInformation',
		// 		validate: ['applicants.'],
		// 	},
		// ],
	},
	{
		title: 'Property & loan',
		key: 'propertyAndLoan',
		// children: [
		// 	{
		// 		title: 'Property information',
		// 		key: 'propertyInformation ',
		// 	},
		// 	{
		// 		title: 'Loan detail',
		// 		key: 'loanDetail ',
		// 	},
		// ],
	},
	{
		title: 'Experience & Third service',
		key: 'experienceAndThirdService',
		// children: [
		// 	{
		// 		title: 'Experience',
		// 		key: 'experience',
		// 	},
		// 	{
		// 		title: 'Third-party services',
		// 		key: 'thirdPartyServices',
		// 	},
		// ],
	},
	{
		title: 'Declaration & Demographics',
		key: 'declarationAndDemographics',
		// children: [
		// 	{
		// 		title: 'About this property and your money for this loan',
		// 		key: 'propertyLoanMoney',
		// 	},
		// 	{
		// 		title: 'Borrower demographic information',
		// 		key: 'borrowerDemographic',
		// 	},
		// ],
	},
	{
		title: 'Review and sign',
		key: 'reviewAndSign',
	},
];

export const pendingSteps = [
	{
		title: 'Quoting',
		key: 'Quoting',
	},
	{
		title: 'Processing',
		key: 'Processing',
	},
	{
		title: 'Underwriting',
		key: 'Underwriting',
	},
	{
		title: 'Closing',
		key: 'Closing',
	},
];

export const lenderDetailsSteps = [
	{
		title: 'GENERAL',
		key: 'General',
		icon: <HomeOutlined />,
	},
	{
		title: 'APPLICATION',
		key: 'Application',
		children: createZeeBeeSteps,
	},
	{
		title: 'PENDING',
		key: 'Pending',
		children: pendingSteps,
	},
	{
		title: 'FUNDED',
		key: 'Funded',
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
	},
];

export const stageColor = {
	Pending: 'var(--base-primary)',
	Active: 'var(--base-success)',
	Archived: 'var(--base-error)',
	Draft: 'var(--base-warning)',
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
	Underwriting: {
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
	Draft: pendingFields,
};
