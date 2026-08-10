'use client';

import React from 'react';
import { Button, Popover, Tooltip } from 'antd';
import { CheckCircleFilled, ClockCircleOutlined } from '@ant-design/icons';
import { useUserStore } from '@repo/shared-methods/store';
import { displayValue, formatAmount } from '@repo/shared-methods/utils';
import ContactSalesContent from './ContactSalesContent';
import { DEFAULT_SALES_EMAIL, DEFAULT_SALES_PHONE } from './constants';
import {
	getPlanCta,
	getPolicyByCycle,
	hasTrial,
	isCurrentPlanLocked,
	isCustomPlan,
	isHighlightedPlan,
} from './utils';

/** 与 CurrentPlanModal 一致：仅组织 owner / admin 可发起付费订阅 */
const BILLING_RESTRICTED_TIP = 'Only the organization owner or admin can subscribe';

/**
 * 单档套餐卡片：价格、描述、试用提示与 CTA。
 * Custom 档 CTA 点击展开 Contact Sales Popover（邮箱 + 电话），不跳转 mailto。
 * 已登录时，付费 CTA 仅 owner / admin 可点；其余角色禁用并 Tooltip 提示。
 *
 * @param {object} props
 * @param {object} props.plan 套餐数据（与 pricing_plans/public 字段一致）
 * @param {string} props.billingCycle 当前计费周期 Monthly / Yearly
 * @param {boolean} [props.isCurrentPlan] 是否为当前订阅档
 * @param {boolean} [props.isLoggedIn] 是否已登录
 * @param {boolean} [props.isCurrentPlanPaid] 当前档是否已付费（由订阅 status 推导）
 * @param {boolean} [props.loading] CTA loading
 * @param {(plan: object) => void} props.onAction 点击 CTA
 * @param {string} [props.salesEmail] 销售邮箱（Custom Popover）
 * @param {string} [props.salesPhone] 销售电话（Custom Popover）
 */
const PlanCard = ({
	plan,
	billingCycle,
	isCurrentPlan = false,
	isLoggedIn = false,
	isCurrentPlanPaid = false,
	loading = false,
	onAction,
	salesEmail = DEFAULT_SALES_EMAIL,
	salesPhone = DEFAULT_SALES_PHONE,
}) => {
	const userInfo = useUserStore((state) => state.userInfo);
	const isOwnerAdmin = Boolean(userInfo?.is_owner || userInfo?.is_admin);
	/** 已登录且非 owner/admin：不可发起支付（未登录仍可点，走注册） */
	const isBillingRestricted = Boolean(isLoggedIn) && !isOwnerAdmin;
	const isCustom = isCustomPlan(plan);
	const isHighlighted = isHighlightedPlan(plan);
	const isLocked = isCurrentPlanLocked(isCurrentPlan, isCurrentPlanPaid);
	const isPayDisabled = isLocked || isBillingRestricted;
	const policy = getPolicyByCycle(plan, billingCycle);
	const ctaLabel = getPlanCta(plan, { isLoggedIn, isCurrentPlan, isCurrentPlanPaid });

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

		return (
			<div className="mt-6 min-h-[72px]">
				<div className="flex flex-wrap items-end gap-x-2 gap-y-1">
					<span className="text-[32px] font-medium leading-10 text-base-text">
						{actualAmount}
					</span>
					<span className="pb-1 text-sm text-base-label">/ month</span>
				</div>
				{originalAmount ? (
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
				'flex h-full flex-col rounded-xl border p-6',
				isCurrentPlan
					? 'border-base-primary bg-base-active shadow-sm'
					: isHighlighted
						? 'border-base-primary bg-white shadow-sm'
						: 'border-base-border bg-white',
			].join(' ')}
		>
			<div className="mb-3 inline-flex h-[22px] gap-x-1">
				{isCurrentPlan && (
					<span
						className={[
							'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium leading-[18px]',
							isCurrentPlanPaid
								? 'border-[var(--green-light-100)] bg-[var(--green-light)] text-[var(--green-medium)]'
								: 'border-[var(--orange-medium)] bg-[var(--orange-light)] text-[var(--orange-dark)]',
						].join(' ')}
					>
						{isCurrentPlanPaid ? (
							<CheckCircleFilled className="text-[11px]" />
						) : (
							<ClockCircleOutlined className="text-[11px]" />
						)}
						{isCurrentPlanPaid ? 'Current plan' : 'Free trial'}
					</span>
				)}
				{isHighlighted && (
					<span className="inline-flex w-fit rounded-full bg-base-active px-3 py-0.5 text-xs font-medium leading-[18px] text-base-primary border border-base-primary">
						Most popular
					</span>
				)}
			</div>

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

			{hasTrial(plan) && !isCurrentPlanPaid ? (
				<p className="m-0 mt-4 text-sm font-medium text-base-primary">
					{plan.trial_days}-day free trial
				</p>
			) : null}

			{isCustom ? (
				<Popover
					trigger="click"
					placement="bottom"
					arrow={false}
					content={<ContactSalesContent email={salesEmail} phone={salesPhone} />}
					styles={{ body: { padding: 0, background: 'transparent' } }}
					getPopupContainer={() => document.body}
				>
					<Button
						type={isHighlighted || hasTrial(plan) ? 'primary' : 'default'}
						size="large"
						block
						className="mt-4"
						loading={loading}
						disabled={isLocked}
						onClick={() => onAction?.(plan)}
					>
						{ctaLabel}
					</Button>
				</Popover>
			) : isBillingRestricted ? (
				<Tooltip title={BILLING_RESTRICTED_TIP}>
					<span className="mt-4 block">
						<Button
							type={isHighlighted || hasTrial(plan) ? 'primary' : 'default'}
							size="large"
							block
							loading={loading}
							disabled={isPayDisabled}
							ghost={hasTrial(plan) && !isCurrentPlan}
						>
							{ctaLabel}
						</Button>
					</span>
				</Tooltip>
			) : (
				<Button
					type={isHighlighted || hasTrial(plan) ? 'primary' : 'default'}
					size="large"
					block
					className="mt-4"
					loading={loading}
					disabled={isPayDisabled}
					ghost={hasTrial(plan) && !isCurrentPlan}
					onClick={() => onAction?.(plan)}
				>
					{ctaLabel}
				</Button>
			)}
		</div>
	);
};

export default PlanCard;
