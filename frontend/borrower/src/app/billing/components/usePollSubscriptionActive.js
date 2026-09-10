'use client';

/**
 * 轮询 organization/subscription/current（via getCurrentPlan）。
 * subscription 为 null 表示尚未落库；有数据且 status=Active 时停止。
 */

import { useRef, useState } from 'react';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { useUserStore } from '@repo/shared-methods/store';
import { SUBSCRIPTION_STATUS } from '@repo/shared-methods/componentsUi/Pricing/constants';
import { normalizeSubscriptionStatus } from '@repo/shared-methods/componentsUi/Pricing/utils';

const DEFAULT_INTERVAL_MS = 2000;
/** webhook 通常数秒；约 60s 后停止，避免无限轮询 */
const DEFAULT_MAX_ATTEMPTS = 30;

/**
 * @param {object} [options]
 * @param {boolean} [options.enabled=true]
 * @param {number} [options.intervalMs]
 * @param {number} [options.maxAttempts]
 * @returns {{ phase: 'processing'|'active'|'timeout'|'error', errorMessage: string }}
 */
export default function usePollSubscriptionActive({
	enabled = true,
	intervalMs = DEFAULT_INTERVAL_MS,
	maxAttempts = DEFAULT_MAX_ATTEMPTS,
} = {}) {
	const getCurrentPlan = useUserStore((s) => s.getCurrentPlan);
	const [phase, setPhase] = useState('processing');
	const [errorMessage, setErrorMessage] = useState('');
	const attemptRef = useRef(0);
	const timerRef = useRef(null);
	const cancelledRef = useRef(false);

	const clearTimer = () => {
		if (timerRef.current != null) {
			clearTimeout(timerRef.current);
			timerRef.current = null;
		}
	};

	useDeferredSyncEffect({
		predicate: () => enabled,
		deferred: () => {
			cancelledRef.current = false;
			attemptRef.current = 0;
			setPhase('processing');
			setErrorMessage('');
			clearTimer();

			const pollOnce = async () => {
				if (cancelledRef.current) return;

				try {
					const subscription = await getCurrentPlan();
					if (cancelledRef.current) return;

					if (
						subscription &&
						normalizeSubscriptionStatus(subscription.status) ===
							SUBSCRIPTION_STATUS.ACTIVE
					) {
						setPhase('active');
						return;
					}

					attemptRef.current += 1;
					if (attemptRef.current >= maxAttempts) {
						setPhase('timeout');
						return;
					}

					timerRef.current = setTimeout(() => {
						void pollOnce();
					}, intervalMs);
				} catch (error) {
					if (cancelledRef.current) return;
					setPhase('error');
					setErrorMessage(error?.message || 'Unable to verify subscription status.');
				}
			};

			void pollOnce();
		},
		destroyCallback: () => {
			cancelledRef.current = true;
			clearTimer();
		},
		deps: [enabled, getCurrentPlan, intervalMs, maxAttempts],
		isUseMicrotask: true,
	});

	return { phase, errorMessage };
}
