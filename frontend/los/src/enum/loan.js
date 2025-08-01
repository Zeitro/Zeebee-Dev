// 创建或编辑loan时的步骤
export const createLoanSteps = [
	{
		title: 'Loan detail',
		key:'loanDetail',
	},
	{
		title: 'Property detail',
		key:'propertyDetail'
	},
	{
		title: 'Loan terms',
		key:'loanTerms'
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
export const loanStatusOptions = [
	{
		label: 'Pending',
		value: 'pending',
	},
	{
		label: 'Active',
		value: 'active',
	},
	{
		label: 'Archived',
		value: 'archived',
	},
];

// loan purpose
export const loanPurposeOptions = [
	{
		label: 'Purchase',
		value: 'purchase',
	},
	{
		label: 'Rate-term refi',
		value: 'rateTermRefi',
	},
	{
		label: 'Cash-out refi',
		value: 'cashOutRefi',
	},
];
// loan type
export const loanTypeOptions = [
	{
		label: 'Bridge',
		value: 'bridge',
	},
	{
		label: 'Rehab',
		value: 'rehab',
	},
	{
		label: 'Ground up construction',
		value: 'groundUpConstruction',
	},
];
