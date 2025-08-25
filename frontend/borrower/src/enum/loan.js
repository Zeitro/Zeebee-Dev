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
			},
			{
				title: 'Borrower information',
				key: 'borrowerInformation',
			},
		],
	},
	{
		title: 'Property and loan',
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
			paramsKey: 'otherHispanic',
			eventType: 'input',
			placeholder: 'Enter ethnicity, for example: Colombian, Dominican, Nicaraguan and so on',
		},
	},
	{ label: 'Not Hispanic or Latino', value: 'Not Hispanic or Latino' },
	{
		label: 'I do not wish to provide this information',
		value: 'I do not wish to provide this information',
	},
];

export const raceOptions = [
	{
		label: 'American Indian or Alaska Native',
		value: 'American Indian or Alaska Native',
		extras: {
			paramsKey: 'americanIndianOrAlaskaNative',
			eventType: 'input',
			placeholder: 'Name of enrolled or principal tribe',
		},
	},
	{
		label: 'Asian',
		value: 'Asian',
		extras: {
			eventType: 'checkbox',
			paramsKey: 'asian',
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
						paramsKey: 'otherAsian',
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
			paramsKey: 'otherPacific',
			eventType: 'input',
			placeholder: 'Enter race, For example: Fijian, Tongan and so on',
		},
	},
	{ label: 'White', value: 'White' },
	{
		label: 'I do not wish to provide this information',
		value: 'I do not wish to provide this information',
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
