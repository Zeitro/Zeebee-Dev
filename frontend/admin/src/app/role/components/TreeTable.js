'use client';
import React, { useEffect } from 'react';
import { FilterTable } from '@repo/shared-methods';
import useColumn from './useColumn';
import { createRequest, eventBus } from '@repo/shared-methods/utils';

const TreeTable = ({
	hasPagination,
	filterWrapClassName,
	dataSource,
	loading,
	type = 'Borrower', // 默认是Borrower
	eventBusPrefix,
	menuKey = 'Api',
	api,
}) => {
	const { getColumnsByMenuKey, tableFilterRef } = useColumn({ eventBusPrefix, api });
	const { request } = createRequest();

	// 更新列表
	useEffect(() => {
		if (type && eventBusPrefix) {
			const searchList = () => tableFilterRef?.current?.handleSearch?.();
			eventBus.on(`${eventBusPrefix}${type}`, searchList);
			return () => eventBus.off(`${eventBusPrefix}${type}`, searchList);
		}
	}, [type, eventBusPrefix]);

	return (
		<FilterTable
			key={type}
			ref={tableFilterRef}
			columns={getColumnsByMenuKey(menuKey)}
			hasFilterField={false}
			pagination={hasPagination}
			filterWrapClassName={filterWrapClassName}
			dataSource={dataSource}
			dataSourceLoading={loading}
			getData={() => {
				const payload = new URLSearchParams({ type });
				const url = `${api}?${payload}`;
				return request(url);
			}}
			hasPaginationRequest={false}
			scroll={{ x: 'max-content' }}
		/>
	);
};
export default TreeTable;
