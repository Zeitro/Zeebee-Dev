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
import { formatFirstLetterName } from '@/lib/format';
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
	const { siberCollapsed, changeSiberCollapsed } = useGlobalStore();
	const [selectedKey, setSelectedKey] = useState(homePathname);
	const [twoToneColor, setTwoToneColor] = useState(initTwoToneColor);
	const config = useMemo(() => getCurrentRouterItem(pathname, navRouter), [pathname]); // 获取当前路由的配置

	const handleSelect = ({ key }) => {
		// if (selectedKey === key) return;
		router.push(key);
		setSelectedKey(key);
	};

	const handleLogout = () => {
		router.push('/login');
	};

	const handleExpand = () => {
		changeSiberCollapsed(!siberCollapsed);
	};
	// 初始化根据路由选中侧边栏
	useEffect(() => {
		const routerItem = items.find((i) => pathname.indexOf(i.key) !== -1);
		setSelectedKey(routerItem?.key || homePathname);
	}, [pathname]);

	return config ? (
		<div
			style={{ width: siberCollapsed ? 56 : 216 }}
			className={[
				'relative shrink-0 h-full flex flex-col border-r border-base-border',
				styles.sidebar,
			].join(' ')}
		>
			<div
				className={[
					'h-14 mx-3.5 shrink-0 border-b border-base-border cursor-pointer flex items-center',
					styles.logo,
				].join(' ')}
				onClick={() => handleSelect({ key: homePathname })}
			>
				<Image src="/images/logo.svg" width={95} height={30} alt="ZeeBee Logo" />
			</div>
			<Menu
				mode="vertical"
				selectedKeys={[selectedKey]}
				items={items}
				className={[
					'flex-1 overflow-hidden border-r-0',
					styles.menuList,
					siberCollapsed && styles.menuList_collapsed,
				].join(' ')}
				onClick={handleSelect}
			/>
			<div
				className={`h-16 border-t border-solid border-base-border text-base-primary flex items-center justify-center p-3.5 gap-2.5 overflow-hidden`}
				onClick={() => handleLogout()}
			>
				{/* <Image
					src="/globe.svg"
					width={32}
					height={32}
					alt="Picture of the author"
					className="shrink-0"
				/> */}
				<div
					className={[
						`w-8 h-8 rounded-full flex items-center justify-center shrink-0 text-sm text-background`,
						styles.bg,
					].join(' ')}
				>
					{formatFirstLetterName(userInfo)}
				</div>
				{!siberCollapsed && (
					<div className="flex-1">
						<p className={`text-base/4 text-black-dark font-medium pb-[3px]`}>
							{userInfo?.first_name} {userInfo?.last_name}
						</p>
						<span className="text-gray-medium text-ellipsis overflow-hidden whitespace-nowrap">
							{userInfo?.email}
						</span>
					</div>
				)}
			</div>
			<div
				className="absolute right-[-8px] top-[18px] cursor-pointer"
				onClick={handleExpand}
				onMouseLeave={() => setTwoToneColor(initTwoToneColor)}
				onMouseEnter={() => setTwoToneColor(['var(--base-primary)', 'var(--base-active)'])}
			>
				{siberCollapsed ? (
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
