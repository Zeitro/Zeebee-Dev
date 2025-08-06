import {formatAmount} from '@/lib/common';
export const sort_asc = 'Ascending';
export const sort_des = 'Descending';

export const propertiesFieldArr = [
	{
		label: 'Total loan',
		value: 'total_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Outstanding loan',
		value: 'total_principal_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Avg. rate',
		value: 'avg_rate',
		render: (v) => `${v}%`,
	},
	{
		label: 'Total monthly payment',
		value: 'total_monthly_payment',
		render: (v) => formatAmount(v),
	},
];
