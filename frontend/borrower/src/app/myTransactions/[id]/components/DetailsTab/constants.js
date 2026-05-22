import { formatAmount, formatDate } from '@repo/shared-methods/utils';

export const BASIC_INFO_FORM_FIELDS = [
	'coe_date',
	'purchase_price',
	'emd',
	'property_type',
	'pm',
	'business_strategy',
	'tc_notes',
];

export const BASIC_INFO_VIEW_FIELDS = [
	{
		label: 'COE',
		key: 'coe_date',
		getValue: (data) => {
			const raw = data?.coe_date ?? data?.coeDate;
			return raw ? formatDate(raw) : null;
		},
	},
	{
		label: 'Purchase price',
		key: 'purchase_price',
		getValue: (data) => {
			const raw = data?.purchase_price ?? data?.purchasePrice;
			return raw != null && raw !== '' ? formatAmount(raw) : null;
		},
	},
	{
		label: 'EMD',
		key: 'emd',
		getValue: (data) => (data?.emd != null && data?.emd !== '' ? formatAmount(data.emd) : null),
	},
	{
		label: 'Property type',
		key: 'property_type',
		getValue: (data) =>
			data?.property_type ?? data?.propertyType ?? data?.property?.property_type,
	},
	{
		label: 'PM',
		key: 'pm',
		getValue: (data) => data?.pm,
	},
	{
		label: 'Strategy',
		key: 'business_strategy',
		getValue: (data) => data?.business_strategy ?? data?.businessStrategy,
	},
	{
		label: 'TC Notes',
		key: 'tc_notes',
		getValue: (data) => data?.tc_notes ?? data?.tcNotes,
		fullWidth: true,
	},
];
