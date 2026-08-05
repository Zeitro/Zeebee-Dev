'use client';

import React from 'react';
import { Button } from 'antd';
import { displayValue, formatAmount } from '@repo/shared-methods/utils';
import { getPlanCta, getPolicyByCycle, hasTrial, isCustomPlan, isHighlightedPlan } from './utils';

/**
 * 单档套餐卡片：价格、描述、试用提示与 CTA。
 *
 * @param {object} props
 * @param {object} props.plan 套餐数据（与 admin pricing_plan 字段一致）
 * @param {string} props.billingCycle 当前计费周期 Monthly / Yearly
 * @param {boolean} [props.isCurrentPlan] 是否为当前订阅
 * @param {boolean} [props.loading] CTA loading
 * @param {(plan: object) => void} props.onAction 点击 CTA
 */
const PlanCard = ({ plan, billingCycle, isCurrentPlan = false, loading = false, onAction }) => {
	const isCustom = isCustomPlan(plan);
	const isHighlighted = isHighlightedPlan(plan);
	const policy = getPolicyByCycle(plan, billingCycle);
	const ctaLabel = isCurrentPlan ? 'Current plan' : getPlanCta(plan);

	const renderPrice = () => {
		if (isCustom) {
			return (
				<div className="mt-6 min-h-[72px]">
					<p className="m-0 text-2xl font-medium leading-8 text-base-text">
						Custom pricing
					</p>
					<p className="m-0 mt-1 text-sm text-base-label">Talk to our sales team</p>
				</div>
			);
		}

		const actualAmount = displayValue(policy?.actual_amount, formatAmount);
		const originalAmount = displayValue(policy?.original_amount, formatAmount);
		const showOriginal =
			policy?.original_amount != null &&
			policy?.actual_amount != null &&
			Number(policy.original_amount) > Number(policy.actual_amount) &&
			originalAmount !== actualAmount;

		return (
			<div className="mt-6 min-h-[72px]">
				<div className="flex flex-wrap items-end gap-x-2 gap-y-1">
					<span className="text-[32px] font-medium leading-10 text-base-text">
						{actualAmount}
					</span>
					<span className="pb-1 text-sm text-base-label">/ month</span>
				</div>
				{showOriginal ? (
					<p className="m-0 mt-1 text-sm text-base-label line-through">
						{originalAmount}
					</p>
				) : null}
			</div>
		);
	};

	return (
		<div
			className={[
				'flex h-full flex-col rounded-xl border bg-white p-6',
				isHighlighted ? 'border-base-primary shadow-sm' : 'border-base-border',
			].join(' ')}
		>
			{isHighlighted ? (
				<span className="mb-3 inline-flex w-fit rounded-full bg-base-active px-3 py-0.5 text-xs font-medium text-base-primary">
					Most popular
				</span>
			) : (
				<span className="mb-3 inline-flex h-[22px]" aria-hidden />
			)}

			<h3 className="m-0 text-xl font-medium leading-7 text-base-text">
				{displayValue(plan?.name)}
			</h3>
			<p className="m-0 mt-1 text-sm leading-5 text-base-gray">
				{displayValue(plan?.subtitle)}
			</p>

			{renderPrice()}

			<p className="m-0 mt-4 flex-1 text-sm leading-6 text-base-gray">
				{displayValue(plan?.description)}
			</p>

			{hasTrial(plan) ? (
				<p className="m-0 mt-4 text-sm font-medium text-base-primary">
					{plan.trial_days}-day free trial
				</p>
			) : null}

			<Button
				type={isHighlighted ? 'primary' : 'default'}
				size="large"
				block
				className="mt-4"
				loading={loading}
				disabled={isCurrentPlan}
				onClick={() => onAction?.(plan)}
			>
				{ctaLabel}
			</Button>
		</div>
	);
};

export default PlanCard;
