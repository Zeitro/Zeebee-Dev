import { formatAmount } from '@/lib/common';
import LoanPendingStatus from './LoanPendingStatus';
import TableOperation from '@/components/TableOperation';
const common = (isNotAddress) => {
	const addressArr = [
		{
			title: 'Property address',
			dataIndex: ['property', 'address', 'address1'],
			key: 'address',
			ellipsis: true,
			sorter: true,
			sortDirections: ['ascend', 'descend'],
			defaultSortOrder: 'descend',
			fixed: 'left',
			minWidth: 200,
			render: (v, { source }) => (
				<div className="flex items-center">
					<p className="overflow-hidden text-ellipsis whitespace-nowrap">{v}</p>
					{source === 'borrower' && (
						<div className="ml-2 px-2 bg-blue-light-1 text-blue-medium-default rounded-xs">
							3rd party
						</div>
					)}
				</div>
			),
		},
		{
			title: 'City',
			dataIndex: ['property', 'address', 'city'],
			key: 'city',
			width: 120,
			ellipsis: true,
		},
		{
			title: 'State',
			dataIndex: ['property', 'address', 'state'],
			key: 'state',
			width: 70,
			ellipsis: true,
		},
		{
			title: 'Zip code',
			dataIndex: ['property', 'address', 'zip_code'],
			key: 'zip_code',
			ellipsis: true,
			width: 120,
		},
	];
	return [
		...(isNotAddress ? [] : addressArr),
		{
			title: 'Borrower',
			dataIndex: 'applicants',
			key: 'applicants',
			ellipsis: true,
			render: (v) => v?.map((i) => i.name)?.join(','),
			width: 160,
		},
	];
};
export const pendingColumns = (isNotAddress, successCallback, showDeleteModal) => {
	return [
		...common(isNotAddress),
		{
			title: 'Loan purpose',
			dataIndex: 'loan_purpose',
			key: 'loan_purpose',
			ellipsis: true,
			width: 120,
		},
		{
			title: 'Target closing date',
			dataIndex: 'closing_date',
			key: 'closing_date',
			ellipsis: true,
			sorter: true,
			width: 170,
		},
		{
			title: 'Business strategy',
			dataIndex: 'business_strategy',
			key: 'business_strategy',
			ellipsis: true,
			render: (v) => v?.join(','),
			width: 160,
		},
		{
			title: 'Target lender',
			dataIndex: 'lender_company',
			key: 'lender_company',
			ellipsis: true,
			width: 180,
		},
		{
			title: 'Total loan amount',
			dataIndex: 'total_amount',
			key: 'total_amount',
			ellipsis: true,
			sorter: true,
			width: 170,
			render: (v) => formatAmount(v),
		},
		{
			title: 'Initial loan',
			dataIndex: 'initial_amount',
			key: 'initial_amount',
			ellipsis: true,
			sorter: true,
			width: 120,
			render: (v) => formatAmount(v),
		},
		{
			title: 'Hold back',
			dataIndex: 'holdback_amount',
			key: 'holdback_amount',
			ellipsis: true,
			sorter: true,
			width: 120,
			render: (v) => formatAmount(v),
		},
		{
			title: 'Estimated Closing Fund',
			dataIndex: 'closing_fund',
			key: 'closing_fund',
			ellipsis: true,
			sorter: true,
			width: 200,
			render: (v) => formatAmount(v),
		},
		{
			title: 'Rate',
			dataIndex: 'rate',
			key: 'rate',
			ellipsis: true,
			sorter: true,
			width: 90,
			render: (v) => `${v || 0}%`,
		},
		{
			title: 'Loan type',
			dataIndex: 'loan_type',
			key: 'loan_type',
			ellipsis: true,
			width: 180,
		},
		{
			title: 'Status',
			dataIndex: 'pending_status',
			key: 'pending_status',
			ellipsis: true,
			width: 150,
			align: 'center',
			fixed: 'right',
			render: (v, record) => (
				<LoanPendingStatus data={record} successCallback={successCallback} />
			),
		},
		{
			title: '',
			dataIndex: 'operation',
			width: 50,
			align: 'left',
			fixed: 'right',
			render: (_, record) => (
				<TableOperation
					record={record}
					getItems={(editItem, deleteItem) => {
						const arr = [{ ...editItem, url: `/loans/${record?.id}` }];
						if (record?.source === 'borrower') {
							arr.push(deleteItem);
						}
						return arr;
					}}
					onSelect={(key) => {
						if (key === 'delete') {
							showDeleteModal(record);
						}
					}}
				/>
			),
		},
	];
};

export const activeColumns = (isNotAddress) => [
	...common(isNotAddress),
	{
		title: 'Lien Position',
		dataIndex: 'lien_position',
		key: 'lien_position',
		ellipsis: true,
		width: 120,
	},
	{
		title: 'Lender',
		dataIndex: 'lender_company',
		key: 'lender_company',
		ellipsis: true,
		width: 180,
	},
	{
		title: 'Loan type',
		dataIndex: 'loan_type',
		key: 'loan_type',
		ellipsis: true,
		width: 180,
	},
	{
		title: 'Business strategy',
		dataIndex: 'business_strategy',
		key: 'business_strategy',
		ellipsis: true,
		render: (v) => v?.join(','),
		width: 160,
	},
	{
		title: 'Rate',
		dataIndex: 'rate',
		key: 'rate',
		ellipsis: true,
		sorter: true,
		width: 90,
		render: (v) => `${v || 0}%`,
	},
	// {
	// 	title: 'Point',
	// 	dataIndex: 'point',
	// 	key: 'point',
	// 	ellipsis: true,
	// },
	{
		title: 'Maturity Date',
		dataIndex: 'maturity_date',
		key: 'maturity_date',
		ellipsis: true,
		sorter: true,
		width: 140,
	},
	{
		title: 'Outstanding Loan Amount',
		dataIndex: 'principal_amount',
		key: 'principal_amount',
		ellipsis: true,
		sorter: true,
		width: 220,
		render: (v) => formatAmount(v),
	},
	{
		title: 'Total loan amount',
		dataIndex: 'total_amount',
		key: 'total_amount',
		ellipsis: true,
		sorter: true,
		width: 180,
		render: (v) => formatAmount(v),
	},
	{
		title: 'Hold back',
		dataIndex: 'holdback_amount',
		key: 'holdback_amount',
		ellipsis: true,
		sorter: true,
		width: 120,
		render: (v) => formatAmount(v),
	},
	// {
	// 	title: 'Released Holdback',
	// 	dataIndex: 'releasedHoldback',
	// 	key: 'releasedHoldback',
	// 	ellipsis: true,
	// 	sorter: true,
	// 	sortDirections: ['ascend', 'descend'],
	// },
	{
		title: 'Monthly Payment',
		dataIndex: 'monthly_payment',
		key: 'monthly_payment',
		ellipsis: true,
		sorter: true,
		width: 160,
		render: (v) => formatAmount(v),
	},
	{
		title: 'Prepayment Penalty',
		dataIndex: 'prepayment_penalty',
		key: 'prepayment_penalty',
		ellipsis: true,
		width: 160,
		render: (v) => `${v} months interest`,
	},
	{
		title: 'Servicing Company',
		dataIndex: 'service_company',
		key: 'service_company',
		ellipsis: true,
		width: 160,
	},
];

export const archivedColumns = (isNotAddress) => [
	...common(isNotAddress),
	{
		title: 'Loan type',
		dataIndex: 'loan_type',
		key: 'loan_type',
		ellipsis: true,
		width: 180,
	},
	{
		title: 'Principal Balance',
		dataIndex: 'principal_amount',
		key: 'principal_amount',
		ellipsis: true,
		sorter: true,
		width: 180,
		render: (v) => formatAmount(v),
	},
	{
		title: 'Lender',
		dataIndex: 'lender_company',
		key: 'lender_company',
		ellipsis: true,
		width: 180,
	},
	{
		title: 'Rate',
		dataIndex: 'rate',
		key: 'rate',
		ellipsis: true,
		sorter: true,
		width: 90,
		render: (v) => `${v || 0}%`,
	},
	{
		title: 'Paid Off Date',
		dataIndex: 'paid_off_date',
		key: 'paid_off_date',
		ellipsis: true,
		sorter: true,
		width: 160,
	},
];
