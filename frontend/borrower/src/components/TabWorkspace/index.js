'use client';

import React, { Suspense } from 'react';
import { Spin } from 'antd';
import TabWorkspaceContent from './TabWorkspaceContent';

const WorkspaceFallback = () => (
	<div className="flex-1 flex items-center justify-center">
		<Spin />
	</div>
);

/**
 * Tab 内容区容器：P1 对已注册路由做 Panel 保活，其余仍渲染 Next children。
 * @param {{ children: React.ReactNode }} props
 */
const TabWorkspace = ({ children }) => {
	return (
		<Suspense fallback={<WorkspaceFallback />}>
			<TabWorkspaceContent>{children}</TabWorkspaceContent>
		</Suspense>
	);
};

export default TabWorkspace;
