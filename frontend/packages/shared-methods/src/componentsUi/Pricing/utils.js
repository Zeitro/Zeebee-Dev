/**
 * Pricing 套餐展示与 CTA 派生逻辑；纯函数，不改入参。
 */

import { PRICING_PLAN_TYPE } from './constants';

/** @param {object} plan */
export function isCustomPlan(plan) {
	return plan?.plan_type === PRICING_PLAN_TYPE.CUSTOM;
}

/** @param {object} plan */
export function isBasePlan(plan) {
	return String(plan?.name || '').toLowerCase() === 'base';
}

/**
 * 按 sort_order 升序排列套餐。
 * @param {object[]} plans
 * @returns {object[]}
 */
export function sortPlans(plans = []) {
	return [...plans].sort((a, b) => (a?.sort_order ?? 0) - (b?.sort_order ?? 0));
}

/**
 * 取指定计费周期对应的 policy。
 * @param {object} plan
 * @param {string} billingCycle policies[].type，如 Monthly / Yearly
 * @returns {object|null}
 */
export function getPolicyByCycle(plan, billingCycle) {
	return plan?.policies?.find((item) => item?.type === billingCycle) ?? null;
}

/**
 * 套餐 CTA 按钮文案。
 * Base + trial_days > 0 → Try Zeebee；Custom → Contact Sales；其余 → Get Pro。
 * @param {object} plan
 * @returns {string}
 */
export function getPlanCta(plan) {
	if (isCustomPlan(plan)) {
		return 'Contact Sales';
	}
	if (isBasePlan(plan) && Number(plan?.trial_days) > 0) {
		return 'Try Zeebee';
	}
	return 'Get Pro';
}

/**
 * 是否展示试用提示。
 * @param {object} plan
 * @returns {boolean}
 */
export function hasTrial(plan) {
	return isBasePlan(plan) && Number(plan?.trial_days) > 0;
}

/**
 * Pro 推荐档高亮：is_default 或 name 为 Pro。
 * @param {object} plan
 * @returns {boolean}
 */
export function isHighlightedPlan(plan) {
	if (plan?.is_default) return true;
	return String(plan?.name || '').toLowerCase() === 'pro';
}
