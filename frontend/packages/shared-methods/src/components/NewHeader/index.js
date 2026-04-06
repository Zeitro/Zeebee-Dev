'use client';
import { Breadcrumb } from 'antd';
import React, { useLayoutEffect, useMemo, useState } from 'react';
import { usePathname } from 'next/navigation';
import { navRouter } from '@config/routes';
import Link from 'next/link';
import {
	eventBus,
	getCurrentRouterItem,
	flattenNavRouter,
	getCrumbs,
} from '@repo/shared-methods/utils';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
/**
 * 面包屑导航
 * isShowCrumbs：true时不受router内配置的hideCrumbs影响，直接显示
 */
export const NewHeader = ({ children, isShowCrumbs }) => {
	const pathname = usePathname();
	const isShow = useMemo(() => {
		const item = getCurrentRouterItem(pathname, navRouter);
		if (isShowCrumbs) return true;
		return item && !item?.hideCrumbs;
	}, [pathname, isShowCrumbs]); // 获取当前路由的配置
	const [crumbs, setCrumbs] = useState([]);
	const itemRender = (currentRoute, params, items) => {
		const isLast = currentRoute?.path === items[items.length - 1]?.path;
		return isLast ? (
			<span>{currentRoute.title}</span>
		) : (
			<Link href={currentRoute.path || ''}>{currentRoute.title}</Link>
		);
	};

	useDeferredSyncEffect({
		deferred: () => {
			const flatRoutes = flattenNavRouter(navRouter);
			const breadcrumbs = getCrumbs(pathname, flatRoutes);
			setCrumbs([navRouter[0], ...breadcrumbs]);
		},
		deps: [pathname],
		isUseMicrotask: true,
	});

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
	return isShow ? (
		<div
			className="px-4 py-4 shrink-0 flex items-center justify-between border-b border-base-border"
			style={{
				boxShadow:
					'0px 1px 2px 0px rgba(0, 0, 0, 0.03),0px 1px 6px -1px rgba(0, 0, 0, 0.02),0px 2px 4px 0px rgba(0, 0, 0, 0.02)',
				height: 56,
			}}
		>
			<div>
				<Breadcrumb items={crumbs} itemRender={itemRender} />
				{/* <div className="font-medium text-gray-medium text-base">
					{crumbs?.[crumbs.length - 1]?.title}
				</div> */}
			</div>
			<div>{children}</div>
		</div>
	) : (
		''
	);
};
export default NewHeader;
