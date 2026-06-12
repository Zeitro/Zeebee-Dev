/**
 * Subscription 列表 mock 数据与展示常量。
 * 接入 `/subscription/page` 等接口后可移除 mock 与 {@link getMockSubscriptionList}。
 */

/** @type {{ TRIALING: 'trialing', ACTIVE: 'active', EXPIRING: 'expiring', EXPIRED: 'expired', SUSPENDED: 'suspended' }} */
export const SUBSCRIPTION_STATUS = {
	TRIALING: 'trialing',
	ACTIVE: 'active',
	EXPIRING: 'expiring',
	EXPIRED: 'expired',
	SUSPENDED: 'suspended',
};

/** 筛选 / 表单下拉；value 与 API 字段 status 一致 */
export const subscriptionStatusOptions = [
	{ label: 'Trialing', value: SUBSCRIPTION_STATUS.TRIALING },
	{ label: 'Active', value: SUBSCRIPTION_STATUS.ACTIVE },
	{ label: 'Expiring', value: SUBSCRIPTION_STATUS.EXPIRING },
	{ label: 'Expired', value: SUBSCRIPTION_STATUS.EXPIRED },
	{ label: 'Suspended', value: SUBSCRIPTION_STATUS.SUSPENDED },
];

export const SUBSCRIPTION_TYPE = {
	MONTHLY: 'monthly',
	YEARLY: 'yearly',
};

export const subscriptionTypeOptions = [
	{ label: 'Monthly', value: SUBSCRIPTION_TYPE.MONTHLY },
	{ label: 'Yearly', value: SUBSCRIPTION_TYPE.YEARLY },
];

export const planNameOptions = [
	{ label: 'Basic', value: 'Basic' },
	{ label: 'Advanced', value: 'Advanced' },
	{ label: 'Custom', value: 'Custom' },
];

/** 状态徽章颜色；API 字段 status */
export const subscriptionStatusColor = {
	trialing: {
		border: 'var(--blue-light)',
		bg: 'var(--blue-light-1)',
		color: 'var(--blue-dark)',
	},
	active: {
		border: 'var(--green-light-100)',
		bg: 'var(--green-light)',
		color: 'var(--green-medium)',
	},
	expiring: {
		border: 'var(--orange-medium)',
		bg: 'var(--orange-light)',
		color: 'var(--orange-dark)',
	},
	expired: {
		border: 'var(--base-border)',
		bg: 'var(--gray-light)',
		color: 'var(--gray-medium)',
	},
	suspended: {
		border: 'var(--orange-medium)',
		bg: 'var(--orange-light)',
		color: 'var(--orange-dark)',
	},
};

/** 与 UI 稿一致的示例数据 */
export const MOCK_SUBSCRIPTIONS = [
	{
		id: '1',
		subscription_no: 'SU202605001',
		customer_name: 'ABC Lending',
		plan_name: 'Advanced',
		subscription_type: SUBSCRIPTION_TYPE.MONTHLY,
		start_date: '2026-05-01',
		expiry_date: '2026-06-01',
		status: SUBSCRIPTION_STATUS.ACTIVE,
	},
	{
		id: '2',
		subscription_no: 'SU202605002',
		customer_name: 'Nova Borrower Hub',
		plan_name: 'Basic',
		subscription_type: SUBSCRIPTION_TYPE.YEARLY,
		start_date: '2026-04-15',
		expiry_date: '2027-04-15',
		status: SUBSCRIPTION_STATUS.TRIALING,
	},
	{
		id: '3',
		subscription_no: 'SU202605003',
		customer_name: 'Alphax Investment',
		plan_name: 'Custom',
		subscription_type: SUBSCRIPTION_TYPE.YEARLY,
		start_date: '2025-12-01',
		expiry_date: '2026-06-08',
		status: SUBSCRIPTION_STATUS.EXPIRING,
	},
	{
		id: '4',
		subscription_no: 'SU202604010',
		customer_name: 'Sunrise Capital',
		plan_name: 'Basic',
		subscription_type: SUBSCRIPTION_TYPE.MONTHLY,
		start_date: '2026-03-01',
		expiry_date: '2026-04-01',
		status: SUBSCRIPTION_STATUS.EXPIRED,
	},
	{
		id: '5',
		subscription_no: 'SU202603015',
		customer_name: 'Metro Finance',
		plan_name: 'Advanced',
		subscription_type: SUBSCRIPTION_TYPE.MONTHLY,
		start_date: '2026-01-15',
		expiry_date: '2026-05-15',
		status: SUBSCRIPTION_STATUS.SUSPENDED,
	},
];

/**
 * 表格展示 subscription_type 文案
 * @param {string} value monthly | yearly
 */
export function formatSubscriptionType(value) {
	const map = {
		[SUBSCRIPTION_TYPE.MONTHLY]: 'Monthly',
		[SUBSCRIPTION_TYPE.YEARLY]: 'Yearly',
	};
	return map[value] || value;
}

/**
 * Mock 分页列表；keyword 匹配 subscription_no / customer_name。
 * @param {object} params FilterTable 请求参数
 * @returns {Promise<{ list: object[], total: number }>}
 */
export function getMockSubscriptionList(params = {}) {
	const { filter = {} } = params;
	const current = params.page || params.current || 1;
	const pageSize = params.limit || params.pageSize || 10;
	const keyword = filter.keyword?.trim().toLowerCase();

	let list = [...MOCK_SUBSCRIPTIONS];

	if (keyword) {
		list = list.filter(
			(item) =>
				item.subscription_no?.toLowerCase().includes(keyword) ||
				item.customer_name?.toLowerCase().includes(keyword),
		);
	}

	if (filter.plan_name) {
		list = list.filter((item) => item.plan_name === filter.plan_name);
	}

	if (filter.subscription_type) {
		list = list.filter((item) => item.subscription_type === filter.subscription_type);
	}

	if (filter.status) {
		list = list.filter((item) => item.status === filter.status);
	}

	if (filter.start_date?.length === 2) {
		const [start, end] = filter.start_date;
		list = list.filter((item) => item.start_date >= start && item.start_date <= end);
	}

	const startIndex = (current - 1) * pageSize;
	return Promise.resolve({
		list: list.slice(startIndex, startIndex + pageSize),
		total: list.length,
	});
}

/**
 * Mock 新增/更新；接入 API 后可删除。
 * @param {object} payload 表单字段
 * @param {string} [id] 编辑时传入
 */
export function upsertMockSubscription(payload, id) {
	if (id) {
		const index = MOCK_SUBSCRIPTIONS.findIndex((item) => item.id === id);
		if (index >= 0) {
			MOCK_SUBSCRIPTIONS[index] = { ...MOCK_SUBSCRIPTIONS[index], ...payload, id };
		}
		return;
	}
	MOCK_SUBSCRIPTIONS.push({
		id: String(Date.now()),
		...payload,
	});
}

/**
 * Mock 暂停 / 恢复；接入 API 后可删除。
 * @param {string} id
 * @param {'suspended' | 'active'} nextStatus
 */
export function toggleMockSubscriptionStatus(id, nextStatus) {
	const index = MOCK_SUBSCRIPTIONS.findIndex((item) => item.id === id);
	if (index >= 0) {
		MOCK_SUBSCRIPTIONS[index].status = nextStatus;
	}
}
