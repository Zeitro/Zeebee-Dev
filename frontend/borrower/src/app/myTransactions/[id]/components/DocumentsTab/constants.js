/** Mock document categories — replace when `/transaction/:id/documents` is wired */
export const MOCK_DOCUMENT_CATEGORIES = [
	{
		key: 'contract',
		name: 'Contract',
		required: true,
		completed: true,
		documents: [
			{ key: 'purchase_agreement', name: 'Purchase Agreement', required: true },
			{ key: 'addendum', name: 'Addendum', required: false },
		],
	},
	{
		key: 'escrow',
		name: 'Escrow',
		required: true,
		completed: true,
		documents: [{ key: 'escrow_instructions', name: 'Escrow Instructions', required: true }],
	},
	{
		key: 'disclosures',
		name: 'Disclosures',
		required: true,
		completed: true,
		documents: [{ key: 'disclosure_package', name: 'Disclosure Package', required: true }],
	},
	{
		key: 'loan',
		name: 'Loan',
		required: true,
		completed: true,
		documents: [{ key: 'loan_estimate', name: 'Loan Estimate', required: true }],
	},
	{
		key: 'insurance',
		name: 'Insurance',
		required: true,
		pendingCount: 3,
		documents: [
			{ key: 'eoi', name: 'EOI', required: true },
			{ key: 'payment_receipt', name: 'Payment Receipt', required: true },
			{ key: 'full_policy', name: 'Full Policy', required: true },
			{ key: 'flood_insurance', name: 'Flood Insurance', required: false },
		],
	},
	{
		key: 'closing_package',
		name: 'Closing Package',
		required: true,
		pendingCount: 5,
		documents: [
			{ key: 'closing_disclosure', name: 'Closing Disclosure', required: true },
			{ key: 'settlement_statement', name: 'Settlement Statement', required: true },
		],
	},
	{
		key: '1031_exchange',
		name: '1031 Exchange',
		required: false,
		documents: [],
	},
	{
		key: 'lease',
		name: 'Lease',
		required: false,
		documents: [],
	},
	{
		key: 'hoa',
		name: 'HOA',
		required: false,
		documents: [],
	},
	{
		key: 'solar',
		name: 'Solar',
		required: false,
		documents: [],
	},
];
