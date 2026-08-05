'use client';

import React, { useMemo, useState } from 'react';
import { Segmented } from 'antd';
import PlanCard from './PlanCard';
import {
	BILLING_CYCLE,
	BILLING_CYCLE_OPTIONS,
	DEFAULT_CHECKOUT_API,
	DEFAULT_CONTACT_SALES_HREF,
	DEFAULT_REGISTER_PATH,
	MOCK_PRICING_PLANS,
} from './constants';
import usePricingActions from './usePricingActions';
import { sortPlans } from './utils';

/**
 * 系统套餐定价展示组件；支持页面与弹框两种布局。
 * 当前使用静态 mock 数据，后续可由父组件传入 plans 或接入接口。
 *
 * @param {object} props
 * @param {object[]} [props.plans] 套餐列表，默认 MOCK_PRICING_PLANS
 * @param {'page'|'modal'} [props.variant] 布局变体
 * @param {string} [props.defaultBillingCycle] 默认计费周期
 * @param {string} [props.registerPath] 未登录跳转注册页
 * @param {string} [props.checkoutApi] 已登录 checkout 接口
 * @param {number|string} [props.currentPlanId] 当前订阅 plan id
 * @param {(plan: object, policy: object|null) => void|Promise<void>} [props.onSelectPlan]
 * @param {(plan: object) => void} [props.onContactSales]
 * @param {string} [props.contactSalesHref]
 * @param {string} [props.className]
 */
const Pricing = ({
	plans = MOCK_PRICING_PLANS,
	variant = 'page',
	defaultBillingCycle = BILLING_CYCLE.MONTHLY,
	registerPath = DEFAULT_REGISTER_PATH,
	checkoutApi = DEFAULT_CHECKOUT_API,
	currentPlanId,
	onSelectPlan,
	onContactSales,
	contactSalesHref = DEFAULT_CONTACT_SALES_HREF,
	className = '',
}) => {
	const [billingCycle, setBillingCycle] = useState(defaultBillingCycle);
	const sortedPlans = useMemo(() => sortPlans(plans), [plans]);

	const { handlePlanAction, loadingPlanId } = usePricingActions({
		registerPath,
		checkoutApi,
		contactSalesHref,
		onSelectPlan,
		onContactSales,
		billingCycle,
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

			<div className="grid grid-cols-1 gap-4 md:grid-cols-3 md:gap-5">
				{sortedPlans.map((plan) => (
					<PlanCard
						key={plan.id}
						plan={plan}
						billingCycle={billingCycle}
						isCurrentPlan={
							currentPlanId != null && String(currentPlanId) === String(plan.id)
						}
						loading={loadingPlanId != null && String(loadingPlanId) === String(plan.id)}
						onAction={handlePlanAction}
					/>
				))}
			</div>
		</div>
	);
};

export default Pricing;
