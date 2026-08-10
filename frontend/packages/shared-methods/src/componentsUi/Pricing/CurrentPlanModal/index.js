'use client';

import React, { useCallback, useRef, useState } from 'react';
import { Button, Modal, Spin, Tooltip, message } from 'antd';
import {
	CheckCircleFilled,
	ClockCircleOutlined,
	ExclamationCircleOutlined,
	SyncOutlined,
} from '@ant-design/icons';
import { NewModal } from '@repo/shared-methods';
import { displayValue, formatAmount, formatDate } from '@repo/shared-methods/utils';
import { useUserStore } from '@repo/shared-methods/store';
import { useDeferredSyncEffect } from '../../../hooks';
import { cancelSubscriptionRenewal, resumeSubscriptionRenewal } from '../../../api/pricing';
import { SUBSCRIPTION_STATUS } from '../constants';
import {
	isPurchasedSubscription,
	isTrialingSubscription,
	normalizeSubscriptionStatus,
} from '../utils';

/**
 * 订阅详情一行。
 * @param {{ label: string, value: React.ReactNode }} props
 */
const DetailRow = ({ label, value }) => (
	<div className="flex items-start justify-between gap-4 py-2.5">
		<span className="shrink-0 text-sm text-base-label">{label}</span>
		<span className="text-right text-sm font-medium text-base-text">{value}</span>
	</div>
);

/**
 * 将 API status 映射为用户可读的类型文案与徽章样式。
 * @param {string|undefined} status
 * @returns {{ label: string, hint: string, badgeClass: string, icon: React.ReactNode }|null}
 */
function getStatusMeta(status) {
	switch (normalizeSubscriptionStatus(status)) {
		case SUBSCRIPTION_STATUS.TRIALING:
			return {
				label: 'Free trial',
				hint: 'You have not been charged yet.',
				badgeClass:
					'border-[var(--orange-medium)] bg-[var(--orange-light)] text-[var(--orange-dark)]',
				icon: <ClockCircleOutlined className="text-[12px]" />,
			};
		case SUBSCRIPTION_STATUS.ACTIVE:
			return {
				label: 'Paid plan',
				hint: 'Your subscription is active and paid.',
				badgeClass:
					'border-[var(--green-light-100)] bg-[var(--green-light)] text-[var(--green-medium)]',
				icon: <CheckCircleFilled className="text-[12px]" />,
			};
		case SUBSCRIPTION_STATUS.EXPIRING:
			return {
				label: 'Paid · Ending soon',
				hint: 'Access continues until the expire date below. Resume auto-renewal to keep this plan.',
				badgeClass:
					'border-[var(--orange-medium)] bg-[var(--orange-light)] text-[var(--orange-dark)]',
				icon: <ExclamationCircleOutlined className="text-[12px]" />,
			};
		case SUBSCRIPTION_STATUS.EXPIRED:
			return {
				label: 'Expired',
				hint: 'This plan has ended. Upgrade to continue.',
				badgeClass: 'border-base-border bg-gray-light text-base-gray',
				icon: <ExclamationCircleOutlined className="text-[12px]" />,
			};
		default:
			return status
				? {
						label: status,
						hint: '',
						badgeClass: 'border-base-border bg-gray-light text-base-gray',
						icon: null,
					}
				: null;
	}
}

/**
 * 当前订阅详情弹框。
 * 必展示过期时间；已付费档额外展示费用、是否自动续费，以及取消/恢复自动续费。
 * 试用期不展示已付费相关字段与续费操作；状态用「Free trial / Paid plan」直观区分。
 *
 * @param {object} props
 * @param {boolean} props.open 是否展示
 * @param {() => void} props.onCancel 关闭回调
 * @param {string} [props.title] 弹框标题
 * @param {number} [props.width] 弹框宽度
 * @param {object} [props.subscription] 覆盖 store 中的 currentPlan
 * @param {() => void} [props.onUpgrade] 「Upgrade plan」回调
 */
const CurrentPlanModal = ({
	open,
	onCancel,
	title = 'Current plan',
	width = 440,
	subscription: subscriptionProp,
	onUpgrade,
}) => {
	const userInfo = useUserStore((state) => state.userInfo);
	const storePlan = useUserStore((state) => state.currentPlan);
	const getCurrentPlan = useUserStore((state) => state.getCurrentPlan);
	const [loading, setLoading] = useState(false);
	const [actionLoading, setActionLoading] = useState(false);
	const [subscription, setSubscription] = useState(() => subscriptionProp || storePlan || null);
	const [messageApi, contextHolder] = message.useMessage();
	const cancelledRef = useRef(false);

	const isOwnerAdmin = Boolean(userInfo?.is_owner || userInfo?.is_admin);
	const isTrialing = isTrialingSubscription(subscription);
	const isPurchased = isPurchasedSubscription(subscription);
	const isAutoRenew = subscription?.is_auto_renew === true;
	/** Active / Expiring：周期内仍可用，展示取消/恢复自动续费，不引导换套餐 */
	const showRenewActions = isPurchased && !isTrialing;

	// 弹框打开时拉取最新订阅；setState 放在 microtask，避免 effect 同步 setState 级联渲染
	useDeferredSyncEffect({
		predicate: () => open,
		deferred: () => {
			cancelledRef.current = false;
			setLoading(true);
			getCurrentPlan()
				.then((next) => {
					if (!cancelledRef.current) {
						setSubscription(next || null);
					}
				})
				.finally(() => {
					if (!cancelledRef.current) {
						setLoading(false);
					}
				});
		},
		destroyCallback: () => {
			cancelledRef.current = true;
		},
		// 仅在弹框打开时拉取；勿依赖 storePlan，避免 getCurrentPlan 写 store 后循环请求
		deps: [open, getCurrentPlan, messageApi],
		isUseMicrotask: true,
	});

	const refreshSubscription = useCallback(async () => {
		const next = await getCurrentPlan();
		setSubscription(next || null);
		return next;
	}, [getCurrentPlan]);

	const handleCancelRenewal = useCallback(() => {
		Modal.confirm({
			title: 'Cancel auto-renewal?',
			icon: <ExclamationCircleOutlined />,
			centered: true,
			content:
				'Your plan will stay active until the current period ends. After that, it will fall back to the free default plan and you will not be charged again.',
			okText: 'Cancel auto-renewal',
			cancelText: 'Keep auto-renewal',
			okButtonProps: { danger: true },
			onOk: async () => {
				setActionLoading(true);
				try {
					await cancelSubscriptionRenewal();
					await refreshSubscription();
					messageApi.success('Auto-renewal has been turned off.');
				} catch (error) {
					messageApi.error('Unable to cancel auto-renewal. Please try again.');
					return Promise.reject(error);
				} finally {
					setActionLoading(false);
				}
			},
		});
	}, [messageApi, refreshSubscription]);

	const handleResumeRenewal = useCallback(async () => {
		setActionLoading(true);
		try {
			await resumeSubscriptionRenewal();
			await refreshSubscription();
			messageApi.success('Auto-renewal has been turned back on.');
		} finally {
			setActionLoading(false);
		}
	}, [messageApi, refreshSubscription]);

	if (!open) {
		return null;
	}

	const planName = subscription?.pricing_plan?.name;
	const expireDate = displayValue(subscription?.expire_date, formatDate);
	const statusMeta = getStatusMeta(subscription?.status);
	const billingCycle =
		subscription?.pricing_plan_policy_type || subscription?.pricing_plan_policy?.type;
	const headerToneClass = isTrialing
		? 'border-[var(--orange-medium)] bg-[var(--orange-light)]'
		: isPurchased
			? 'border-[var(--green-light-100)] bg-[var(--green-light)]'
			: 'border-base-border bg-gray-light';

	const renewButton = isAutoRenew ? (
		<Button
			block
			size="large"
			danger
			loading={actionLoading}
			disabled={!isOwnerAdmin}
			onClick={handleCancelRenewal}
		>
			Cancel auto-renewal
		</Button>
	) : (
		<Button
			block
			size="large"
			type="primary"
			loading={actionLoading}
			disabled={!isOwnerAdmin}
			onClick={handleResumeRenewal}
			icon={<SyncOutlined />}
		>
			Resume auto-renewal
		</Button>
	);

	return (
		<NewModal title={title} width={width} onCancel={onCancel} footer={null} destroyOnHidden>
			{contextHolder}
			{loading ? (
				<div className="flex min-h-[200px] items-center justify-center">
					<Spin />
				</div>
			) : (
				<div className="flex flex-col gap-4">
					{/* 套餐头图：名称 + 试用/已付费徽章 */}
					<div className={['rounded-xl border p-4', headerToneClass].join(' ')}>
						<div className="flex items-start justify-between gap-3">
							<div className="min-w-0">
								<p className="m-0 text-xs font-medium uppercase tracking-wide text-base-label">
									Your plan
								</p>
								<h3 className="m-0 mt-1 truncate text-xl font-medium leading-7 text-base-text">
									{displayValue(planName)}
								</h3>
							</div>
							{statusMeta ? (
								<span
									className={[
										'inline-flex shrink-0 items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-medium leading-[18px]',
										statusMeta.badgeClass,
									].join(' ')}
								>
									{statusMeta.icon}
									{statusMeta.label}
								</span>
							) : null}
						</div>
						{statusMeta?.hint ? (
							<p className="m-0 mt-2 text-sm leading-5 text-base-gray">
								{statusMeta.hint}
							</p>
						) : null}
						<p className="m-0 mt-3 text-sm text-base-text">
							<span className="text-base-label">
								{isTrialing ? 'Trial ends on' : 'Expires on'}
							</span>{' '}
							<span className="font-medium">{expireDate}</span>
						</p>
					</div>

					{/* 已付费详情 */}
					{!isTrialing ? (
						<div className="rounded-xl border border-base-border bg-white px-4 py-1">
							<DetailRow
								label="Amount paid"
								value={displayValue(subscription?.actual_amount, formatAmount)}
							/>
							{billingCycle ? (
								<DetailRow
									label="Billing cycle"
									value={displayValue(billingCycle)}
								/>
							) : null}
							<DetailRow
								label="Auto-renewal"
								value={
									<span
										className={
											isAutoRenew
												? 'text-[var(--green-medium)]'
												: 'text-base-gray'
										}
									>
										{isAutoRenew ? 'On' : 'Off'}
									</span>
								}
							/>
						</div>
					) : (
						<p className="m-0 text-sm leading-5 text-base-gray">
							Billing amount and auto-renewal only apply after you upgrade to a paid
							plan.
						</p>
					)}

					{/* 取消 / 恢复自动续费 */}
					{showRenewActions ? (
						<div className="rounded-xl border border-base-border bg-base-active p-4">
							<p className="m-0 text-sm font-medium text-base-text">Auto-renewal</p>
							<p className="m-0 mt-1 text-xs leading-5 text-base-gray">
								{isAutoRenew
									? `Your plan renews automatically. Cancel anytime — access continues until ${expireDate}.`
									: `Auto-renewal is off. Resume before ${expireDate} to keep your plan without interruption.`}
							</p>
							<div className="mt-3">
								{isOwnerAdmin ? (
									renewButton
								) : (
									<Tooltip title="Only the organization owner or admin can manage auto-renewal">
										<span className="block">{renewButton}</span>
									</Tooltip>
								)}
							</div>
						</div>
					) : null}

					{isTrialing && onUpgrade ? (
						<Button block type="primary" size="large" onClick={onUpgrade}>
							Upgrade plan
						</Button>
					) : null}

					{!isTrialing && !isPurchased && onUpgrade ? (
						<Button block type="primary" size="large" onClick={onUpgrade}>
							Choose a plan
						</Button>
					) : null}
				</div>
			)}
		</NewModal>
	);
};

export default CurrentPlanModal;
