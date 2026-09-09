/**
 * 套餐定价相关接口；路径相对 `/api/v1`（由 createRequest 拼接）。
 */

import { createRequest } from '../lib/fetch';

const { request } = createRequest();

/**
 * 获取公开套餐列表（含 policies）。
 * GET `/pricing_plans/public` → `{ plans: object[] }`
 *
 * @returns {Promise<{ plans?: object[] }>}
 */
export const getPublicPricingPlans = () => {
	return request('/pricing_plans/public', {
		method: 'GET',
	});
};

/**
 * 发起 Stripe Checkout（仅组织 owner；目标 policy 须 actual_amount > 0）。
 * POST `/organization/subscription/checkout`
 *
 * @param {{ pricing_plan_id: number, pricing_plan_policy_id: number }} body
 * @returns {Promise<{ checkout_url?: string }>}
 */
export const createSubscriptionCheckout = (body) => {
	return request('/organization/subscription/checkout', {
		method: 'POST',
		body,
	});
};

/**
 * 取消自动续费（仅组织 owner）。
 * POST `/organization/subscription/cancel_renewal`，body `{}`
 *
 * @returns {Promise<object>}
 */
export const cancelSubscriptionRenewal = () => {
	return request('/organization/subscription/cancel_renewal', {
		method: 'POST',
		body: {},
	});
};

/**
 * 恢复自动续费（仅组织 owner；须在当前周期到期前）。
 * POST `/organization/subscription/resume_renewal`，body `{}`
 *
 * @returns {Promise<object>}
 */
export const resumeSubscriptionRenewal = () => {
	return request('/organization/subscription/resume_renewal', {
		method: 'POST',
		body: {},
	});
};
