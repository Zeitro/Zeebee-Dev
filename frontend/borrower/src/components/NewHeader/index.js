'use client';
import { Breadcrumb } from 'antd';
import React, { useLayoutEffect, useMemo, useState } from 'react';
import { usePathname } from 'next/navigation';
import { navRouter } from '@config/routes';
import { getCurrentRouterItem, flattenNavRouter, getCrumbs } from '@config/common';
import Link from 'next/link';
import styles from './index.module.scss';
import { CustomContent } from './customContent';
import eventBus from '@/lib/bus';
/**
 * 面包屑导航
 */
const NewHeader = () => {
	const pathname = usePathname();
	const config = useMemo(() => getCurrentRouterItem(pathname, navRouter), [pathname]); // 获取当前路由的配置
	const [crumbs, setCrumbs] = useState([]);
	const itemRender = (currentRoute, params, items) => {
		const isLast = currentRoute?.path === items[items.length - 1]?.path;
		return isLast ? (
			<span>{currentRoute.title}</span>
		) : (
			<Link href={currentRoute.path || ''}>{currentRoute.title}</Link>
		);
	};
	useLayoutEffect(() => {
		const flatRoutes = flattenNavRouter(navRouter);
		const breadcrumbs = getCrumbs(pathname, flatRoutes);
		setCrumbs([navRouter[0], ...breadcrumbs]);
	}, [pathname]);

	// 有动态的面包屑内容时，可调用changeCrumbs修改，参数config = {1: '动态内容'}
	useLayoutEffect(() => {
		const changeCrumbs = (config = {}) => {
			const indexArr = Object.keys(config);
			const newCrumbs = crumbs.map((item, index) => {
				if (indexArr.includes(index.toString())) {
					item.title = config[index];
				}
				return item;
			});
			setCrumbs([...newCrumbs]);
		};
		eventBus.on('changeCrumbs', changeCrumbs);
		return () => {
			eventBus.off('changeCrumbs', changeCrumbs);
		};
	}, [crumbs]);
	return config && !config?.hideCrumbs ? (
		<div className={`h-[84px] p-4 shrink-0 flex items-center justify-between ${styles.crumbs}`}>
			<div>
				<Breadcrumb items={crumbs} itemRender={itemRender} />
				<div className="text-lg font-medium text-gray-medium">
					{crumbs?.[crumbs.length - 1]?.title}
				</div>
			</div>
			<CustomContent />
		</div>
	) : (
		''
	);
};
export default NewHeader;
