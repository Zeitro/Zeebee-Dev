'use client';

import { useCallback, useState } from 'react';
import { message } from 'antd';
import { useRouter } from 'next/navigation';
import { useUserStore } from '@repo/shared-methods/store';
import { createSubscriptionCheckout } from '../../api/pricing';
import { DEFAULT_REGISTER_PATH } from './constants';
import { getPolicyByCycle, isCustomPlan } from './utils';

/**
 * Pricing CTA 行为：未登录跳转注册；已登录且 owner/admin 调 checkout 并跳转 Stripe；
 * Custom 由 PlanCard Popover 展示联系方式，此处仅转发 onContactSales。
 *
 * @param {object} options
 * @param {string} [options.registerPath]
 * @param {(plan: object, policy: object|null) => void|Promise<void>} [options.onSelectPlan]
 * @param {(plan: object) => void} [options.onContactSales]
 * @param {string} [options.billingCycle]
 * @param {boolean} [options.isLoggedIn]
 */
export default function usePricingActions({
	registerPath = DEFAULT_REGISTER_PATH,
	onSelectPlan,
	onContactSales,
	billingCycle,
	isLoggedIn,
}) {
	const router = useRouter();
	const userInfo = useUserStore((state) => state.userInfo);
	const [loadingPlanId, setLoadingPlanId] = useState(null);

	const handleContactSales = useCallback(
		(plan) => {
			onContactSales?.(plan);
		},
		[onContactSales],
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
				const { checkout_url: checkoutUrl } = await createSubscriptionCheckout({
					pricing_plan_id: plan.id,
					pricing_plan_policy_id: policy.id,
				});
				if (!checkoutUrl) {
					message.error('Unable to start checkout. Please try again later.');
					return;
				}
				// Stripe Hosted Checkout；支付成功/取消由后端配置的 success/cancel URL 回跳
				const popupWidth = 500;
				const popupHeight = 600;
				const left = (screen.availWidth - popupWidth) / 2;
				const top = (screen.availHeight - popupHeight) / 2;

				window.open(
					checkoutUrl,
					'_blank',
					`width=${popupWidth},height=${popupHeight},left=${left},top=${top},resizable=yes,scrollbars=yes`,
				);
			} finally {
				setLoadingPlanId(null);
			}
		},
		[onSelectPlan],
	);

	const handlePlanAction = useCallback(
		async (plan) => {
			if (isCustomPlan(plan)) {
				handleContactSales(plan);
				return;
			}

			if (!isLoggedIn) {
				router.push(registerPath);
				return;
			}

			// 与 PlanCard / 后端约定一致：仅组织 owner / admin 可发起支付
			if (!(userInfo?.is_owner || userInfo?.is_admin)) {
				message.warning('Only the organization owner or admin can subscribe.');
				return;
			}

			const policy = getPolicyByCycle(plan, billingCycle);
			await handleCheckout(plan, policy);
		},
		[
			billingCycle,
			handleCheckout,
			handleContactSales,
			isLoggedIn,
			registerPath,
			router,
			userInfo?.is_admin,
			userInfo?.is_owner,
		],
	);

	return {
		handlePlanAction,
		loadingPlanId,
	};
}
