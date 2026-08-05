/**
 * Pricing 套餐静态 mock 与计费周期常量。
 * 字段与 admin `/pricing_plan` API 一致，便于后续替换为接口数据。
 */

/** @type {{ STANDARD: 'Standard', CUSTOM: 'Custom' }} */
export const PRICING_PLAN_TYPE = {
	STANDARD: 'Standard',
	CUSTOM: 'Custom',
};

/** 计费周期；与 policies[].type 一致 */
export const BILLING_CYCLE = {
	MONTHLY: 'Monthly',
	YEARLY: 'Yearly',
};

/** Segmented 选项；Yearly 暂不可用 */
export const BILLING_CYCLE_OPTIONS = [
	{ label: 'Monthly', value: BILLING_CYCLE.MONTHLY },
	{ label: 'Yearly', value: BILLING_CYCLE.YEARLY, disabled: true },
];

/** 未登录跳转注册页默认路径 */
export const DEFAULT_REGISTER_PATH = '/login/register';

/** 已登录 checkout 占位接口 */
export const DEFAULT_CHECKOUT_API = '/organization_subscription/checkout';

/** 定制版默认联系销售 mailto */
export const DEFAULT_CONTACT_SALES_HREF = 'mailto:sales@zeebee.com';

/**
 * 静态 mock 套餐列表（Base / Pro / Custom）。
 * @type {Array<object>}
 */
export const MOCK_PRICING_PLANS = [
	{
		id: 1,
		name: 'Base',
		subtitle: 'For small teams getting started',
		description:
			'Core loan management tools, document uploads, and basic reporting for growing teams.',
		plan_type: PRICING_PLAN_TYPE.STANDARD,
		trial_days: 14,
		sort_order: 1,
		is_default: false,
		policies: [
			{
				id: 101,
				type: BILLING_CYCLE.MONTHLY,
				actual_amount: 49,
				original_amount: 69,
			},
		],
	},
	{
		id: 2,
		name: 'Pro',
		subtitle: 'For teams that need full visibility',
		description:
			'Everything in Base, plus advanced workflows, integrations, and priority support.',
		plan_type: PRICING_PLAN_TYPE.STANDARD,
		trial_days: 0,
		sort_order: 2,
		is_default: true,
		policies: [
			{
				id: 201,
				type: BILLING_CYCLE.MONTHLY,
				actual_amount: 99,
				original_amount: 129,
			},
		],
	},
	{
		id: 3,
		name: 'Custom',
		subtitle: 'Tailored for your organization',
		description:
			'Dedicated onboarding, custom permissions, SLA, and enterprise-grade security controls.',
		plan_type: PRICING_PLAN_TYPE.CUSTOM,
		trial_days: 0,
		sort_order: 3,
		is_default: false,
		policies: [],
	},
];
