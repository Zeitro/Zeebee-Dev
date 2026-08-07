'use client';

/**
 * Stripe Checkout 支付成功回跳页。
 * 订阅由 webhook 异步落库，需轮询 current 直至 status=Active；
 * Active 后按登录流程刷新菜单 / organization / currentPlan，再延迟关闭窗口。
 */

import { useRef } from 'react';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import { useUserStore, usePermissionStore } from '@repo/shared-methods/store';
import { menuSetting } from '@config/routes';
import BillingCallbackView from '../components/BillingCallbackView';
import usePollSubscriptionActive from '../components/usePollSubscriptionActive';

/** 与 driveAuth / gmailAuth 一致：成功态展示后再关窗 */
const CLOSE_DELAY_MS = 500;

const CONTENT_BY_PHASE = {
	processing: {
		status: 'processing',
		title: 'Activating your subscription…',
		description: 'Payment received. This usually takes a few seconds.',
	},
	active: {
		status: 'success',
		title: 'Subscription activated',
		description: 'Your plan is ready. This window will close shortly.',
	},
	timeout: {
		status: 'error',
		title: 'Still processing',
		description:
			'Payment went through, but activation is taking longer than expected. Please check again shortly.',
	},
	error: {
		status: 'error',
		title: 'Unable to confirm subscription',
		description: 'Please try again later or contact support.',
	},
};

export default function BillingSuccessPage() {
	const { phase, errorMessage } = usePollSubscriptionActive();
	const { getUserMenus, getOrganization, getCurrentPlan } = useUserStore();
	const { setPermissionTags, setAccessPaths } = usePermissionStore();
	const closeTimerRef = useRef(null);
	const content = CONTENT_BY_PHASE[phase] ?? CONTENT_BY_PHASE.error;

	useDeferredSyncEffect({
		predicate: () => phase === 'active',
		deferred: () => {
			const refreshAndClose = async () => {
				try {
					// 与登录 useData 一致：菜单权限 + organization + 最新 currentPlan
					const result = await getUserMenus(menuSetting);
					await getCurrentPlan();
					setPermissionTags(result?.permissionTags);
					setAccessPaths(result?.accessPaths ?? []);
				} catch {
					// 刷新失败仍关闭，避免成功页卡住
				}
				closeTimerRef.current = setTimeout(() => {
					window.close();
				}, CLOSE_DELAY_MS);
			};
			void refreshAndClose();
		},
		destroyCallback: () => {
			if (closeTimerRef.current != null) {
				clearTimeout(closeTimerRef.current);
				closeTimerRef.current = null;
			}
		},
		deps: [
			phase,
			getUserMenus,
			getOrganization,
			getCurrentPlan,
			setPermissionTags,
			setAccessPaths,
		],
		isUseMicrotask: true,
	});

	return (
		<BillingCallbackView
			status={content.status}
			title={content.title}
			description={phase === 'error' && errorMessage ? errorMessage : content.description}
		/>
	);
}
