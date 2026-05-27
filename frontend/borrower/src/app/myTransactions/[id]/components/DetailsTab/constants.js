import { formatAmount, formatDate } from '@repo/shared-methods/utils';

export const BASIC_INFO_VIEW_FIELDS = [
	{
		label: 'COE',
		key: 'coe_date',
		getValue: (data) => {
			const raw = data?.coe_date;
			return raw ? formatDate(raw) : null;
		},
	},
	{
		label: 'Purchase price',
		key: 'purchase_amount',
		getValue: (data) => {
			const raw = data?.purchase_amount;
			return raw != null && raw !== '' ? formatAmount(raw) : null;
		},
	},
	{
		label: 'EMD',
		key: 'emd_amount',
		getValue: (data) =>
			data?.emd_amount != null && data?.emd_amount !== ''
				? formatAmount(data.emd_amount)
				: null,
	},
	{
		label: 'Property type',
		key: 'property_type',
		getValue: (data) => data?.property_type,
	},
	{
		label: 'PM',
		key: 'pm',
		getValue: (data) => data?.pm,
	},
	{
		label: 'Business strategy',
		key: 'business_strategy',
		getValue: (data) => data?.business_strategy,
	},
	{
		label: 'TC Notes',
		key: 'note',
		getValue: (data) => data?.note,
		fullWidth: true,
	},
];
