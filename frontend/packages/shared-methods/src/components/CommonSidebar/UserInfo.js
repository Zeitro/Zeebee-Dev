'use client';

import React, { useCallback, useMemo, useState } from 'react';
import { Button, Dropdown } from 'antd';
import { RocketOutlined } from '@ant-design/icons';
import { displayValue, formatFirstLetterName } from '@repo/shared-methods/utils';
import { useUserStore } from '@repo/shared-methods/store';
import { usePathname } from 'next/navigation';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import PricingModal from '../../componentsUi/Pricing/PricingModal';
import CurrentPlanModal from '../../componentsUi/Pricing/CurrentPlanModal';

/**
 * 侧栏底部用户区：常显姓名、当前工作区、套餐；
 * 下拉顶部为账号信息区（姓名 / 邮箱 + Upgrade 同栏）；Upgrade 打开 PricingModal，
 * 套餐徽章打开 CurrentPlanModal。
 *
 * @param {object} props
 * @param {object} props.userInfo
 * @param {boolean} props.sidebarCollapsed
 * @param {object[]} props.personalList 应用侧组装的菜单项（会剥离 key=personal 的邮箱行，改由顶部账号区展示）
 * @param {(info: { key: string, domEvent: Event }) => void} props.handleClickPersonalItem
 * @param {string} [props.planName] 当前套餐名（后端未接时默认 Base）
 * @param {number|string} [props.currentPlanId] 透传 PricingModal
 * @param {boolean} [props.isLoggedIn] 是否已登录，透传 Pricing
 * @param {object} [props.pricingProps] 透传 PricingModal.pricingProps
 * @param {boolean} [props.pricingOpen] 受控：PricingModal 显隐
 * @param {(open: boolean) => void} [props.onPricingOpenChange] 受控变更回调
 * @param {object} [props.currentPlan] 当前订阅，透传 CurrentPlanModal
 */
const Userinfo = ({
	userInfo,
	sidebarCollapsed,
	personalList,
	handleClickPersonalItem,
	planName = '',
	currentPlanId,
	isLoggedIn,
	pricingProps = {},
	pricingOpen: pricingOpenProp,
	onPricingOpenChange,
	currentPlan,
}) => {
	const pathname = usePathname();
	const organization = useUserStore((state) => state.organization);
	const [selectedKey, setSelectedKey] = useState('');
	const [pricingOpenInner, setPricingOpenInner] = useState(false);
	const [currentPlanOpen, setCurrentPlanOpen] = useState(false);
	const isPricingControlled = typeof onPricingOpenChange === 'function';
	const pricingOpen = isPricingControlled ? Boolean(pricingOpenProp) : pricingOpenInner;

	const setPricingOpen = useCallback(
		(open) => {
			if (isPricingControlled) {
				onPricingOpenChange(open);
				return;
			}
			setPricingOpenInner(open);
		},
		[isPricingControlled, onPricingOpenChange],
	);

	const orgName =
		organization?.name?.trim() ||
		userInfo?.organization_name?.trim() ||
		userInfo?.organizations?.find((item) => item.is_default)?.name?.trim() ||
		'';

	const displayName = userInfo?.first_name
		? `${userInfo.first_name} ${userInfo.last_name || ''}`.trim()
		: '';

	useDeferredSyncEffect({
		deferred: () => {
			const selected = personalList.some((i) => i.key === pathname);
			if (selected) {
				setSelectedKey(pathname);
			}
		},
		deps: [pathname, personalList],
		delay: 0,
	});

	const openPricing = useCallback(
		(e) => {
			e?.preventDefault?.();
			e?.stopPropagation?.();
			setCurrentPlanOpen(false);
			setPricingOpen(true);
		},
		[setPricingOpen],
	);

	const openCurrentPlan = useCallback((e) => {
		e?.preventDefault?.();
		e?.stopPropagation?.();
		setCurrentPlanOpen(true);
	}, []);

	const menuItems = useMemo(() => {
		/** 去掉应用侧的邮箱行，改由顶部账号区统一展示 */
		const restItems = (personalList || []).filter((item) => item?.key !== 'personal');
		const trimmedRest = restItems[0]?.type === 'divider' ? restItems.slice(1) : restItems;

		const accountHeader = {
			key: 'account-header',
			className:
				'!h-auto !cursor-default !bg-transparent hover:!bg-transparent [&_.ant-dropdown-menu-title-content]:!whitespace-normal',
			label: (
				<div
					className="w-[220px] px-1 py-1"
					onClick={(e) => e.stopPropagation()}
					onMouseDown={(e) => e.stopPropagation()}
				>
					{displayName ? (
						<div className="truncate text-sm font-medium leading-5 text-base-text">
							{displayName}
						</div>
					) : null}
					<div
						className={[
							'truncate text-xs leading-5 text-base-label',
							displayName ? 'mt-0.5' : '',
						].join(' ')}
					>
						{displayValue(userInfo?.email)}
					</div>
					{!!currentPlan?.id && (
						<Button
							className="mt-2.5 flex w-full items-center justify-center gap-1.5 rounded-lg border border-base-border bg-white px-3 py-1.5 text-sm font-medium text-base-text hover:border-base-primary hover:text-base-primary"
							onClick={openPricing}
						>
							<RocketOutlined className="text-xs" />
							Upgrade plan
						</Button>
					)}
				</div>
			),
		};

		return [accountHeader, { type: 'divider' }, ...trimmedRest];
	}, [currentPlan?.id, displayName, openPricing, personalList, userInfo?.email]);

	const handleMenuClick = (info) => {
		if (info.key === 'account-header') {
			return;
		}
		handleClickPersonalItem?.(info);
	};

	return (
		<>
			<Dropdown
				menu={{
					items: menuItems,
					selectedKeys: [selectedKey],
					onClick: handleMenuClick,
				}}
				placement="topLeft"
				trigger={['click']}
				className="overflow-hidden"
			>
				<div
					className="flex h-16 cursor-pointer items-center gap-2.5 border-t border-solid border-base-border p-3.5 text-base-primary"
					onClick={(e) => e.preventDefault()}
				>
					<div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-base-active text-sm text-base-label">
						{formatFirstLetterName(userInfo) || ''}
					</div>
					<div
						className="min-w-0 flex-1 overflow-hidden"
						style={{
							opacity: sidebarCollapsed ? 0 : 1,
							transition: 'opacity 0.3s cubic-bezier(0.2, 0, 0, 1) 0s',
						}}
					>
						<p className="m-0 truncate pb-[3px] text-base/4 font-medium text-black-dark">
							{displayValue(displayName)}
						</p>
						<div className="flex min-w-0 items-center gap-1.5">
							{orgName ? (
								<span
									className="min-w-0 flex-1 truncate text-xs text-base-gray"
									title={orgName}
								>
									{orgName}
								</span>
							) : (
								<span className="min-w-0 flex-1 truncate text-xs text-base-label">
									{displayValue(userInfo?.email)}
								</span>
							)}
							{planName && (
								<button
									type="button"
									className="shrink-0 rounded border border-base-border bg-gray-light px-1.5 py-0 text-[11px] leading-5 text-base-gray hover:border-base-primary hover:text-base-primary"
									title="View current plan"
									onClick={openCurrentPlan}
								>
									{planName}
								</button>
							)}
						</div>
					</div>
				</div>
			</Dropdown>

			<PricingModal
				open={pricingOpen}
				onCancel={() => setPricingOpen(false)}
				pricingProps={{
					currentPlanId,
					isLoggedIn,
					...pricingProps,
				}}
			/>

			<CurrentPlanModal
				open={currentPlanOpen}
				onCancel={() => setCurrentPlanOpen(false)}
				subscription={currentPlan}
				onUpgrade={() => {
					setCurrentPlanOpen(false);
					setPricingOpen(true);
				}}
			/>
		</>
	);
};

export default Userinfo;
