import { formatAmount, formatNumber } from '@repo/shared-methods/utils';

export const ADDRESS_FIELDS = [
	{ label: 'Street 1', getValue: (data) => data?.address?.address1 },
	{ label: 'Street 2', getValue: (data) => data?.address?.address2 },
	{ label: 'City', getValue: (data) => data?.address?.city },
	{ label: 'State', getValue: (data) => data?.address?.state },
	{ label: 'Zip code', getValue: (data) => data?.address?.zip_code },
];

export const PROPERTY_FIELDS = [
	{ label: 'Property type', getValue: (data) => data?.property_type },
	{
		label: 'Living Area',
		getValue: (data) => {
			const v = data?.living_area_sqft;
			return v != null && v !== '' ? `${formatNumber(v)} sqft` : null;
		},
	},
	{ label: 'Units', getValue: (data) => data?.unit },
	{ label: 'Beds', getValue: (data) => data?.bedrooms },
	{ label: 'Baths', getValue: (data) => data?.bathrooms },
];

export const LOAN_INFO_FIELDS = [
	{ label: 'Total loan', getValue: (data) => formatAmount(data?.total_amount) },
	{
		label: 'Outstanding loan',
		getValue: (data) => formatAmount(data?.total_principal_amount),
	},
	{
		label: 'Avg. rate',
		getValue: (data) => (data?.avg_rate ? `${data.avg_rate}%` : null),
	},
	{
		label: 'Total monthly payment',
		getValue: (data) => formatAmount(data?.total_monthly_payment_amount),
	},
];
