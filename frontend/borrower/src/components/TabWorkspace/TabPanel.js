'use client';

import React from 'react';

/**
 * Tab 保活容器：非激活 Tab 用 hidden 保留 DOM，关闭 Tab 时随父级卸载。
 * @param {{ isActive: boolean, keepAlive: boolean, children: React.ReactNode }} props
 */
const TabPanel = ({ isActive, keepAlive, children }) => {
	if (!keepAlive && !isActive) {
		return null;
	}

	return (
		<div
			className={[
				'flex flex-col h-full flex-1 overflow-hidden',
				isActive ? '' : 'hidden',
			].join(' ')}
			aria-hidden={!isActive}
		>
			{children}
		</div>
	);
};

export default TabPanel;
