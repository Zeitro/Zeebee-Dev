import { formatAmount, formatDate, formatNumber, formatPhone } from '@repo/shared-methods/utils';

/**
 * Details Tab 各模块只读字段配置。
 * Property Info 数值字段取自 rowData.property；Title/Escrow 取自 rowData.escrows[]。
 */

/** @param {boolean|null|undefined} value */
function formatYesNo(value) {
	if (value == null) return null;
	return value ? 'Yes' : 'No';
}

/** @param {object|null} data transaction rowData */
function getProperty(data) {
	return data?.property;
}

/** Form.List 空 escrow 占位（无数据时仍展示一行表单） */
export const EMPTY_ESCROW = {
	escrow_company: '',
	officer_name: '',
	officer_phone: '',
	officer_email: '',
	no: '',
	checklist: {
		has_wire_instruction: false,
		has_emd_receipt: false,
		has_prelim_report: false,
		has_nhd: false,
		has_colored_map: false,
	},
};

export const PROPERTY_INFO_VIEW_FIELDS = [
	{
		label: 'Property type',
		key: 'property_type',
		getValue: (data) => getProperty(data)?.property_type,
	},
	{
		label: 'Strategy/Business line',
		key: 'business_strategy',
		getValue: (data) => data?.business_strategy?.join(', '),
	},
	{
		label: 'Title holder',
		key: 'title_holder',
		getValue: (data) => data?.title_holder,
	},
	{
		label: 'Living area',
		key: 'living_area_sqft',
		getValue: (data) => {
			const raw = getProperty(data)?.living_area_sqft;
			return raw != null && raw !== '' ? formatNumber(raw, 'sqft') : null;
		},
	},
	{
		label: 'Units',
		key: 'unit',
		getValue: (data) => getProperty(data)?.unit,
	},
	{
		label: 'Lot area',
		key: 'lot_area_sqft',
		getValue: (data) => {
			const raw = getProperty(data)?.lot_area_sqft;
			return raw != null && raw !== '' ? formatNumber(raw, 'sqft') : null;
		},
	},
	{
		label: 'Beds',
		key: 'bedrooms',
		getValue: (data) => getProperty(data)?.bedrooms,
	},
	{
		label: 'Baths',
		key: 'bathrooms',
		getValue: (data) => getProperty(data)?.bathrooms,
	},
];

/** Transaction Info 按行分组，只读态与编辑态栅格对齐 */
export const TRANSACTION_INFO_VIEW_ROWS = [
	{
		gridClassName: 'grid w-full grid-cols-1 gap-6 sm:grid-cols-2',
		fields: [
			{
				label: 'Purchase price',
				key: 'purchase_amount',
				getValue: (data) =>
					data?.purchase_amount != null && data?.purchase_amount !== ''
						? formatAmount(data.purchase_amount)
						: null,
			},
			{
				label: 'EMD',
				key: 'emd_amount',
				getValue: (data) =>
					data?.emd_amount != null && data?.emd_amount !== ''
						? formatAmount(data.emd_amount)
						: null,
			},
		],
	},
	{
		gridClassName: 'grid w-full grid-cols-1 gap-6 sm:grid-cols-2 xl:grid-cols-4',
		fields: [
			{
				label: 'Contract date',
				key: 'contract_date',
				getValue: (data) => (data?.contract_date ? formatDate(data.contract_date) : null),
			},
			{
				label: 'EMD due date',
				key: 'emd_due_date',
				getValue: (data) => (data?.emd_due_date ? formatDate(data.emd_due_date) : null),
			},
			{
				label: 'Contingency removal date',
				key: 'contingency_removal_date',
				getValue: (data) =>
					data?.contingency_removal_date
						? formatDate(data.contingency_removal_date)
						: null,
			},
			{
				label: 'COE date',
				key: 'coe_date',
				getValue: (data) => (data?.coe_date ? formatDate(data.coe_date) : null),
			},
		],
	},
	{
		gridClassName: 'grid w-full grid-cols-1 gap-6 sm:grid-cols-2',
		fields: [
			{ label: 'PM', key: 'pm', getValue: (data) => data?.pm },
			{ label: 'TC', key: 'tc', getValue: (data) => data?.tc },
		],
	},
	{
		gridClassName: 'grid w-full grid-cols-1 gap-6 sm:grid-cols-2',
		fields: [
			{
				label: 'Who pays what',
				key: 'cost_allocation',
				getValue: (data) => data?.cost_allocation,
			},
			{ label: 'TC notes', key: 'note', getValue: (data) => data?.note },
		],
	},
];

/**
 * Title/Escrow 字段配置；getValue 入参为单条 escrow（rowData.escrows[i]）。
 * 与 Form.List 字段名对齐：escrow_company / officer_* / no / checklist.has_*。
 */
export const TITLE_ESCROW_VIEW_ROWS = [
	{
		gridClassName: 'grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5',
		fields: [
			{
				label: 'Title company',
				key: 'escrow_company',
				getValue: (escrow) => escrow?.escrow_company,
			},
			{
				label: 'Escrow officer',
				key: 'officer_name',
				getValue: (escrow) => escrow?.officer_name,
			},
			{
				label: 'Escrow officer phone',
				key: 'officer_phone',
				getValue: (escrow) => formatPhone(escrow?.officer_phone),
			},
			{
				label: 'Escrow officer email',
				key: 'officer_email',
				getValue: (escrow) => escrow?.officer_email,
			},
			{
				label: 'Escrow #',
				key: 'no',
				getValue: (escrow) => escrow?.no,
			},
		],
	},
	{
		gridClassName: 'grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5',
		fields: [
			{
				label: 'Wire instruction',
				key: 'has_wire_instruction',
				getValue: (escrow) => formatYesNo(escrow?.checklist?.has_wire_instruction),
			},
			{
				label: 'EMD receipt',
				key: 'has_emd_receipt',
				getValue: (escrow) => formatYesNo(escrow?.checklist?.has_emd_receipt),
			},
			{
				label: 'Prelim report',
				key: 'has_prelim_report',
				getValue: (escrow) => formatYesNo(escrow?.checklist?.has_prelim_report),
			},
			{
				label: 'NHD',
				key: 'has_nhd',
				getValue: (escrow) => formatYesNo(escrow?.checklist?.has_nhd),
			},
			{
				label: 'Colored map',
				key: 'has_colored_map',
				getValue: (escrow) => formatYesNo(escrow?.checklist?.has_colored_map),
			},
		],
	},
];
