/**
 * Pricing plan 列表展示常量；筛选 value 与 `/pricing_plan/page` 等接口字段一致。
 */

/** @type {{ ACTIVE: 'Active', INACTIVE: 'Inactive' }} */
export const PRICING_PLAN_STATUS = {
	ACTIVE: 'Active',
	DISABLED: 'Disabled',
};

export const planStatusOptions = [
	{ label: 'Active', value: PRICING_PLAN_STATUS.ACTIVE },
	{ label: 'Disabled', value: PRICING_PLAN_STATUS.DISABLED },
];

/** @type {{ STANDARD: 'Standard', CUSTOM: 'Custom' }} */
export const PRICING_PLAN_TYPE = {
	STANDARD: 'Standard',
	CUSTOM: 'Custom',
};

/** 筛选下拉；value 与 API 字段 plan_type 一致 */
export const planTypeOptions = [
	{ label: 'Standard', value: PRICING_PLAN_TYPE.STANDARD },
	{ label: 'Custom', value: PRICING_PLAN_TYPE.CUSTOM },
];

/** Borrower / Lender / Admin 三端 type，与 filter.type、新增 body.type 一致 */
export const PRICING_PLAN_SYSTEM_TYPES = ['Borrower', 'Lender', 'Admin'];
