'use client';

import React from 'react';

/**
 * Tab 内容区容器。P0 仍渲染 Next {children}；P1 起在此按 tabRegistry 做保活。
 * @param {{ children: React.ReactNode }} props
 */
const TabWorkspace = ({ children }) => {
	return <div className="flex-1 overflow-hidden flex flex-col">{children}</div>;
};

export default TabWorkspace;
