'use client';

import React, { Suspense } from 'react';
import TabBar from './TabBar';

/**
 * 替代 NewHeader 面包屑：左侧 TabBar + 右侧自定义操作区。
 * TabBar 依赖 useSearchParams，需 Suspense 包裹。
 */
const TabAppHeader = ({ children }) => {
	return (
		<div
			className="px-4 py-2 shrink-0 flex items-center justify-between gap-4 border-b border-base-border"
			style={{
				boxShadow:
					'0px 1px 2px 0px rgba(0, 0, 0, 0.03),0px 1px 6px -1px rgba(0, 0, 0, 0.02),0px 2px 4px 0px rgba(0, 0, 0, 0.02)',
				minHeight: 56,
			}}
		>
			<Suspense fallback={<div className="flex-1 h-8" />}>
				<TabBar />
			</Suspense>
			<div className="shrink-0">{children}</div>
		</div>
	);
};

export default TabAppHeader;
