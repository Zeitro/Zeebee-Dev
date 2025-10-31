import { formatAmount } from '@/lib/format';
import { LoanStage, LoanPendingStatus } from './ChangeStageStatus';
import TableOperation from '@/components/TableOperation';
export const get3rdPartyMark = (v, { source }) => (
	<div className="flex items-center">
		<p className="overflow-hidden text-ellipsis whitespace-nowrap" title={v}>
			{v}
		</p>
		{source === 'borrower' && (
			<div className="ml-2 px-2 bg-blue-light-1 text-blue-medium-default rounded-xs">
				3rd party
			</div>
		)}
	</div>
);
const common = (isNotAddress) => {
	const addressArr = [
		{
			title: 'Property address',
			dataIndex: ['loan_property', 'property_base', 'address', 'address1'],
			key: 'address',
			ellipsis: true,
			sorter: true,
			sortDirections: ['ascend', 'descend'],
			// defaultSortOrder: 'descend',
			fixed: 'left',
			minWidth: 200,
			render: get3rdPartyMark,
		},
		{
			title: 'City',
			dataIndex: ['loan_property', 'property_base', 'address', 'city'],
			key: 'city',
			width: 120,
			ellipsis: true,
		},
		{
			title: 'State',
			dataIndex: ['loan_property', 'property_base', 'address', 'state'],
			key: 'state',
			width: 70,
			ellipsis: true,
		},
		{
			title: 'Zip code',
			dataIndex: ['loan_property', 'property_base', 'address', 'zip_code'],
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
			render: (v, row) => {
				const value = v?.map((i) => i.full_name)?.join(',');
				if (isNotAddress) return get3rdPartyMark(value, row);
				return value;
			},
			width: isNotAddress ? 200 : 160,
		},
	];
};

const getOpt = (showDeleteModal) => {
	return {
		title: '',
		dataIndex: 'operation',
		width: 50,
		align: 'left',
		fixed: 'right',
		render: (_, record) => (
			<TableOperation
				record={record}
				getItems={(editItem, deleteItem, detailsItem) => {
					const arr = [];
					if (record?.source === 'borrower') {
						arr.push({ ...editItem, url: `/loans/${record?.id}` });
						arr.push(deleteItem);
					} else {
						arr.push({
							...detailsItem,
							url: `/loans/zeeBeeDetails?id=${record?.id}`,
						});
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
	};
};

const getStageOpt = (successCallback) => {
	return {
		title: 'Stage',
		dataIndex: 'stage',
		key: 'stage',
		ellipsis: true,
		width: 150,
		align: 'center',
		fixed: 'right',
		render: (v, record) => <LoanStage data={record} successCallback={successCallback} />,
	};
};

export const pendingColumns = (isNotAddress, successCallback, showDeleteModal, stage) => {
	const status = {
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
	};
	return [
		...common(isNotAddress),
		{
			title: 'Loan purpose',
			dataIndex: 'loan_purpose',
			key: 'loan_purpose',
			ellipsis: true,
			width: 120,
			render: (v) => v || '--',
		},
		{
			title: 'Target closing date',
			dataIndex: 'closing_date',
			key: 'closing_date',
			ellipsis: true,
			sorter: true,
			width: 170,
			render: (v) => v || '--',
		},
		{
			title: 'Business strategy',
			dataIndex: 'business_strategy',
			key: 'business_strategy',
			ellipsis: true,
			render: (v) => v?.join(',') || '--',
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
			render: (v) => (v ? `${v}%` : '--'),
		},
		{
			title: 'Loan type',
			dataIndex: 'loan_type',
			key: 'loan_type',
			ellipsis: true,
			width: 180,
		},
		...(stage !== 'Draft' ? [getStageOpt(successCallback)] : []),
		...(stage === 'Pending' ? [status] : []),
		getOpt(showDeleteModal, stage),
	];
};

export const activeColumns = (isNotAddress, successCallback, showDeleteModal) => [
	...common(isNotAddress),
	{
		title: 'Lien Position',
		dataIndex: 'lien_position',
		key: 'lien_position',
		ellipsis: true,
		width: 120,
		render: (v) => v || '--',
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
		render: (v) => (v ? `${v}%` : '--'),
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
		render: (v) => v || '--',
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
		render: (v) => v || '--',
	},
	getStageOpt(successCallback),
	getOpt(showDeleteModal),
];

export const archivedColumns = (isNotAddress, successCallback, showDeleteModal) => [
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
		render: (v) => (v ? `${v}%` : '--'),
	},
	{
		title: 'Paid Off Date',
		dataIndex: 'paid_off_date',
		key: 'paid_off_date',
		ellipsis: true,
		sorter: true,
		width: 160,
		render: (v) => v || '--',
	},
	getStageOpt(successCallback),
	getOpt(showDeleteModal),
];
