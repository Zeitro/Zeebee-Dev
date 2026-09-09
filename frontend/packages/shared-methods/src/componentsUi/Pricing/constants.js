/**
 * Pricing 计费周期与默认路径常量。
 * 套餐列表由 GET `/pricing_plans/public` 返回，字段与接口 `plans[]` 一致。
 */

/** @type {{ STANDARD: 'Standard', CUSTOM: 'Custom' }} */
export const PRICING_PLAN_TYPE = {
	STANDARD: 'Standard',
	CUSTOM: 'Custom',
};

/** 计费周期；与 policies[].type 一致 */
export const BILLING_CYCLE = {
	MONTHLY: 'Monthly',
	ANNUAL: 'Annual',
};

/** Segmented 选项；Yearly 暂不可用 */
export const BILLING_CYCLE_OPTIONS = [
	{ label: 'Monthly', value: BILLING_CYCLE.MONTHLY },
	{ label: 'Annual', value: BILLING_CYCLE.ANNUAL, disabled: true },
];

/** 未登录跳转注册页默认路径 */
export const DEFAULT_REGISTER_PATH = '/login/register';

/** 定制版销售联系邮箱 */
export const DEFAULT_SALES_EMAIL = 'sales@zeebee.com';

/** 定制版销售联系电话 */
export const DEFAULT_SALES_PHONE = '(800)933-2333';

/** Contact Sales Popover 提示语 */
export const CONTACT_SALES_TIP =
	'Looking for a custom plan? Reach out and our team will tailor pricing to your organization.';

/** 定制版默认联系销售 mailto（兼容旧用法） */
export const DEFAULT_CONTACT_SALES_HREF = `mailto:${DEFAULT_SALES_EMAIL}`;

/** 组织订阅状态；与 organization/subscription/current → subscription.status 一致 */
export const SUBSCRIPTION_STATUS = {
	TRIALING: 'Trialing',
	EXPIRED: 'Expired',
	ACTIVE: 'Active',
	EXPIRING: 'Expiring',
};
