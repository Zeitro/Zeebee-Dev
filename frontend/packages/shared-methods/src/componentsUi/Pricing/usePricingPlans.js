'use client';

import { useEffect, useState } from 'react';
import { useUserStore } from '@repo/shared-methods/store';
import { getPublicPricingPlans } from '../../api/pricing';
import { getCurrentPlanIdFromSubscription } from './utils';

/**
 * 拉取公开套餐列表，并在登录后拉取当前订阅。
 * enabled=false 时不请求套餐列表，但仍会在已登录时尝试 getCurrentPlan。
 *
 * @param {object} [options]
 * @param {boolean} [options.enabled=true] 是否请求公开套餐列表
 * @returns {{
 *   plans: object[],
 *   loading: boolean,
 *   currentPlan: object|null,
 *   currentPlanId: number|string|undefined,
 *   isLoggedIn: boolean,
 * }}
 */
export default function usePricingPlans({ enabled = true, isLoggedIn } = {}) {
	const { getCurrentPlan } = useUserStore();
	const [plans, setPlans] = useState([]);
	const [currentPlan, setCurrentPlan] = useState(null);
	const [loading, setLoading] = useState(Boolean(enabled));

	useEffect(() => {
		let cancelled = false;
		const load = async () => {
			if (enabled) {
				setLoading(true);
			}

			try {
				const tasks = [];
				if (enabled) {
					tasks.push(
						getPublicPricingPlans().then((res) => ({
							type: 'plans',
							data: Array.isArray(res?.plans) ? res.plans : [],
						})),
					);
				}
				if (isLoggedIn) {
					tasks.push(
						getCurrentPlan()
							.then((subscription) => ({
								type: 'currentPlan',
								data: subscription || null,
							}))
							.catch(() => ({
								type: 'currentPlan',
								data: null,
							})),
					);
				}

				if (tasks.length === 0) {
					return;
				}

				const results = await Promise.all(tasks);
				if (cancelled) {
					return;
				}

				results.forEach((result) => {
					if (result.type === 'plans') {
						setPlans(result.data);
					}
					if (result.type === 'currentPlan') {
						setCurrentPlan(result.data);
					}
				});
			} catch {
				if (!cancelled && enabled) {
					setPlans([]);
				}
			} finally {
				if (!cancelled && enabled) {
					setLoading(false);
				}
			}
		};

		load();

		return () => {
			cancelled = true;
		};
	}, [enabled, getCurrentPlan, isLoggedIn]);

	return {
		plans,
		loading,
		currentPlan,
		currentPlanId: getCurrentPlanIdFromSubscription(currentPlan),
		isLoggedIn,
	};
}
