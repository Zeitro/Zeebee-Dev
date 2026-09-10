/**
 * Pricing 套餐展示与 CTA 派生逻辑；纯函数，不改入参。
 */

import { formatDate, getRemainingDaysUntilExpire } from '../../lib/date.js';
import { PRICING_PLAN_TYPE, SUBSCRIPTION_STATUS } from './constants';

/** @param {object} plan */
export function isCustomPlan(plan) {
	return plan?.plan_type === PRICING_PLAN_TYPE.CUSTOM;
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
 * 是否有试用期（trial_days > 0）。
 * @param {object} plan
 * @returns {boolean}
 */
export function hasTrial(plan) {
	return Number(plan?.trial_days) > 0;
}

/**
 * 套餐 CTA 按钮文案。
 * - Custom → Contact Sales
 * - 已登录：当前档且已付费 → Your current plan；其余 → Choose plan
 * - 未登录：有试用 → Try Zeebee；否则 → Get + name
 *
 * @param {object} plan
 * @param {object} [options]
 * @param {boolean} [options.isLoggedIn] 是否已登录
 * @param {boolean} [options.isCurrentPlan] 是否当前订阅档
 * @param {boolean} [options.isCurrentPlanPaid] 当前档是否已付费（非试用）
 * @returns {string}
 */
export function getPlanCta(plan, options = {}) {
	const { isLoggedIn = false, isCurrentPlan = false, isCurrentPlanPaid = false } = options;
	if (isCustomPlan(plan)) {
		return 'Contact Sales';
	}
	if (isLoggedIn) {
		if (isCurrentPlan && isCurrentPlanPaid) {
			return 'Your current plan';
		}
		return 'Choose plan';
	}
	const planName = plan?.name || '';
	return hasTrial(plan) ? 'Try Zeebee' : `Get ${planName}`;
}

/**
 * 当前档是否锁定为「已购不可再点」。
 * @param {boolean} isCurrentPlan
 * @param {boolean} isCurrentPlanPaid
 * @returns {boolean}
 */
export function isCurrentPlanLocked(isCurrentPlan, isCurrentPlanPaid) {
	return Boolean(isCurrentPlan && isCurrentPlanPaid);
}

/**
 * 规范化订阅 status，兼容大小写差异（如 expiring / Expiring）。
 * @param {string|undefined|null} status
 * @returns {string}
 */
export function normalizeSubscriptionStatus(status) {
	if (status == null || status === '') {
		return '';
	}
	const raw = String(status).trim();
	const matched = Object.values(SUBSCRIPTION_STATUS).find(
		(item) => item.toLowerCase() === raw.toLowerCase(),
	);
	return matched || raw;
}

/**
 * 当前订阅是否处于试用期。
 * @param {object|null|undefined} subscription organization/subscription/current
 * @returns {boolean}
 */
export function isTrialingSubscription(subscription) {
	return normalizeSubscriptionStatus(subscription?.status) === SUBSCRIPTION_STATUS.TRIALING;
}

/**
 * 当前订阅是否已付费购买（非试用、未过期）。
 * subscription 为 null 表示未订阅；status 为 Trialing / Expired 视为未付费。
 * Active / Expiring（含大小写变体）视为仍在付费有效期内。
 * @param {object|null|undefined} subscription organization/subscription/current
 * @returns {boolean}
 */
export function isPurchasedSubscription(subscription) {
	if (!subscription) {
		return false;
	}
	const status = normalizeSubscriptionStatus(subscription.status);
	return status === SUBSCRIPTION_STATUS.ACTIVE || status === SUBSCRIPTION_STATUS.EXPIRING;
}

/**
 * 从当前订阅解析 pricing_plan_id。
 * @param {object|null|undefined} subscription
 * @returns {number|string|undefined}
 */
export function getCurrentPlanIdFromSubscription(subscription) {
	return subscription?.pricing_plan_id ?? undefined;
}

/**
 * 订阅/试用即将结束时的顶部提示文案；不满足条件返回 null。
 * - Trialing：试用即将结束
 * - Expiring 且 is_auto_renew=false：已付费订阅即将过期
 *
 * @param {object|null|undefined} subscription organization/subscription/current
 * @param {boolean} isLoggedIn
 * @returns {{ message: string }|null}
 */
export function getSubscriptionExpiryAlert(subscription, isLoggedIn) {
	if (!isLoggedIn || !subscription?.expire_date) {
		return null;
	}

	const remainingDays = getRemainingDaysUntilExpire(subscription.expire_date);
	if (remainingDays == null) {
		return null;
	}

	const dayLabel = remainingDays === 1 ? '1 day' : `${remainingDays} days`;
	const status = normalizeSubscriptionStatus(subscription.status);

	if (status === SUBSCRIPTION_STATUS.TRIALING) {
		return {
			message: `Your free trial is ending soon. ${dayLabel} remaining.`,
		};
	}

	if (status === SUBSCRIPTION_STATUS.EXPIRING && subscription.is_auto_renew === false) {
		const expireDateFormatted = formatDate(subscription.expire_date);
		return {
			message: `Your subscription is expiring on ${expireDateFormatted}. ${dayLabel} remaining.`,
		};
	}

	return null;
}

/**
 * Pro 推荐档高亮：is_default 或 name 为 Pro。
 * @param {object} plan
 * @returns {boolean}
 */
export function isHighlightedPlan(plan) {
	if (plan?.is_recommend) return true;
	return false;
}
