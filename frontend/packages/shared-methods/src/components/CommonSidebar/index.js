'use client';
import { Menu } from 'antd';
import React, { useMemo, useState } from 'react';
import { useGlobalStore } from '@/store/global';
import styles from './index.module.scss';
import { useRouter, usePathname } from 'next/navigation';
import { RightCircleTwoTone, LeftCircleTwoTone } from '@ant-design/icons';
import { useUserStore } from '@repo/shared-methods/store';
import UserInfo from './UserInfo';
import { useDeferredSyncEffect, useElementSize } from '@repo/shared-methods/hooks';
import { generateMenuBar, resolveSidebarMenuState } from '@repo/shared-methods/utils';
import { iconData } from '@repo/shared-methods/enum';

const homePathname = '/';
const initTwoToneColor = ['var(--base-label)', 'var(--gray-light)'];
/**
 * 侧边栏组件
 */
export const CommonSidebar = ({
	routerList,
	personalList,
	handleClickPersonalItem,
	renderLogo,
}) => {
	const router = useRouter();
	const pathname = usePathname();
	const { userInfo } = useUserStore();
	const { sidebarCollapsed, changeSidebarCollapsed } = useGlobalStore();
	const [selectedKey, setSelectedKey] = useState('');
	const [openKey, setOpenKey] = useState('');
	const [twoToneColor, setTwoToneColor] = useState(initTwoToneColor);
	const [rect, rectDom] = useElementSize();

	// 详情/设置等未展示菜单的路由 → 选中并展开最近可见父级
	const menuState = useMemo(
		() => resolveSidebarMenuState(pathname, routerList, homePathname),
		[pathname, routerList],
	);
	const items = useMemo(() => {
		return generateMenuBar([...routerList], iconData, sidebarCollapsed);
	}, [sidebarCollapsed, routerList]);
	const handleSelect = ({ key }) => {
		// if (selectedKey === key) return;
		router.push(key);
		setSelectedKey(key);
	};

	const handleExpand = () => {
		changeSidebarCollapsed(!sidebarCollapsed);
	};
	// 初始化根据路由选中侧边栏（含无菜单页回落父级）
	useDeferredSyncEffect({
		predicate: () => selectedKey !== menuState.selectedKey,
		deferred: () => {
			setSelectedKey(menuState.selectedKey);
		},
		deps: [pathname, menuState.selectedKey, selectedKey],
		delay: 0,
	});

	useDeferredSyncEffect({
		predicate: () => !sidebarCollapsed,
		deferred: () => {
			setOpenKey(menuState.openKey);
		},
		deps: [pathname, menuState.openKey, sidebarCollapsed],
		delay: 0,
	});

	// hideSidebar隐藏左侧菜单栏；无精确路由配置时仍展示，以便详情页回落父级选中
	return !menuState.hideSidebar ? (
		<div
			style={{ width: sidebarCollapsed ? 63 : 216 }}
			className={[
				'relative shrink-0 h-full flex flex-col border-r border-base-border',
				styles.sidebar,
			].join(' ')}
			ref={rectDom}
		>
			<div
				className={[
					'h-14 mx-3.5 shrink-0 border-b border-base-border cursor-pointer flex items-center',
					styles.logo,
				].join(' ')}
				onClick={() => handleSelect({ key: homePathname })}
			>
				{renderLogo?.({ rect })}
			</div>
			<Menu
				mode="inline"
				selectedKeys={[selectedKey]}
				openKeys={[openKey]}
				items={items}
				className={['flex-1 overflow-hidden'].join(' ')}
				onClick={handleSelect}
				onOpenChange={(keys) => {
					setOpenKey(keys[keys.length - 1]);
				}}
				inlineCollapsed={sidebarCollapsed}
				inlineIndent={3}
				getPopupContainer={() => {
					// 将弹出菜单挂载到 body 或指定容器，避免被父级遮挡
					const customRoot = document.querySelector('.customRoot');
					return customRoot;
				}}
			/>
			<UserInfo
				userInfo={userInfo}
				sidebarCollapsed={sidebarCollapsed}
				personalList={personalList}
				handleClickPersonalItem={handleClickPersonalItem}
			/>
			<div
				className="absolute right-[-8px] top-[18px] cursor-pointer"
				onClick={handleExpand}
				onMouseLeave={() => setTwoToneColor(initTwoToneColor)}
				onMouseEnter={() => setTwoToneColor(['var(--base-primary)', 'var(--base-active)'])}
			>
				{sidebarCollapsed ? (
					<RightCircleTwoTone twoToneColor={twoToneColor} />
				) : (
					<LeftCircleTwoTone twoToneColor={twoToneColor} />
				)}
			</div>
		</div>
	) : (
		''
	);
};
export default CommonSidebar;
