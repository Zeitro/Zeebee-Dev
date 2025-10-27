'use client';
import { Menu } from 'antd';
import React, { useEffect, useMemo, useState } from 'react';
import { useGlobalStore } from '@/store/global';
import styles from './index.module.scss';
import { navRouter } from '@config/routes';
import { generateMenuBar, getCurrentRouterItem } from '@config/common';
import Image from 'next/image';
import { useRouter, usePathname } from 'next/navigation';
import { RightCircleTwoTone, LeftCircleTwoTone } from '@ant-design/icons';
import { useUserStore } from '@/store/user';
import { sidebarIcon } from './icon';
import useElementSize from '@/hooks/useElementSize';
import UserInfo from './userInfo';
const homePathname = '/';
const initTwoToneColor = ['var(--base-label)', 'var(--gray-light)'];
/**
 * 侧边栏组件
 */
const SidebarDom = () => {
	const router = useRouter();
	const pathname = usePathname();
	const items = generateMenuBar(navRouter, sidebarIcon).filter((i) => i.key !== '/');
	const { userInfo } = useUserStore();
	const { sidebarCollapsed, changeSidebarCollapsed } = useGlobalStore();
	const [selectedKey, setSelectedKey] = useState(homePathname);
	const [twoToneColor, setTwoToneColor] = useState(initTwoToneColor);
	const config = useMemo(() => getCurrentRouterItem(pathname, navRouter), [pathname]); // 获取当前路由的配置
	const [rect, rectDom] = useElementSize();
	const handleSelect = ({ key }) => {
		// if (selectedKey === key) return;
		router.push(key);
		setSelectedKey(key);
	};

	const handleExpand = () => {
		changeSidebarCollapsed(!sidebarCollapsed);
	};
	// 初始化根据路由选中侧边栏
	useEffect(() => {
		const routerItem = items.find((i) => pathname.indexOf(i.key) !== -1);
		setSelectedKey(routerItem?.key || homePathname);
	}, [items, pathname]);

	return config ? (
		<div
			style={{ width: sidebarCollapsed ? 56 : 216 }}
			className={[
				'relative shrink-0 h-full flex flex-col border-r border-base-border',
				styles.sidebar,
			].join(' ')}
			ref={rectDom}
		>
			<div
				className={[
					'h-14 mx-3.5 shrink-0 border-b border-base-border cursor-pointer flex items-center overflow-hidden',
					styles.logo,
				].join(' ')}
				onClick={() => handleSelect({ key: homePathname })}
			>
				<div
					className={[styles.logoBox, rect?.width < 150 && styles.logoBox_collapsed].join(
						' ',
					)}
				>
					<Image
						src="/images/logo.png"
						width={149}
						height={30}
						alt="ZeeBee Logo"
						className={[
							styles.logoImg,
							rect?.width < 150 && styles.logoImg_hidden,
						].join(' ')}
						priority
					/>
					<Image
						src="/images/logo_abbreviated.png"
						width={27}
						height={27}
						alt="ZeeBee Logo"
						className={[
							styles.logoImgSmall,
							rect?.width >= 150 && styles.logoImgSmall_hidden,
						].join(' ')}
						priority
					/>
				</div>
			</div>
			<Menu
				mode="vertical"
				selectedKeys={[selectedKey]}
				items={items}
				className={[
					'flex-1 overflow-hidden border-r-0',
					styles.menuList,
					sidebarCollapsed && styles.menuList_collapsed,
				].join(' ')}
				onClick={handleSelect}
			/>
			<UserInfo userInfo={userInfo} sidebarCollapsed={sidebarCollapsed} />
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
export default SidebarDom;
