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
const homePathname = '/';
/**
 * 侧边栏组件
 */
const SidebarDom = () => {
	const router = useRouter();
	const pathname = usePathname();
	const { siberCollapsed, changeSiberCollapsed } = useGlobalStore();
	const [selectedKey, setSelectedKey] = useState(homePathname);
	const config = useMemo(() => getCurrentRouterItem(pathname, navRouter), [pathname]); // 获取当前路由的配置

	const handleSelect = ({ key }) => {
		if (selectedKey === key) return;
		router.push(key);
		setSelectedKey(key);
	};

	const handleExpand = () => {
		changeSiberCollapsed(!siberCollapsed);
	};
	// 初始化根据路由选中侧边栏
	useEffect(() => {
		setSelectedKey((v) => (v !== pathname ? pathname : homePathname));
	}, []);

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
					'h-14 mx-3.5 shrink-0 border-b border-base-border cursor-pointer',
					styles.logo,
				].join(' ')}
				onClick={() => handleSelect({ key: homePathname })}
			>
				Zeebee
			</div>
			<Menu
				mode="vertical"
				selectedKeys={[selectedKey]}
				items={generateMenuBar(navRouter).filter((i) => i.key !== '/')}
				className={[
					'flex-1 overflow-hidden border-r-0',
					styles.menuList,
					siberCollapsed && styles.menuList_collapsed,
				].join(' ')}
				onSelect={handleSelect}
			/>
			<div
				className={`h-16 border-t border-solid border-base-border text-base-primary flex items-center justify-center p-3.5 gap-2.5 overflow-hidden`}
			>
				<Image
					src="/globe.svg"
					width={32}
					height={32}
					alt="Picture of the author"
					className="shrink-0"
				/>
				{!siberCollapsed && (
					<div className="flex-1">
						<p className={`text-base/4 text-black-dark font-medium pb-[3px]`}>
							username
						</p>
						<span className="text-gray-medium">kevinmo@email.com</span>
					</div>
				)}
			</div>
			<div
				className="absolute right-[-10px] top-[19px] text-base-primary cursor-pointer"
				onClick={handleExpand}
			>
				{siberCollapsed ? (
					<RightCircleTwoTone className="text-base-primary" />
				) : (
					<LeftCircleTwoTone className="text-base-primary" />
				)}
			</div>
		</div>
	) : (
		''
	);
};
export default SidebarDom;
