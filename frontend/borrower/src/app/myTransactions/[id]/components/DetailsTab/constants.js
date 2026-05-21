import { assemblyAddress, formatDate } from '@repo/shared-methods/utils';

export const TRANSACTION_FIELDS = [
	{ label: 'Transaction type', getValue: (data) => data?.type },
	{ label: 'Status', getValue: (data) => data?.status },
	{
		label: 'COE date',
		getValue: (data) => {
			const raw = data?.coe_date ?? data?.coeDate;
			return raw ? formatDate(raw) : null;
		},
	},
	{
		label: 'Created',
		getValue: (data) => {
			const raw = data?.created_at ?? data?.createdAt;
			return raw ? formatDate(raw) : null;
		},
	},
	{
		label: 'Updated',
		getValue: (data) => {
			const raw = data?.updated_at ?? data?.updatedAt;
			return raw ? formatDate(raw) : null;
		},
	},
];

export const PROPERTY_FIELDS = [
	{
		label: 'Property',
		getValue: (data) => assemblyAddress(data?.property?.address),
	},
	{ label: 'Street 1', getValue: (data) => data?.property?.address?.address1 },
	{ label: 'Street 2', getValue: (data) => data?.property?.address?.address2 },
	{ label: 'City', getValue: (data) => data?.property?.address?.city },
	{ label: 'State', getValue: (data) => data?.property?.address?.state },
	{ label: 'Zip code', getValue: (data) => data?.property?.address?.zip_code },
];
