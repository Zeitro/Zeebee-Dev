'use client';

import { useCallback, useMemo, useState } from 'react';
import { message } from 'antd';
import { useRouter } from 'next/navigation';
import { createRequest, getAccessTokenFromStorage } from '@repo/shared-methods/utils';
import { useUserStore } from '@repo/shared-methods/store';
import {
	DEFAULT_CHECKOUT_API,
	DEFAULT_CONTACT_SALES_HREF,
	DEFAULT_REGISTER_PATH,
} from './constants';
import { getPolicyByCycle, isCustomPlan } from './utils';

/**
 * Pricing CTA 行为：未登录跳转注册；已登录调 checkout；Custom 联系销售。
 *
 * @param {object} options
 * @param {string} [options.registerPath]
 * @param {string} [options.checkoutApi]
 * @param {string} [options.contactSalesHref]
 * @param {(plan: object, policy: object|null) => void|Promise<void>} [options.onSelectPlan]
 * @param {(plan: object) => void} [options.onContactSales]
 * @param {string} [options.billingCycle]
 */
export default function usePricingActions({
	registerPath = DEFAULT_REGISTER_PATH,
	checkoutApi = DEFAULT_CHECKOUT_API,
	contactSalesHref = DEFAULT_CONTACT_SALES_HREF,
	onSelectPlan,
	onContactSales,
	billingCycle,
}) {
	const router = useRouter();
	const userInfo = useUserStore((state) => state.userInfo);
	const { request } = useMemo(() => createRequest(), []);
	const [loadingPlanId, setLoadingPlanId] = useState(null);

	const isLoggedIn = useCallback(() => {
		return Boolean(getAccessTokenFromStorage() || userInfo?.id);
	}, [userInfo?.id]);

	const handleContactSales = useCallback(
		(plan) => {
			if (onContactSales) {
				onContactSales(plan);
				return;
			}
			window.location.href = contactSalesHref;
		},
		[contactSalesHref, onContactSales],
	);

	const handleCheckout = useCallback(
		async (plan, policy) => {
			if (onSelectPlan) {
				await onSelectPlan(plan, policy);
				return;
			}

			if (!policy?.id) {
				message.error('Pricing is not available for the selected billing cycle.');
				return;
			}

			setLoadingPlanId(plan.id);
			try {
				await request(checkoutApi, {
					method: 'POST',
					body: {
						pricing_plan_id: plan.id,
						pricing_plan_policy_id: policy.id,
					},
				});
				message.success('Redirecting to checkout...');
			} catch {
				message.error('Unable to start checkout. Please try again later.');
			} finally {
				setLoadingPlanId(null);
			}
		},
		[checkoutApi, onSelectPlan, request],
	);

	const handlePlanAction = useCallback(
		async (plan) => {
			if (isCustomPlan(plan)) {
				handleContactSales(plan);
				return;
			}

			if (!isLoggedIn()) {
				router.push(registerPath);
				return;
			}

			const policy = getPolicyByCycle(plan, billingCycle);
			await handleCheckout(plan, policy);
		},
		[billingCycle, handleCheckout, handleContactSales, isLoggedIn, registerPath, router],
	);

	return {
		handlePlanAction,
		loadingPlanId,
	};
}
