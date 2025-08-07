import LoanPendingStatus from './LoanPendingStatus';
import TableOperation from '@/components/TableOperation';


const common = (isNotAddress) => {
	const addressArr = [
		{
		title: 'Property address',
		dataIndex: ['property','address','address1'],
		key: 'address',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
		defaultSortOrder: 'descend',
		fixed: 'left',
	},
	{
		title: 'City',
		dataIndex: ['property','address','city'],
		key: 'city',
		ellipsis: true,
	},
	{
		title: 'State',
		dataIndex: ['property','address','state'],
		key: 'state',
		ellipsis: true,
	},
	{
		title: 'Zip code',
		dataIndex: ['property','address','zip_code'],
		key: 'zip_code',
		ellipsis: true,
	},
	]
	return [
	...(isNotAddress?[]:addressArr),
	{
		title: 'Borrower',
		dataIndex: 'applicants',
		key: 'applicants',
		ellipsis: true,
		render:(v)=>v?.map(i=>i.full_name)?.join('、')
	},
]
};
export const pendingColumns = (isNotAddress,successCallback) => {
	return [
	...common(isNotAddress),
	{
		title: 'Loan purpose',
		dataIndex: 'loan_purpose',
		key: 'loan_purpose',
		ellipsis: true,
	},
	{
		title: 'Target closing date',
		dataIndex: 'closing_date',
		key: 'closing_date',
		ellipsis: true,
		sorter: true,
	},
	{
		title: 'Business strategy',
		dataIndex: 'business_strategy',
		key: 'business_strategy',
		ellipsis: true,
		render:(v)=>v?.join('、')
	},
	{
		title: 'Target lender',
		dataIndex: 'lender_company',
		key: 'lender_company',
		ellipsis: true,
	},
	{
		title: 'Total loan amount',
		dataIndex: 'total_amount',
		key: 'total_amount',
		ellipsis: true,
		sorter: true,
	},
	{
		title: 'Initial loan',
		dataIndex: 'initial_amount',
		key: 'initial_amount',
		ellipsis: true,
		sorter: true,
	},
	{
		title: 'Hold back',
		dataIndex: 'holdback_amount',
		key: 'holdback_amount',
		ellipsis: true,
		sorter: true,
	},
	{
		title: 'Estimated Closing Fund',
		dataIndex: 'closing_fund',
		key: 'closing_fund',
		ellipsis: true,
		sorter: true,
	},
	{
		title: 'Rate',
		dataIndex: 'rate',
		key: 'rate',
		ellipsis: true,
		sorter: true,
	},
	{ title: 'Loan type', dataIndex: 'loan_type', key: 'loan_type', ellipsis: true },
	{
		title: 'Status',
		dataIndex: 'pending_status',
		key: 'pending_status',
		ellipsis: true,
		width: 160,
		render: (v, record) => <LoanPendingStatus data={record} successCallback={successCallback}/>,
	},
	{
		title: '',
		dataIndex: 'operation',
		width: 80,
		render: (_, record) => (
			<TableOperation
				record={record}
				getItems={(editItem, deleteItem) => {
					return [{ ...editItem, url: `/loans/${record?.id}` }, deleteItem];
				}}
				onSelect={(key) => {
					
				}}
			/>
		),
	},
]
};

export const activeColumns = (isNotAddress)=>[
	...common(isNotAddress),
	{
		title: 'Lien Position',
		dataIndex: 'lien_position',
		key: 'lien_position',
		ellipsis: true,
	},
	{
		title: 'Lender',
		dataIndex: 'lender_company',
		key: 'lender_company',
		ellipsis: true,
	},
	{
		title: 'Loan type',
		dataIndex: 'loan_type',
		key: 'loan_type',
		ellipsis: true,
	},
	{
		title: 'Business strategy',
		dataIndex: 'business_strategy',
		key: 'business_strategy',
		ellipsis: true,
		render:(v)=>v?.join('、')
	},
	{
		title: 'Rate',
		dataIndex: 'rate',
		key: 'rate',
		ellipsis: true,
		sorter: true,
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
	},
	{
		title: 'Outstanding Loan Amount',
		dataIndex: 'principal_amount',
		key: 'principal_amount',
		ellipsis: true,
		sorter: true,
	},
	{
		title: 'Total loan amount',
		dataIndex: 'total_amount',
		key: 'total_amount',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
	{
		title: 'Hold back',
		dataIndex: 'holdback_amount',
		key: 'holdback_amount',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
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
	},
	{
		title: 'Prepayment Penalty',
		dataIndex: 'prepayment_penalty',
		key: 'prepayment_penalty',
		ellipsis: true,
	},
	{
		title: 'Servicing Company',
		dataIndex: 'service_company',
		key: 'service_company',
		ellipsis: true,
	},
];

export const archivedColumns = (isNotAddress)=>[
	...common(isNotAddress),
	{
		title: 'Loan type',
		dataIndex: 'loan_type',
		key: 'loan_type',
		ellipsis: true,
	},
	{
		title: 'Principal Balance',
		dataIndex: 'principal_amount',
		key: 'principal_amount',
		ellipsis: true,
		sorter: true,
	},
	{
		title: 'Lender',
		dataIndex: 'lender_company',
		key: 'lender_company',
		ellipsis: true,
	},
	{
		title: 'Rate',
		dataIndex: 'rate',
		key: 'rate',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
	{
		title: 'Paid Off Date',
		dataIndex: 'paid_off_date',
		key: 'paid_off_date',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
];
