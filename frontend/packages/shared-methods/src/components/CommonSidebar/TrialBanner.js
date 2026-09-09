'use client';

import React from 'react';
import { Tooltip } from 'antd';
import { ClockCircleOutlined, RightOutlined } from '@ant-design/icons';

/**
 * 侧栏试用期剩余天数提示；仅 trialDaysRemaining > 0 时展示。
 * 点击整条可触发升级（打开 PricingModal）。
 *
 * @param {object} props
 * @param {number} props.trialDaysRemaining 剩余试用天数
 * @param {boolean} props.sidebarCollapsed 侧栏是否折叠
 * @param {() => void} [props.onUpgradeClick] 点击升级回调
 */
const TrialBanner = ({ trialDaysRemaining = 7, sidebarCollapsed, onUpgradeClick }) => {
	const days = Number(trialDaysRemaining);
	if (!Number.isFinite(days) || days <= 0) {
		return null;
	}

	const isUrgent = true;
	const dayLabel = days === 1 ? '1 day left' : `${days} days left`;

	const toneClass = isUrgent
		? 'border-[var(--orange-medium)] bg-[var(--orange-light)] text-[var(--orange-dark)]'
		: 'border-base-border bg-base-active text-base-primary';

	if (sidebarCollapsed) {
		return (
			<div className="flex shrink-0 justify-center px-2 pb-2">
				<Tooltip title={`${dayLabel} in trial`} placement="right">
					<button
						type="button"
						className={[
							'flex h-8 w-8 items-center justify-center rounded-lg border',
							toneClass,
						].join(' ')}
						onClick={onUpgradeClick}
						aria-label={`${dayLabel} in trial`}
					>
						<ClockCircleOutlined className="text-sm" />
					</button>
				</Tooltip>
			</div>
		);
	}

	return (
		<div className="shrink-0 px-3.5 pb-2">
			<button
				type="button"
				onClick={onUpgradeClick}
				className={[
					'group flex w-full items-center gap-2 rounded-lg border px-2.5 py-2 text-left transition-opacity hover:opacity-90',
					toneClass,
				].join(' ')}
			>
				<ClockCircleOutlined className="shrink-0 text-sm" />
				<div className="min-w-0 flex-1">
					<div className="truncate text-xs font-medium leading-4">{dayLabel}</div>
					<div
						className={[
							'mt-0.5 truncate text-[11px] leading-4',
							isUrgent ? 'text-[var(--orange-dark)] opacity-80' : 'text-base-gray',
						].join(' ')}
					>
						in your free trial
					</div>
				</div>
				<span className="flex shrink-0 items-center gap-0.5 text-[11px] font-medium leading-4">
					Upgrade
					<RightOutlined className="text-[9px] opacity-70 transition-transform group-hover:translate-x-0.5" />
				</span>
			</button>
		</div>
	);
};

export default TrialBanner;
