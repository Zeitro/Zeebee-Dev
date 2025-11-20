import { formatAmount } from '@repo/shared-methods/utils';
export const sort_asc = 'Ascending';
export const sort_des = 'Descending';

export const propertiesFieldArr = [
	{
		label: 'Total active loan',
		value: 'total_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Outstanding loan',
		value: 'total_principal_amount',
		render: (v) => formatAmount(v),
	},
	{
		label: 'Average active rate',
		value: 'avg_rate',
		render: (v) => `${v || 0}%`,
	},
	{
		label: 'Total monthly payment',
		value: 'total_monthly_payment_amount',
		render: (v) => formatAmount(v),
	},
];
