/** 右侧内容区展示类型；非 upload 时可在此扩展并在 CategoryContentPanel 中分支 */
export const CATEGORY_CONTENT_TYPE = {
	UPLOAD: 'upload',
};

export function resolveCategoryContentType(category) {
	return category?.contentType ?? CATEGORY_CONTENT_TYPE.UPLOAD;
}

/** 转为 CollapsibleDocumentList 所需数据结构 */
export function toDocumentListItems(documents = []) {
	return documents.map((doc, index) => ({
		id: doc.id ?? doc.key ?? index,
		document_type: doc.document_type ?? doc.name ?? doc.file_type,
		updated_at: doc.updated_at ?? doc.last_uploaded_at ?? null,
		status: doc.status ?? 'Pending upload',
		active: doc.active ?? true,
		files: doc.files ?? [],
	}));
}

/** Mock document categories — replace when `/transaction/:id/documents` is wired */
export const MOCK_DOCUMENT_CATEGORIES = [
	{
		key: 'contract',
		name: 'Contract',
		required: true,
		completed: true,
		documents: [
			{
				key: 'purchase_agreement',
				name: 'Purchase Agreement',
				required: true,
				status: 'Ready for review',
				updated_at: '2026-04-10',
				active: true,
				files: [
					{
						id: 1001,
						file_url: 'https://example.com/uploads/purchase_agreement.pdf',
						updated_at: '2026-04-10',
						status: 'Ready for review',
					},
				],
			},
			{
				key: 'addendum',
				name: 'Addendum',
				required: false,
				status: 'Pending upload',
				updated_at: null,
				active: false,
				files: [],
			},
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
			{
				key: 'eoi',
				name: 'EOI',
				required: true,
				status: 'Pending upload',
				updated_at: null,
				active: true,
				files: [],
			},
			{
				key: 'payment_receipt',
				name: 'Payment Receipt',
				required: true,
				status: 'Ready for review',
				updated_at: '2026-04-15',
				active: true,
				files: [
					{
						id: 2001,
						file_url: 'https://example.com/uploads/payment_receipt.pdf',
						updated_at: '2026-04-15',
						status: 'Ready for review',
					},
				],
			},
			{
				key: 'full_policy',
				name: 'Full Policy',
				required: true,
				status: 'Requested',
				updated_at: null,
				active: false,
				files: [],
			},
			{
				key: 'flood_insurance',
				name: 'Flood Insurance',
				required: false,
				status: 'Pending upload',
				updated_at: null,
				active: true,
				files: [],
			},
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
