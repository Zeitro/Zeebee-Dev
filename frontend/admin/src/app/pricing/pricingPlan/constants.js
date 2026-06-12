/**
 * Pricing plan 列表 mock 数据与展示常量。
 * 接入 `/pricing_plan/page` 等接口后可移除 mock 与 {@link getMockPricingPlanList}。
 */

/** @type {{ ACTIVE: 'active', INACTIVE: 'inactive' }} */
export const PRICING_PLAN_STATUS = {
	ACTIVE: 'active',
	INACTIVE: 'inactive',
};

export const planStatusOptions = [
	{ label: 'Active', value: PRICING_PLAN_STATUS.ACTIVE },
	{ label: 'Inactive', value: PRICING_PLAN_STATUS.INACTIVE },
];

/** 计划状态徽章颜色；API 字段 status：active / inactive */
export const planStatusColor = {
	active: {
		border: 'var(--green-light-100)',
		bg: 'var(--green-light)',
		color: 'var(--green-medium)',
	},
	inactive: {
		border: 'var(--orange-medium)',
		bg: 'var(--orange-light)',
		color: 'var(--orange-dark)',
	},
};

/** 与 UI 稿一致的示例数据 */
export const MOCK_PRICING_PLANS = [
	{
		id: '1',
		name: 'Basic',
		subtitle: 'Best for fast launch teams',
		description: 'Fixed pricing with monthly and yearly limits',
		monthly_price: 49,
		yearly_price: 499,
		trial_days: 7,
		status: PRICING_PLAN_STATUS.ACTIVE,
	},
	{
		id: '2',
		name: 'Advanced',
		subtitle: 'For scaling organizations',
		description: 'Volume discounts and advanced billing',
		monthly_price: 129,
		yearly_price: 1299,
		trial_days: 14,
		status: PRICING_PLAN_STATUS.ACTIVE,
	},
	{
		id: '3',
		name: 'Custom',
		subtitle: 'Tailored enterprise solution',
		description: 'A flexible plan',
		monthly_price: null,
		yearly_price: null,
		trial_days: 30,
		status: PRICING_PLAN_STATUS.INACTIVE,
	},
];

/**
 * Mock 分页列表；keyword 匹配 name / subtitle / description，status 精确匹配。
 * @param {object} params FilterTable 请求参数
 * @returns {Promise<{ list: object[], total: number }>}
 */
export function getMockPricingPlanList(params = {}) {
	const { filter = {}, current = 1, pageSize = 10 } = params;
	const keyword = filter.keyword?.trim().toLowerCase();

	let list = [...MOCK_PRICING_PLANS];

	if (keyword) {
		list = list.filter(
			(item) =>
				item.name?.toLowerCase().includes(keyword) ||
				item.subtitle?.toLowerCase().includes(keyword) ||
				item.description?.toLowerCase().includes(keyword),
		);
	}

	if (filter.status) {
		list = list.filter((item) => item.status === filter.status);
	}

	const start = (current - 1) * pageSize;
	return Promise.resolve({
		list: list.slice(start, start + pageSize),
		total: list.length,
	});
}

/**
 * Mock 新增/更新；接入 API 后可删除。
 * @param {object} payload 表单字段
 * @param {string} [id] 编辑时传入
 */
export function upsertMockPricingPlan(payload, id) {
	if (id) {
		const index = MOCK_PRICING_PLANS.findIndex((item) => item.id === id);
		if (index >= 0) {
			MOCK_PRICING_PLANS[index] = { ...MOCK_PRICING_PLANS[index], ...payload, id };
		}
		return;
	}
	MOCK_PRICING_PLANS.push({
		id: String(Date.now()),
		...payload,
	});
}
