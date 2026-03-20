export const asc = 'asc';
export const des = 'desc';
export const ascend = 'ascend';
export const descend = 'descend';
export const loan_pending = 'Pending';
export const loan_active = 'Active';
export const loan_archived = 'Archived';
// contact email验证状态颜色
export const authorizedColor = {
	authorized: {
		color: 'var(--green-medium)',
		bg: 'var(--green-light)',
		border: 'var(--green-light-100)',
	},
	unauthorized: {
		color: 'var(--orange-dark)',
		bg: 'var(--orange-light)',
		border: 'var(--orange-medium)',
	},
};

export const titleCol = {
	xs: 24,
	sm: 24,
};

export const layoutCol = {
	...titleCol,
	md: 12,
};

export const yesNoOptions = [
	{ value: true, label: 'Yes' },
	{ value: false, label: 'No' },
];

// 系统用户状态颜色
export const userStatusColor = {
	Active: {
		border: '#B7EB8F',
		bg: '#F6FFED',
		color: '#52C41A',
	},
	'Not activated': {
		border: '#FFD591',
		bg: '#FFF7E6',
		color: '#FA8C16',
	},
	Disabled: {
		border: 'var(--base-border)',
		bg: '#FAFAFA',
		color: 'rgba(0, 0, 0, 0.85)',
	},
};

export const systemTypeOptions = [
	{ value: 'Admin', label: 'Admin' },
	{ value: 'Lender', label: 'Lender' },
	{ value: 'Borrower', label: 'Borrower' },
];
