import LoanStatus from '@/components/LoanStatus';
import TableOperation from '@/components/TableOperation';


const common = [
	{
		title: 'Property address',
		dataIndex: ['property','address','address1'],
		key: ['property','address','address1'].join('.'),
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
		defaultSortOrder: 'descend',
		fixed: 'left',
	},
	{
		title: 'City',
		dataIndex: ['property','address','city'],
		key: ['property','address','city'].join('.'),
		ellipsis: true,
	},
	{
		title: 'State',
		dataIndex: ['property','address','state'],
		key: ['property','address','state'].join('.'),
		ellipsis: true,
	},
	{
		title: 'Zip code',
		dataIndex: ['property','address','zip_code'],
		key: ['property','address','zip_code'].join('.'),
		ellipsis: true,
	},
	{
		title: 'Borrower',
		dataIndex: 'applicants',
		key: 'applicants',
		ellipsis: true,
		render:(v)=>v?.map(i=>i.full_name)?.join('、')
	},
];
export const pendingColumns = [
	...common,
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
		render: (v, { source }) => <LoanStatus status={v} is3rdParty={source==='borrower'} />,
	},
	{
		title: '',
		dataIndex: 'operation',
		width: 80,
		render: (_, record) => (
			<TableOperation
				record={record}
				getItems={(editItem, deleteItem) => {
					return [{ ...editItem, url: `/loans/${record?.id || '111'}` }, deleteItem];
				}}
				onSelect={(key) => {
					
				}}
			/>
		),
	},
];

export const activeColumns = [
	...common,
	{
		title: 'Lien Position',
		dataIndex: 'lienPosition',
		key: 'lienPosition',
		ellipsis: true,
	},
	{
		title: 'Lender',
		dataIndex: 'lender',
		key: 'lender',
		ellipsis: true,
	},
	{
		title: 'Loan type',
		dataIndex: 'loanType',
		key: 'loanType',
		ellipsis: true,
	},
	{
		title: 'Business strategy',
		dataIndex: 'businessStrategy',
		key: 'businessStrategy',
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
		title: 'Point',
		dataIndex: 'point',
		key: 'point',
		ellipsis: true,
	},
	{
		title: 'Maturity Date',
		dataIndex: 'maturityDate',
		key: 'maturityDate',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
	{
		title: 'Outstanding Loan Amount',
		dataIndex: 'outstandingLoanAmount',
		key: 'outstandingLoanAmount',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
	{
		title: 'Total loan amount',
		dataIndex: 'totalLoanAmount',
		key: 'totalLoanAmount',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
	{
		title: 'Hold back',
		dataIndex: 'holdback',
		key: 'holdback',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
	{
		title: 'Released Holdback',
		dataIndex: 'releasedHoldback',
		key: 'releasedHoldback',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
	{
		title: 'Monthly Payment',
		dataIndex: 'monthlyPayment',
		key: 'monthlyPayment',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
	{
		title: 'Prepayment Penalty',
		dataIndex: 'prepaymentPenalty',
		key: 'prepaymentPenalty',
		ellipsis: true,
	},
	{
		title: 'Servicing Company',
		dataIndex: 'servicingCompany',
		key: 'servicingCompany',
		ellipsis: true,
	},
];

export const archivedColumns = [
	...common,
	{
		title: 'Loan type',
		dataIndex: 'loanType',
		key: 'loanType',
		ellipsis: true,
	},
	{
		title: 'Principal Balance',
		dataIndex: 'principalBalance',
		key: 'principalBalance',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
	{
		title: 'Lender',
		dataIndex: 'lender',
		key: 'lender',
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
		dataIndex: 'paidOffDate',
		key: 'paidOffDate',
		ellipsis: true,
		sorter: true,
		sortDirections: ['ascend', 'descend'],
	},
];
