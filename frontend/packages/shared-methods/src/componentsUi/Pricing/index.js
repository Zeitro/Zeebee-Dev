'use client';

import React, { useMemo, useState } from 'react';
import { Alert, Segmented, Spin } from 'antd';
import { NoData } from '@repo/shared-methods';
import PlanCard from './PlanCard';
import {
	BILLING_CYCLE,
	BILLING_CYCLE_OPTIONS,
	DEFAULT_REGISTER_PATH,
	DEFAULT_SALES_EMAIL,
	DEFAULT_SALES_PHONE,
} from './constants';
import usePricingActions from './usePricingActions';
import usePricingPlans from './usePricingPlans';
import { getSubscriptionExpiryAlert, isPurchasedSubscription, sortPlans } from './utils';

/**
 * 系统套餐定价展示组件；支持页面与弹框两种布局。
 * 默认请求 GET `/pricing_plans/public`；登录后额外请求 getCurrentPlan 获取当前订阅。
 * 也可由父组件传入 plans / currentPlanId / isLoggedIn 覆盖接口结果。
 * Custom 档 Contact Sales 展开 Popover（邮箱 + 电话），不再跳转 mailto。
 *
 * @param {object} props
 * @param {object[]} [props.plans] 套餐列表；传入则不再请求公开套餐接口
 * @param {'page'|'modal'} [props.variant] 布局变体
 * @param {string} [props.defaultBillingCycle] 默认计费周期
 * @param {string} [props.registerPath] 未登录跳转注册页
 * @param {number|string} [props.currentPlanId] 当前订阅 plan id；未传时使用 getCurrentPlan 结果
 * @param {boolean} [props.isLoggedIn] 是否已登录；未传时根据 token 自动判断
 * @param {(plan: object, policy: object|null) => void|Promise<void>} [props.onSelectPlan]
 * @param {(plan: object) => void} [props.onContactSales] Custom CTA 点击时可选回调
 * @param {string} [props.salesEmail] 销售邮箱
 * @param {string} [props.salesPhone] 销售电话
 * @param {string} [props.className]
 */
const Pricing = ({
	variant = 'page',
	defaultBillingCycle = BILLING_CYCLE.MONTHLY,
	registerPath = DEFAULT_REGISTER_PATH,
	currentPlanId,
	isLoggedIn,
	onSelectPlan,
	onContactSales,
	salesEmail = DEFAULT_SALES_EMAIL,
	salesPhone = DEFAULT_SALES_PHONE,
	className = '',
}) => {
	const [billingCycle, setBillingCycle] = useState(defaultBillingCycle);
	const {
		plans,
		loading,
		currentPlan,
		currentPlanId: fetchedCurrentPlanId,
	} = usePricingPlans({ enabled: true, isLoggedIn });
	const sortedPlans = useMemo(() => sortPlans(plans), [plans]);
	const resolvedCurrentPlanId = currentPlanId ?? fetchedCurrentPlanId;
	const isCurrentPlanPaid = isPurchasedSubscription(currentPlan);
	const expiryAlert = useMemo(
		() => getSubscriptionExpiryAlert(currentPlan, Boolean(isLoggedIn)),
		[currentPlan, isLoggedIn],
	);

	const { handlePlanAction, loadingPlanId } = usePricingActions({
		registerPath,
		onSelectPlan,
		onContactSales,
		billingCycle,
		isLoggedIn,
	});

	const isPageVariant = variant === 'page';

	return (
		<div
			className={[
				isPageVariant ? 'mx-auto w-full max-w-5xl px-4 py-8' : 'w-full',
				className,
			].join(' ')}
		>
			<div className={isPageVariant ? 'mb-8 text-center' : 'mb-6 text-center'}>
				{isPageVariant ? (
					<h2 className="m-0 text-[28px] font-medium leading-9 text-base-text">
						Choose your plan
					</h2>
				) : null}
				<p
					className={[
						'm-0 text-sm leading-6 text-base-gray',
						isPageVariant ? 'mt-2' : 'mt-0',
					].join(' ')}
				>
					Select the plan that fits your team. Switch or cancel anytime.
				</p>
				<div className="mt-5 flex justify-center">
					<Segmented
						value={billingCycle}
						options={BILLING_CYCLE_OPTIONS}
						onChange={setBillingCycle}
					/>
				</div>
			</div>

			{loading ? (
				<div className="flex min-h-[280px] items-center justify-center">
					<Spin />
				</div>
			) : sortedPlans.length === 0 ? (
				<NoData text="No pricing plans available" className="min-h-[280px]" />
			) : (
				<>
					<div className="grid grid-cols-1 gap-4 md:grid-cols-3 md:gap-5 mb-6">
						{sortedPlans.map((plan) => {
							const isCurrentPlan =
								resolvedCurrentPlanId != null &&
								String(resolvedCurrentPlanId) === String(plan.id);

							return (
								<PlanCard
									key={plan.id}
									plan={plan}
									billingCycle={billingCycle}
									isCurrentPlan={isCurrentPlan}
									isLoggedIn={isLoggedIn}
									isCurrentPlanPaid={isCurrentPlan && isCurrentPlanPaid}
									loading={
										loadingPlanId != null &&
										String(loadingPlanId) === String(plan.id)
									}
									onAction={handlePlanAction}
									salesEmail={salesEmail}
									salesPhone={salesPhone}
								/>
							);
						})}
					</div>
					{expiryAlert ? (
						<Alert
							className="text-left"
							message={expiryAlert.message}
							type="warning"
							showIcon
						/>
					) : null}
				</>
			)}
		</div>
	);
};

export default Pricing;
