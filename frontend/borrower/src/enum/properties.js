import { formatAmount, formatNumber } from '@repo/shared-methods/utils';

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
		render: (v) => (v ? `${v}%` : '--'),
	},
	{
		label: 'Total monthly payment',
		value: 'total_monthly_payment_amount',
		render: (v) => formatAmount(v),
	},
];

export const propertyList = [
	{
		label: 'Property type',
		value: 'property_type',
		className: 'w-2/3',
	},
	{
		label: 'Living Area',
		value: 'living_area_sqft',
		render: (v) => formatNumber(v),
		className: 'w-1/3',
	},
	{
		label: 'Units',
		value: 'unit',
		className: 'w-1/3',
	},
	{
		label: 'Beds',
		value: 'bedrooms',
		className: 'w-1/3',
	},
	{
		label: 'Baths',
		value: 'bathrooms',
		className: 'w-1/3',
	},
];
